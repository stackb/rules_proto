package zio_grpc

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const ZioGrpcPluginName = "scalapb:zio-grpc:protoc-gen-zio-grpc"

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenZioGrpcPlugin{})
}

// ProtocGenZioGrpcPlugin implements Plugin for the zio-grpc plugin.
type ProtocGenZioGrpcPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenZioGrpcPlugin) Name() string {
	return ZioGrpcPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenZioGrpcPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !protoc.HasServices(ctx.ProtoLibrary.Files()...) {
		return nil
	}

	srcjar := ctx.ProtoLibrary.BaseName() + "_zio_grpc.srcjar"
	if ctx.Rel != "" {
		srcjar = path.Join(ctx.Rel, srcjar)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/scalapb/zio-grpc", "protoc-gen-zio-grpc"),
		Outputs: []string{srcjar},
		Options: ctx.PluginConfig.GetOptions(),
	}
}
