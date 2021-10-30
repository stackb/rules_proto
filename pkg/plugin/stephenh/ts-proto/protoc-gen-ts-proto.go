package ts_proto

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenTsProto{})
}

// ProtocGenTsProto implements Plugin for the built-in protoc js/library plugin.
type ProtocGenTsProto struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenTsProto) Name() string {
	return "stephenh:ts-proto:protoc-gen-ts-proto"
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenTsProto) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	basename := strings.ToLower(ctx.ProtoLibrary.BaseName())
	tsFile := basename + ".ts"
	if ctx.Rel != "" {
		tsFile = path.Join(ctx.Rel, tsFile)
	}
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/stephenh/ts-proto", "protoc-gen-ts-proto"),
		Outputs: []string{tsFile},
		Options: ctx.PluginConfig.GetOptions(),
	}
}
