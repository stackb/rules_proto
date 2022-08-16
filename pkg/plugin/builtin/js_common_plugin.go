// for some bizarre reason, naming this file 'protoc_js.go' makes it be ignored
// by the compiler?
package builtin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&JsCommonPlugin{})
}

// JsCommonPlugin implements Plugin for the built-in protoc js/library plugin.
type JsCommonPlugin struct{}

// Name implements part of the Plugin interface.
func (p *JsCommonPlugin) Name() string {
	return "builtin:js:common"
}

// Configure implements part of the Plugin interface.
func (p *JsCommonPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	basename := strings.ToLower(ctx.ProtoLibrary.BaseName())
	library := basename + "_pb.js"
	if ctx.Rel != "" {
		library = path.Join(ctx.Rel, library)
	}

	jsOutputs := protoc.FlatMapFiles(
		jsGeneratedFileName(ctx.Rel),
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)

	outputs := []string{library}
	if true {
		outputs = jsOutputs
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/builtin", "commonjs"),
		Outputs: outputs,
		Options: append(ctx.PluginConfig.GetOptions(), "import_style=commonjs"),
	}
}

// jsGeneratedFileName is a utility function that returns a function that
// computes the name of a predicted generated file having the given extension(s)
// relative to the given dir.
func jsGeneratedFileName(reldir string) func(f *protoc.File) []string {
	return func(f *protoc.File) []string {
		name := strings.ReplaceAll(f.Name, "-", "_")
		if reldir != "" {
			name = path.Join(reldir, name)
		}
		return []string{name + "_pb.js"}
	}
}
