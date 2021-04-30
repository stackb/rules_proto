package grpc_gateway

import (
	"path"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const GrpcGatewayPluginName = "grpc_gateway"

func init() {
	protoc.Plugins().MustRegisterPlugin(GrpcGatewayPluginName, &GrpcGatewayPlugin{})
}

// GrpcGatewayPlugin implements Plugin for protoc-gen-grpc-gateway.
type GrpcGatewayPlugin struct{}

// Label implements part of the Plugin interface.
func (p *GrpcGatewayPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "grpc-ecosystem/grpc-gateway", "grpc_gateway_plugin")
}

func (p *GrpcGatewayPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface
func (p *GrpcGatewayPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !f.HasServices() {
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

// Options implements part of the optional PluginOptionsProvider interface.  If
// the library contains services, apply the grpc plugin.
func (p *GrpcGatewayPlugin) Options(rel string, c protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	return nil
}
