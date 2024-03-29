package builtin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&PyiPlugin{})
}

// PyiPlugin implements Plugin for the built-in protoc --pyi_out.
type PyiPlugin struct{}

// Name implements part of the Plugin interface.
func (p *PyiPlugin) Name() string {
	return "builtin:pyi"
}

// Configure implements part of the Plugin interface.
func (p *PyiPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/builtin", "pyi"),
		Outputs: protoc.FlatMapFiles(
			pyiGeneratedFileName(ctx.Rel),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// pyiGeneratedFileName is a utility function that returns a function that
// computes the name of a predicted generated file having the given
// extension(s) relative to the given dir.
func pyiGeneratedFileName(reldir string) func(f *protoc.File) []string {
	return func(f *protoc.File) []string {
		name := strings.ReplaceAll(f.Name, "-", "_")
		if reldir != "" {
			name = path.Join(reldir, name)
		}
		return []string{name + "_pb2.pyi"}
	}
}
