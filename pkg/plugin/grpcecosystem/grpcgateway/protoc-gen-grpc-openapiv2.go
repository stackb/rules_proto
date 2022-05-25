package grpcgateway

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&protocGenGrpcOpenapiv2Plugin{})
}

// protocGenGrpcOpenapiv2Plugin implements Plugin for protoc-gen-grpc-openapiv2.
type protocGenGrpcOpenapiv2Plugin struct{}

// Name implements part of the Plugin interface.
func (p *protocGenGrpcOpenapiv2Plugin) Name() string {
	return "grpc-ecosystem:grpc-gateway:protoc-gen-grpc-openapiv2"
}

// Configure implements part of the Plugin interface.
func (p *protocGenGrpcOpenapiv2Plugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/grpc-ecosystem/grpc-gateway", "protoc-gen-grpc-openapiv2"),
		Outputs: p.outputs(ctx.Rel, ctx.ProtoLibrary),
		Options: ctx.PluginConfig.GetOptions(),
	}
}

func (p *protocGenGrpcOpenapiv2Plugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

func (p *protocGenGrpcOpenapiv2Plugin) outputs(rel string, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !f.HasServices() {
			continue
		}
		base := f.Name
		if rel != "" {
			base = path.Join(rel, base)
		}
		srcs = append(srcs, base+".swagger.json")
	}
	return srcs
}
