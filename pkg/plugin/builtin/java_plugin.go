package builtin

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&JavaPlugin{})
}

// JavaPlugin implements Plugin for the built-in protoc java plugin.
type JavaPlugin struct{}

// Name implements part of the Plugin interface.
func (p *JavaPlugin) Name() string {
	return "builtin:java"
}

// Configure implements part of the Plugin interface.
func (p *JavaPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	srcjar := path.Join(ctx.Rel, ctx.ProtoLibrary.BaseName()+".srcjar")
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/builtin", "java"),
		Outputs: []string{srcjar},
		Out:     srcjar,
		Options: ctx.PluginConfig.GetOptions(),
	}
}
