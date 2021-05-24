package plugin

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocJavaPlugin{})
}

// ProtocJavaPlugin implements Plugin for the built-in protoc java plugin.
type ProtocJavaPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocJavaPlugin) Name() string {
	return "protoc:java"
}

// Configure implements part of the Plugin interface.
func (p *ProtocJavaPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	cfg.Label = label.New("build_stack_rules_proto", "plugin/protoc", "java")
	cfg.Outputs = []string{srcjarFile(ctx.Rel, ctx.ProtoLibrary.BaseName())}
	cfg.Out = srcjarFile(ctx.Rel, ctx.ProtoLibrary.BaseName())
}

func srcjarFile(dir, name string) string {
	return path.Join(dir, name+".srcjar")
}
