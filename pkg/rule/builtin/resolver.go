package builtin

import (
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

type DepsProvider interface {
	Deps() []string
}

type DepsResolver func(impl DepsProvider, pc *protoc.ProtocConfiguration, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label)

func ResolveWithSuffix(suffix string) DepsResolver {
	return func(impl DepsProvider, pc *protoc.ProtocConfiguration, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
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
