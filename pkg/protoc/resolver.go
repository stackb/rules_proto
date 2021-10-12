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

const debugResolver = false

const (
	// ResolveProvidesKey is the key expected to store a string slice that
	// informs what imports a rule provides.
	ResolveProvidesKey = "_resolve_provides"
	// RuleProviderKey stores the rule provider implementation for a given rule.
	RuleProviderKey = "_rule_provider"
)

type ImportResolver interface {
	// Resolve returns any previously provided labels associated with the given
	// kind and import.
	Resolve(kind, attr, imp string) []resolve.FindResult
	// Provide records the association between a rule kind+attr, the value of
	// the attr, and the label that provides the value.
	Provide(kind string, attr, val string, location label.Label)
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
			lbl, err := label.Parse(parts[0])
			if err != nil {
				return fmt.Errorf("malformed label at position 4 in %s: %v", line, err)
			}
			kind := parts[1]
			attr := parts[2]
			imp := parts[3]
			r.Provide(kind, attr, imp, lbl)
		}
	}
	return nil
}

func (r *resolver) SaveImports(out io.Writer, repoName string) {
	keys := make([]string, 0)
	for k := range r.known {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		imports := r.known[key]
		imps := make([]string, 0)
		for imp := range imports {
			imps = append(imps, imp)
		}
		sort.Strings(imps)
		kind, attr := keyKind(key)
		for _, imp := range imps {
			labels := imports[imp]
			for _, lbl := range labels {
				l := label.New(lbl.Repo, lbl.Pkg, lbl.Name)
				if l.Repo == "" {
					l.Repo = repoName
				}
				fmt.Fprintf(out, "%s,%s,%s,%s\n", l, kind, attr, imp)
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
	fmt.Fprintf(f, "# lang,imp.lang,imp,label\n")
	r.SaveImports(f, repoName)
	return nil
}

// CrossResolve provides dependency resolution logic for the proto language extension.
func (r *resolver) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	switch lang {
	case "go":
		return r.Resolve("go_library", "importpath", imp.Imp)
	case "proto":
		return r.Resolve("proto_library", "srcs", imp.Imp)
	default:
		return r.Resolve(lang, "", imp.Imp)
	}
}

func (r *resolver) Resolve(kind, attr, imp string) []resolve.FindResult {
	key := kindKey(kind, attr)
	known := r.known[key]
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

func (r *resolver) Provide(kind, attr, imp string, loc label.Label) {
	key := kindKey(kind, attr)
	known, ok := r.known[key]
	if !ok {
		known = make(map[string][]label.Label)
		r.known[key] = known
	}
	for _, v := range known[imp] {
		if v == loc {
			if debugResolver {
				log.Println(key, imp, "PROVIDES (duplicate)", loc)
			}
			return
		}
	}
	if debugResolver {
		log.Println(key, imp, "PROVIDES", loc)
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
func ResolveImports(resolver ImportResolver, kind, attr string, imports []string) []label.Label {
	deps := make([]label.Label, 0)
	for _, imp := range imports {
		result := resolver.Resolve(kind, attr, imp)
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
func ResolveImportsString(resolver ImportResolver, rel, kind, attr string, imports []string) []string {
	match := ResolveImports(resolver, kind, attr, imports)
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

	// create new array: FindRuleWithOverride searches last entries first, so
	// respect the users own resolve directives by putting them last
	newOverrides := make([]overrideSpec, 0)
	newOverrides = append(newOverrides, more...)
	newOverrides = append(newOverrides, *overrides...)

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

func kindKey(kind, attr string) string {
	return kind + "//" + attr
}

func keyKind(key string) (string, string) {
	parts := strings.SplitN(key, "//", 2)
	return parts[0], parts[1]
}
