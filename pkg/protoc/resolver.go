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

// globalImportResolver is the default resolver singleton.
var globalImportResolver = NewImportResolver(&ImportResolverOptions{
	Debug:  false,
	Printf: log.Printf,
}).(*resolver)

const (
	// ResolveProvidesKey is the key expected to store a string slice that
	// informs what imports a rule provides.
	ResolveProvidesKey = "_resolve_provides"
)

type ImportResolver interface {
	// Resolve returns any previously provided labels associated with the given
	// kind and import.
	Resolve(lang, impLang, imp string) []resolve.FindResult
	// Provide records the association between a rule kind+attr, the value of
	// the attr, and the label that provides the value.
	Provide(lang string, impLang, val string, location label.Label)
}

// ImportCrossResolver handles dependency resolution.
type ImportCrossResolver interface {
	resolve.CrossResolver
	ImportResolver

	// LoadFile loads csv file to populate the resolver
	LoadFile(filename string) error
	// SaveFile writes a csv file
	SaveFile(filename, repoName string) error
	// Install adds configured resolve entries into the resolve config.
	Install(c *config.Config)
}

// GlobalResolver returns a reference to the global ImportResolver
func GlobalResolver() ImportCrossResolver {
	return globalImportResolver
}

type ImportResolverOptions struct {
	Printf func(format string, args ...interface{})
	Debug  bool
}

func NewImportResolver(options *ImportResolverOptions) ImportResolver {
	return &resolver{
		known:   make(map[string]importLabels),
		options: options,
	}
}

// importLabels records which labels are associated with a given proto import
// statement.
type importLabels map[string][]label.Label

// resolver implements ImportResolver.
type resolver struct {
	options *ImportResolverOptions
	// known is a mapping between lang and importLabel map
	known map[string]importLabels
}

// LoadFile reads a protoresolve csv file.
func (r *resolver) LoadFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return r.Load(f)
}

// Load reads input and returns a list of items.  Comment lines beginning
// with '#' are ignored.
func (r *resolver) Load(in io.Reader) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ",", 4)
		if len(parts) != 4 {
			log.Printf("warn: parse %q, expected 4 items, got %d", line, len(parts))
			continue
		}
		lang := parts[0]
		impLang := parts[1]
		imp := parts[2]
		lbl, err := label.Parse(parts[3])
		if err != nil {
			return fmt.Errorf("malformed label at position 4 in %s: %v", line, err)
		}
		r.Provide(lang, impLang, imp, lbl)
	}
	return nil
}

func (r *resolver) Save(out io.Writer, repoName string) {
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
		lang, impLang := keyLang(key)
		for _, imp := range imps {
			labels := imports[imp]
			for _, lbl := range labels {
				// skip external labels, these represent externally loaded
				// entries and we don't write transitive resolves
				if lbl.Repo != "" {
					continue
				}
				l := label.New(repoName, lbl.Pkg, lbl.Name)
				fmt.Fprintf(out, "%s,%s,%s,%s\n", lang, impLang, imp, l)
			}
		}
	}
}

func (r *resolver) SaveFile(filename, repoName string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("save imports file: %w", err)
	}

	fmt.Fprintf(f, "# GENERATED FILE, DO NOT EDIT (created by gazelle)\n")
	fmt.Fprintf(f, "# lang,imp.lang,imp,label\n")

	r.Save(f, repoName)
	if err := f.Close(); err != nil {
		return err
	}

	// log.Println("Wrote resolve file:", filename)
	return nil
}

// CrossResolve provides dependency resolution logic for the protobuf language extension.
func (r *resolver) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	res := r.Resolve(lang, imp.Lang, imp.Imp)
	if r.options.Debug {
		r.options.Printf("cross-resolve %s %s %s (%d results)", lang, imp.Lang, imp.Imp, len(res))
	}
	return res
}

func (r *resolver) Resolve(lang, impLang, imp string) []resolve.FindResult {
	key := langKey(lang, impLang)
	known := r.known[key]
	if known == nil {
		known = r.known[lang]
	}
	if known == nil {
		if r.options.Debug {
			r.options.Printf("resolve miss %s: no records under language %q", imp, key)
		}
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

func (r *resolver) Provide(lang, impLang, imp string, from label.Label) {
	key := langKey(lang, impLang)
	known, ok := r.known[key]
	if !ok {
		known = make(map[string][]label.Label)
		r.known[key] = known
	}
	for _, v := range known[imp] {
		if v == from {
			return
		}
	}
	if r.options.Debug {
		r.options.Printf("resolver %v provides %s %s", from, key, imp)
	}
	known[imp] = append(known[imp], from)
}

func (r *resolver) Install(c *config.Config) {
	// The resolve config has already processed resolve directives, and there's
	// no public API. Take somewhat extreme measures to augment it's internal
	// override list via unsafe memory reallocation.
	overrides := make([]overrideSpec, 0)

	for key, known := range r.known {
		lang, impLang := keyLang(key)
		for imp, lbls := range known {
			for _, lbl := range lbls {
				overrides = append(overrides, overrideSpec{
					imp: resolve.ImportSpec{
						Lang: impLang,
						Imp:  imp,
					},
					lang: lang,
					dep:  lbl,
				})
			}
		}
	}

	if len(overrides) == 0 {
		return
	}

	rewriteResolveConfigOverrides(getResolveConfig(c), overrides)
}

// ResolveImports is a utility function that returns a matching list of labels
// for the given import list.
func ResolveImports(resolver ImportResolver, lang, impLang string, imports []string) []label.Label {
	deps := make([]label.Label, 0)
	for _, imp := range imports {
		result := resolver.Resolve(lang, impLang, imp)
		if len(result) > 0 {
			first := result[0]
			deps = append(deps, first.Label)
			// if r.options.Debug {
			// 	r.options.Printf(lang, imp, "HIT", first.Label)
			// }
			// } else {
			// 	if r.options.Debug {
			// 		r.options.Printf(lang, imp, "MISS", resolver)
			// 	}
		}
	}
	return deps
}

// ResolveImportsString is a utility function that returns a matching list of labels
// for the given import list.
func ResolveImportsString(resolver ImportResolver, rel, kind, impLang string, imports []string) []string {
	match := ResolveImports(resolver, kind, impLang, imports)
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

func langKey(lang, impLang string) string {
	return lang + " " + impLang
}

func keyLang(key string) (string, string) {
	parts := strings.SplitN(key, " ", 2)
	return parts[0], parts[1]
}
