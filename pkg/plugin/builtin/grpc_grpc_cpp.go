package builtin

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&GrpcGrpcCppPlugin{})
}

// GrpcGrpcCppPlugin implements Plugin for the built-in protoc python plugin.
type GrpcGrpcCppPlugin struct{}

// Name implements part of the Plugin interface.
func (p *GrpcGrpcCppPlugin) Name() string {
	return "grpc:grpc:cpp"
}

// Configure implements part of the Plugin interface.
func (p *GrpcGrpcCppPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !protoc.HasServices(ctx.ProtoLibrary.Files()...) {
		return nil
	}

	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/grpc/grpc", "protoc-gen-grpc-cpp"),
		Outputs: protoc.FlatMapFiles(
			protoc.ImportPrefixRelativeFileNameWithExtensions(ctx.ProtoLibrary.StripImportPrefix(), ctx.Rel, ".grpc.pb.cc", ".grpc.pb.h"),
			protoc.HasService,
			ctx.ProtoLibrary.Files()...,
		),
		Options: ctx.PluginConfig.GetOptions(),
	}
}
