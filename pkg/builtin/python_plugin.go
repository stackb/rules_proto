package builtin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&PythonPlugin{})
}

// PythonPlugin implements Plugin for the built-in protoc python plugin.
type PythonPlugin struct{}

// Name implements part of the Plugin interface.
func (p *PythonPlugin) Name() string {
	return "protoc:python"
}

// Configure implements part of the Plugin interface.
func (p *PythonPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	cfg.Label = label.New("build_stack_rules_proto", "plugin/protoc", "python")
	cfg.Outputs = protoc.FlatMapFiles(
		pythonGeneratedFileName(ctx.Rel),
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)
}

// pythonGeneratedFileName is a utility function that returns a fucntion that
// compuutes the name of a predicted generated file having the given
// extension(s) relative to the given dir.
func pythonGeneratedFileName(reldir string) func(f *protoc.File) []string {
	return func(f *protoc.File) []string {
		name := strings.ReplaceAll(f.Name, "-", "_")
		if reldir != "" {
			name = path.Join(reldir, name)
		}
		return []string{name + "_pb2.py"}
	}
}
