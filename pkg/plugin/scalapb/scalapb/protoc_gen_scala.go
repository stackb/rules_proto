package scalapb

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const ScalaPBPluginName = "scalapb:scalapb:protoc-gen-scala"

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenScalaPlugin{})
}

// ProtocGenScalaPlugin implements Plugin for the scala plugin.
type ProtocGenScalaPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenScalaPlugin) Name() string {
	return ScalaPBPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenScalaPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	options := ctx.PluginConfig.GetOptions()

	srcjar := ctx.ProtoLibrary.BaseName() + "_scala"
	if containsOption(options, "grpc") {
		srcjar += "_grpc"
	}
	srcjar += ".srcjar"
	if ctx.Rel != "" {
		srcjar = path.Join(ctx.Rel, srcjar)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/scalapb/scalapb", "protoc-gen-scala"),
		Outputs: []string{srcjar},
		Options: options,
	}
}

// containsOption returns a copy of the slice will all elements matching 'key'
// removed.  If at least item was removed, return true.
func containsOption(src []string, key string) bool {
	for _, item := range src {
		if item == key {
			return true
		}
	}
	return false
}
