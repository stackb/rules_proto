package rust_protobuf

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&protocGenRust{})
}

// protocGenRust implements Plugin for the protoc-gen-rust plugin.
type protocGenRust struct{}

// Name implements part of the Plugin interface.
func (p *protocGenRust) Name() string {
	return "stepancheg:rust-protobuf:protoc-gen-rust"
}

// Configure implements part of the Plugin interface.
func (p *protocGenRust) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/stepancheg/rust-protobuf", "protoc-gen-rust"),
		Outputs: protoc.FlatMapFiles(
			protoc.ImportPrefixRelativeFileNameWithExtensions(ctx.ProtoLibrary.StripImportPrefix(), ctx.Rel, ".rs"),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
		Options: ctx.PluginConfig.GetOptions(),
		Out:     ctx.Rel,
	}
}
