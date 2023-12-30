package grpcgateway

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugintest"
)

const serviceEmpty = `
service S{}
`

const serviceWithUnboundMethods = `
service S{rpc R (I) returns (O) {}}
message I{}
message O{}
`

const serviceWithBoundMethods = `
import "google/api/annotations.proto";
service S{rpc R (I) returns (O) {
  option (google.api.http) = {
    get: "/path"
   };
}}
message I{}
message O{}
`

func TestProtocGenGoPlugin(t *testing.T) {
	plugintest.Cases(t, &protocGenGrpcGatewayPlugin{}, map[string]plugintest.Case{
		"no service": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway",
			),
			PluginName:      "protoc-gen-grpc-gateway",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"empty service": {
			Input: serviceEmpty,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway",
			),
			PluginName:      "protoc-gen-grpc-gateway",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"empty service and generate_unbound_methods": {
			Input: serviceEmpty,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway",
				"proto_plugin", "protoc-gen-grpc-gateway option generate_unbound_methods=true",
			),
			PluginName:      "protoc-gen-grpc-gateway",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"service with unbound methods": {
			Input: serviceWithUnboundMethods,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway",
			),
			PluginName:      "protoc-gen-grpc-gateway",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"service with unbound methods and generate_unbound_methods": {
			Input: serviceWithUnboundMethods,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway",
				"proto_plugin", "protoc-gen-grpc-gateway option generate_unbound_methods=true",
			),
			PluginName: "protoc-gen-grpc-gateway",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/grpc-ecosystem/grpc-gateway:protoc-gen-grpc-gateway"),
				plugintest.WithOptions("generate_unbound_methods=true"),
				plugintest.WithOutputs("test.pb.gw.go"),
			),
			SkipIntegration: true,
		},
		"service with bound methods": {
			Input: serviceWithBoundMethods,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway",
			),
			PluginName: "protoc-gen-grpc-gateway",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/grpc-ecosystem/grpc-gateway:protoc-gen-grpc-gateway"),
				plugintest.WithOutputs("test.pb.gw.go"),
			),
			SkipIntegration: true,
		},
		"service with bound methods and generate_unbound_methods": {
			Input: serviceWithBoundMethods,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway",
				"proto_plugin", "protoc-gen-grpc-gateway option generate_unbound_methods=true",
			),
			PluginName: "protoc-gen-grpc-gateway",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/grpc-ecosystem/grpc-gateway:protoc-gen-grpc-gateway"),
				plugintest.WithOptions("generate_unbound_methods=true"),
				plugintest.WithOutputs("test.pb.gw.go"),
			),
			SkipIntegration: true,
		},
	})
}
