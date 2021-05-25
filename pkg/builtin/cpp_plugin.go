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
func (p *CppPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	cfg.Label = label.New("build_stack_rules_proto", "plugin/builtin", "cpp")
	cfg.Outputs = protoc.FlatMapFiles(
		protoc.RelativeFileNameWithExtensions(ctx.Rel, ".pb.cc", ".pb.h"),
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)
}
