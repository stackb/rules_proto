package plugin

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocCppPlugin{})
}

// ProtocCppPlugin implements Plugin for the built-in protoc C++ plugin.
type ProtocCppPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocCppPlugin) Name() string {
	return "protoc:cpp"
}

// Configure implements part of the Plugin interface.
func (p *ProtocCppPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	cfg.Skip = false
	cfg.Label = label.New("build_stack_rules_proto", "plugin/protoc", "cpp")
	cfg.Outputs = protoc.FlatMapFiles(
		protoc.PackageFileNameWithExtensions(".pb.cc", ".pb.h"),
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)
}
