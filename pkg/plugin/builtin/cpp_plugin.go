package builtin

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&CppPlugin{})
}

// CppPlugin implements Plugin for the built-in protoc C++ plugin.
type CppPlugin struct{}

// Name implements part of the Plugin interface.
func (p *CppPlugin) Name() string {
	return "builtin:cpp"
}

// Configure implements part of the Plugin interface.
func (p *CppPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/builtin", "cpp"),
		Outputs: protoc.FlatMapFiles(
			protoc.ImportPrefixRelativeFileNameWithExtensions(ctx.ProtoLibrary.StripImportPrefix(), ctx.Rel, ".pb.cc", ".pb.h"),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
		Options: ctx.PluginConfig.GetOptions(),
	}
}
