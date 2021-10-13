package builtin

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&RubyPlugin{})
}

// RubyPlugin implements Plugin for the built-in protoc ruby plugin.
type RubyPlugin struct{}

// Name implements part of the Plugin interface.
func (p *RubyPlugin) Name() string {
	return "builtin:ruby"
}

// Configure implements part of the Plugin interface.
func (p *RubyPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/builtin", "ruby"),
		Outputs: protoc.FlatMapFiles(
			protoc.RelativeFileNameWithExtensions(ctx.Rel, "_pb.rb"),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
		Options: ctx.PluginConfig.GetOptions(),
	}
}
