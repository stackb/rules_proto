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
	return "builtin:python"
}

// Configure implements part of the Plugin interface.
func (p *PythonPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/builtin", "python"),
		Outputs: protoc.FlatMapFiles(
			pythonGeneratedFileName(ctx.Rel),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
		Imports: protoc.FlatMapFiles(
			pyImports(),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
	}
}

// pythonGeneratedFileName is a utility function that returns a function that
// computes the name of a predicted generated file having the given
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

// pyImports is a utility function that returns a function that
// computes the name of a predicted imports for a given proto file.
func pyImports() func(f *protoc.File) []string {
	return func(f *protoc.File) []string {
		imports := make([]string, 0)

		pkg := f.Name + "_pb"
		if f.Package().Name != "" {
			pkg = f.Package().Name + "." + pkg
		}
		for _, m := range f.Messages() {
			imports = append(imports, getPythonImportName(pkg, m.Name))
		}
		for _, e := range f.Enums() {
			imports = append(imports, getPythonImportName(pkg, e.Name))
		}

		return imports
	}
}

func getPythonImportName(pkg, name string) string {
	if pkg == "" {
		return name
	}
	return pkg + "." + name
}
