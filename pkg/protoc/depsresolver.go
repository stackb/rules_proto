package protoc

import (
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type DepsProvider interface {
	Deps() []string
}

type DepsResolver func(impl DepsProvider, pc *ProtocConfiguration, c *PackageConfig, r *rule.Rule, imports []string, from label.Label)

func ResolveDepsExcludingWellKnownTypes() DepsResolver {
	return func(impl DepsProvider, pc *ProtocConfiguration, c *PackageConfig, r *rule.Rule, imports []string, from label.Label) {
		deps := impl.Deps()

		for _, imp := range imports {
			if strings.HasPrefix(imp, "google/protobuf/") {
				continue
			}
			result := c.Requires(r.Kind(), imp)
			if len(result) > 0 {
				first := result[0]
				deps = append(deps, first.Label.String())
			}
		}

		if len(deps) > 0 {
			r.SetAttr("deps", deps)
		}
	}
}

// StripRel removes the rel prefix from a filename (if has matching prefix)
func StripRel(rel string, filename string) string {
	if !strings.HasPrefix(filename, rel) {
		return filename
	}
	filename = filename[len(rel):]
	return strings.TrimPrefix(filename, "/")
}
