package ts_proto

import (
	"path"

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
	tsFiles := make([]string, 0)
	for _, file := range ctx.ProtoLibrary.Files() {
		tsFile := file.Name + ".ts"
		if ctx.Rel != "" {
			tsFile = path.Join(ctx.Rel, tsFile)
		}
		tsFiles = append(tsFiles, tsFile)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/stephenh/ts-proto", "protoc-gen-ts-proto"),
		Outputs: tsFiles,
		Options: ctx.PluginConfig.GetOptions(),
	}
}
