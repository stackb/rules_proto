package java

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenGrpcJavaPlugin{})
}

// ProtocGenGrpcJavaPlugin implements Plugin for the grpc java plugin.
type ProtocGenGrpcJavaPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenGrpcJavaPlugin) Name() string {
	return "grpc:grpc-java:protoc-gen-grpc-java"
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenGrpcJavaPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}
	srcjar := path.Join(ctx.Rel, ctx.ProtoLibrary.BaseName()+"_grpc.srcjar")
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/grpc/grpc-java", "protoc-gen-grpc-java"),
		Outputs: []string{srcjar},
		Out:     srcjar,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

func (p *ProtocGenGrpcJavaPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}
