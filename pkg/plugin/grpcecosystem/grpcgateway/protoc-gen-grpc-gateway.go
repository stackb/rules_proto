package grpcgateway

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&protocGenGrpcGatewayPlugin{})
}

// protocGenGrpcGatewayPlugin implements Plugin for protoc-gen-grpc-gateway.
type protocGenGrpcGatewayPlugin struct{}

// Name implements part of the Plugin interface.
func (p *protocGenGrpcGatewayPlugin) Name() string {
	return "grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway"
}

// Configure implements part of the Plugin interface.
func (p *protocGenGrpcGatewayPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	generateUnbound := ctx.PluginConfig.Options["generate_unbound_methods=true"]
	if !p.shouldApply(ctx.ProtoLibrary, generateUnbound) {
		return nil
	}
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/grpc-ecosystem/grpc-gateway", "protoc-gen-grpc-gateway"),
		Outputs: p.outputs(ctx.Rel, ctx.ProtoLibrary, generateUnbound),
		Options: ctx.PluginConfig.GetOptions(),
	}
}

func (p *protocGenGrpcGatewayPlugin) shouldApply(lib protoc.ProtoLibrary, generateUnbound bool) bool {
	for _, f := range lib.Files() {
		if p.shouldOutputForFile(f, generateUnbound) {
			return true
		}
	}
	return false
}

func (p *protocGenGrpcGatewayPlugin) outputs(rel string, lib protoc.ProtoLibrary, generateUnbound bool) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !p.shouldOutputForFile(f, generateUnbound) {
			continue
		}
		base := f.Name
		if rel != "" {
			base = path.Join(rel, base)
		}
		srcs = append(srcs, base+".pb.gw.go")
	}
	return srcs
}

func (p *protocGenGrpcGatewayPlugin) shouldOutputForFile(f *protoc.File, generateUnbound bool) bool {
	return f.HasServices() && (generateUnbound || f.HasRPCOption("(google.api.http)"))
}
