package akka_grpc

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const AkkaGrpcPluginName = "akka:akka-grpc:protoc-gen-akka-grpc"

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenAkkaGrpcPlugin{})
}

// ProtocGenAkkaGrpcPlugin implements Plugin for the akka-grpc plugin.
type ProtocGenAkkaGrpcPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenAkkaGrpcPlugin) Name() string {
	return AkkaGrpcPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenAkkaGrpcPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !protoc.HasServices(ctx.ProtoLibrary.Files()...) {
		return nil
	}

	srcjar := ctx.ProtoLibrary.BaseName() + "_akka_grpc.srcjar"
	if ctx.Rel != "" {
		srcjar = path.Join(ctx.Rel, srcjar)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/akka/akka-grpc", "protoc-gen-akka-grpc"),
		Outputs: []string{srcjar},
		Options: ctx.PluginConfig.GetOptions(),
	}
}
