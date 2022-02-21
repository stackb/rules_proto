package akka_grpc

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const AkkaGrpcPluginName = "akka:akka-grpc:protoc-gen-akka-grpc"

func init() {
	protoc.Plugins().MustRegisterPlugin(&protocGenAkkaGrpcPlugin{})
}

// protocGenAkkaGrpcPlugin implements Plugin for the akka-grpc plugin.
type protocGenAkkaGrpcPlugin struct{}

// Name implements part of the Plugin interface.
func (p *protocGenAkkaGrpcPlugin) Name() string {
	return AkkaGrpcPluginName
}

// Configure implements part of the Plugin interface.
func (p *protocGenAkkaGrpcPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !protoc.HasServices(ctx.ProtoLibrary.Files()...) {
		return nil
	}

	srcjar := ctx.ProtoLibrary.BaseName() + "_akka.srcjar"
	if ctx.Rel != "" {
		srcjar = path.Join(ctx.Rel, srcjar)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/akka/akka-grpc", "protoc-gen-akka-grpc"),
		Outputs: []string{srcjar},
		Options: ctx.PluginConfig.GetOptions(),
	}
}
