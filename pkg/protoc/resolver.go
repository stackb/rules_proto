package protoc

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"sort"
	"strings"
	"unsafe"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
)

const protoImportType = "protobuf:import"

const debugResolver = true

type ImportResolver interface {
	// Requires returns any previously provided labels associated with the given
	// kind and import.
	Requires(kind string, imp string) []resolve.FindResult
	// Provides records the association between a kind of rule, a proto import
	// statement, and the label that provides the import.
	Provides(kind string, imp string, location label.Label)
}

// ImportCrossResolver handles dependency resolution.
type ImportCrossResolver interface {
	resolve.CrossResolver
	ImportResolver

	// LoadImportsFile loads csv file to populate the resolver
	LoadImportsFile(filename string) error
	// SaveImportsFile write a csv file
	SaveImportsFile(repoName, filename string) error
	// InstallResolveOverrides adds configured resolve entries into the resolve config.
	InstallResolveOverrides(c *config.Config)
}

// GlobalResolver returns a reference to the global ImportResolver
func GlobalResolver() ImportCrossResolver {
	return globalImportResolver
}

// importLabels records which labels are associated with a given proto import
// statement.
type importLabels map[string][]label.Label

// globalImportResolver is the default resolver singleton.
var globalImportResolver = &resolver{
	// known is a mapping between kind and importLabel map
	known: make(map[string]importLabels),
}

// resolver implements ImportResolver.
type resolver struct {
	known map[string]importLabels
}

// LoadImportsFile reads a protoresolve csv file.
func (r *resolver) LoadImportsFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return r.LoadReader(f)
}

// LoadImportsFile reads input and returns a list of items.  Comment lines beginning
// with '#' are ignored.
func (r *resolver) LoadReader(in io.Reader) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ",", 3)
		if len(parts) == 4 {
			lang := parts[0]
			// TODO: perhaps there will other types in the future
			if lang != protoImportType {
				continue
			}
			kind := parts[1]
			imp := parts[2]
			lbl, err := label.Parse(parts[3])
			if err != nil {
				return fmt.Errorf("malformed label at position 4 in %s: %v", line, err)
			}
			r.Provides(kind, imp, lbl)
		}
	}
	return nil
}

func (r *resolver) SaveImports(out io.Writer, repoName string) {
	kinds := make([]string, 0)
	for k := range r.known {
		kinds = append(kinds, k)
	}
	sort.Strings(kinds)
	for _, kind := range kinds {
		imports := r.known[kind]
		imps := make([]string, 0)
		for imp := range imports {
			imps = append(imps, imp)
		}
		sort.Strings(imps)
		for _, imp := range imps {
			labels := imports[imp]
			for _, lbl := range labels {
				l := label.New(lbl.Repo, lbl.Pkg, lbl.Name)
				if l.Repo == "" {
					l.Repo = repoName
				}
				fmt.Fprintf(out, "%s,%s,%s,%s\n", protoImportType, kind, imp, l)
			}
		}
	}
}

func (r *resolver) SaveImportsFile(repoName, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("save imports file: %w", err)
	}
	defer f.Close()

	fmt.Fprintf(f, "# GENERATED FILE, DO NOT EDIT (created by gazelle)\n")
	fmt.Fprintf(f, "# type,kind,key,value\n")
	r.SaveImports(f, repoName)
	return nil
}

// CrossResolve provides dependency resolution logic for the proto language extension.
func (r *resolver) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	switch lang {
	case "go":
		return r.Requires("go_library", imp.Imp)
	case "proto":
		return r.Requires("proto_library", imp.Imp)
	default:
		return r.Requires(lang, imp.Imp)
	}
}

func (r *resolver) Requires(kind, imp string) []resolve.FindResult {
	known := r.known[kind]
	if known == nil {
		return nil
	}
	if got, ok := known[imp]; ok {
		res := make([]resolve.FindResult, len(got))
		for i, l := range got {
			res[i] = resolve.FindResult{Label: l}
		}
		return res
	}
	return nil
}

func (r *resolver) Provides(kind, imp string, loc label.Label) {
	known, ok := r.known[kind]
	if !ok {
		known = make(map[string][]label.Label)
		r.known[kind] = known
	}
	for _, v := range known[imp] {
		if v == loc {
			if debugResolver {
				log.Println(kind, imp, "PROVIDES (duplicate)", loc)
			}
			return
		}
	}
	if debugResolver {
		log.Println(kind, imp, "PROVIDES", loc)
	}
	known[imp] = append(known[imp], loc)
}

func (r *resolver) InstallResolveOverrides(c *config.Config) {
	// The resolve config has already processed resolve directives, and there's
	// no public API. Take somewhat extreme measures to augment it's internal
	// override list via unsafe memory reallocation.
	overrides := make([]overrideSpec, 0)

	for imp, labels := range r.known["proto_library"] {
		for _, lbl := range labels {
			overrides = append(overrides, overrideSpec{
				imp: resolve.ImportSpec{
					Lang: "proto",
					Imp:  imp,
				},
				lang: "proto",
				dep:  lbl,
			})
		}
	}

	if len(overrides) > 0 {
		rewriteResolveConfigOverrides(getResolveConfig(c), overrides)
	}
}

// ResolveImports is a utility function that returns a matching list of labels
// for the given import list.
func ResolveImports(resolver ImportResolver, kind string, imports []string) []label.Label {
	deps := make([]label.Label, 0)
	for _, imp := range imports {
		result := resolver.Requires(kind, imp)
		if len(result) > 0 {
			first := result[0]
			deps = append(deps, first.Label)
			if debugResolver {
				log.Println(kind, imp, "HIT", first.Label)
			}
		} else {
			if debugResolver {
				log.Println(kind, imp, "MISS", resolver)
			}
		}
	}
	return deps
}

// ResolveImportsString is a utility function that returns a matching list of labels
// for the given import list.
func ResolveImportsString(resolver ImportResolver, rel, kind string, imports []string) []string {
	match := ResolveImports(resolver, kind, imports)
	deps := make([]string, len(match))
	for i, l := range match {
		deps[i] = l.Rel("", rel).String()
	}
	return deps
}

// getResolveConfig returns the resolve.resolveConfig
func getResolveConfig(c *config.Config) interface{} {
	return c.Exts["_resolve"]
}

// rewriteResolveConfigOverrides reads the existing private attribute and
// appends more overrides.
func rewriteResolveConfigOverrides(rc interface{}, more []overrideSpec) {
	rcv := reflect.ValueOf(rc).Elem()
	val := reflect.Indirect(rcv)
	member := val.FieldByName("overrides")
	ptrToOverrides := unsafe.Pointer(member.UnsafeAddr())
	overrides := (*[]overrideSpec)(ptrToOverrides)

	// create new array
	newOverrides := make([]overrideSpec, 0)
	newOverrides = append(newOverrides, *overrides...)
	newOverrides = append(newOverrides, more...)

	// reassign memory value
	*overrides = newOverrides
}

// overrideSpec is a copy of the same private type in resolve/config.go.  It must be
// kept in sync with the original to avoid discrepancy with the expected memory
// layout.
type overrideSpec struct {
	imp  resolve.ImportSpec
	lang string
	dep  label.Label
}
