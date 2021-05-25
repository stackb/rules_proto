package builtin

import (
	"fmt"
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
func (p *JavaPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	srcjar := path.Join(ctx.Rel, fmt.Sprintf("%s.srcjar", ctx.ProtoLibrary.BaseName()))
	cfg.Label = label.New("build_stack_rules_proto", "plugin/builtin", "java")
	cfg.Outputs = []string{srcjar}
}

// func srcjarFile(dir, name string) string {
// 	return path.Join(dir, name+".srcjar")
// }
