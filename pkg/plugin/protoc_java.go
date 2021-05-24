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
	cfg.Skip = !p.shouldApply(ctx.ProtoLibrary)
	if cfg.Skip {
		return
	}

	cfg.Label = label.New("build_stack_rules_proto", "plugin/protoc", "java")
	cfg.Outputs = p.outputs(ctx.Rel, ctx.ProtoLibrary)
	cfg.Out = p.out(ctx.Rel, ctx.ProtoLibrary)
}

func (p *ProtocJavaPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

func (p *ProtocJavaPlugin) outputs(rel string, lib protoc.ProtoLibrary) []string {
	return []string{srcjarFile(rel, lib.BaseName())}
}

func (p *ProtocJavaPlugin) out(rel string, lib protoc.ProtoLibrary) string {
	return srcjarFile(rel, lib.BaseName())
}

func srcjarFile(dir, name string) string {
	return path.Join(dir, name+".srcjar")
}
