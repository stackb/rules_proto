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
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/scalapb/scalapb", "protoc-gen-scala"),
		Outputs: []string{path.Join(ctx.Rel, ctx.ProtoLibrary.BaseName()+"_scala.srcjar")},
	}
}
