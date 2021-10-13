package grpcgo

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/plugin/golang/protobuf"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenGoGrpcPlugin{})
}

// ProtocGenGoGrpcPlugin implements Plugin for the the gogo_* family of plugins.
type ProtocGenGoGrpcPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenGoGrpcPlugin) Name() string {
	return "grpc:grpc-go:protoc-gen-go-grpc"
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenGoGrpcPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}
	options := ctx.PluginConfig.GetOptions()
	mappings := protobuf.GetImportMappings(options)
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/grpc/grpc-go", "protoc-gen-go-grpc"),
		Outputs: p.outputs(ctx.ProtoLibrary, mappings),
		Options: protobuf.FilterImportMappingOptions(options, mappings, ctx.ProtoLibrary.Imports()),
	}
}

func (p *ProtocGenGoGrpcPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

func (p *ProtocGenGoGrpcPlugin) outputs(lib protoc.ProtoLibrary, importMappings map[string]string) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !f.HasServices() {
			continue
		}
		srcs = append(srcs, protobuf.GetGoOutputBaseName(f, importMappings)+"_grpc.pb.go")
	}
	return srcs
}
