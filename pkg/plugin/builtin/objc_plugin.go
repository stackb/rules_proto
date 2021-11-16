package builtin

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ObjcPlugin{})
}

// ObjcPlugin implements Plugin for the built-in protoc C# plugin.
type ObjcPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ObjcPlugin) Name() string {
	return "builtin:objc"
}

// Configure implements part of the Plugin interface.
func (p *ObjcPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/builtin", "objc"),
		Outputs: protoc.FlatMapFiles(
			objcFileName(ctx.Rel, ctx.PluginConfig),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
		Options: ctx.PluginConfig.GetOptions(),
	}
}

func objcFileName(rel string, cfg protoc.LanguagePluginConfig) func(*protoc.File) []string {
	return func(f *protoc.File) []string {
		// setup the file extension
		base := path.Join(rel, protoc.ToPascalCase(f.Name))
		return []string{base + ".pbobjc.h", base + ".pbobjc.m"}
	}

}
