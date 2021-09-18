package protoc

import (
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type DepsProvider interface {
	Deps() []string
}

type DepsResolver func(impl DepsProvider, pc *ProtocConfiguration, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label)

func ResolveDepsWithSuffix(suffix string) DepsResolver {
	return func(impl DepsProvider, pc *ProtocConfiguration, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
		deps := impl.Deps()

		for _, d := range pc.Library.Deps() {
			if strings.HasPrefix(d, "@com_google_protobuf//") {
				continue
			}
			d = strings.TrimSuffix(d, "_proto")
			deps = append(deps, d+suffix)
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
	return filename[len(rel)+1:] // +1 for slash separator
}
