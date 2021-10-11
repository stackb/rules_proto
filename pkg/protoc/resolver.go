package protoc

import (
	"log"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
)

const debugResolver = false

// ImportResolver handles dependency resolution.
type ImportResolver interface {
	// Requires returns any previously provided labels associated with the given
	// kind and import.
	Requires(kind string, imp string) []resolve.FindResult
	// Provides records the association between a kind of rule, a proto import
	// statement, and the label that provides the import.
	Provides(kind string, imp string, location label.Label)
}

// Resolver returns a reference to the global ImportResolver
func Resolver() ImportResolver {
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
func ResolveImportsString(resolver ImportResolver, kind string, imports []string) []string {
	match := ResolveImports(resolver, kind, imports)
	deps := make([]string, len(match))
	for i, l := range match {
		deps[i] = l.String()
	}
	return deps
}
