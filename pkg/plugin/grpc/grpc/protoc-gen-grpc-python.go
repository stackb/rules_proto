package grpc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenGrpcPython{})
}

// ProtocGenGrpcPython implements Plugin for protoc-gen-grpc-python in the
// grpc/grpc repo.
type ProtocGenGrpcPython struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenGrpcPython) Name() string {
	return "grpc:grpc:protoc-gen-grpc-python"
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenGrpcPython) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !protoc.HasServices(ctx.ProtoLibrary.Files()...) {
		return nil
	}
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/grpc/grpc", "protoc-gen-grpc-python"),
		Outputs: protoc.FlatMapFiles(
			grpcGeneratedFileName(ctx.Rel),
			protoc.HasService,
			ctx.ProtoLibrary.Files()...,
		),
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// grpcGeneratedFileName is a utility function that returns a function that
// computes the name of a predicted generated file having the given extension(s)
// relative to the given dir.
func grpcGeneratedFileName(reldir string) func(f *protoc.File) []string {
	return func(f *protoc.File) []string {
		name := strings.ReplaceAll(f.Name, "-", "_")
		if reldir != "" {
			name = path.Join(reldir, name)
		}
		return []string{name + "_pb2_grpc.py"}
	}
}
