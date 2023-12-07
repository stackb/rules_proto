package grpcgateway

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugintest"
)

const serviceWithoutBindings = `
service S{rpc R (I) returns (O) {}}
message I{}
message O{}
`

const serviceWithBindings = `
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
		"service without bindings": {
			Input: serviceWithoutBindings,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-grpc-gateway implementation grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway",
			),
			PluginName:      "protoc-gen-grpc-gateway",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"service without bindings and generate_unbound_methods": {
			Input: serviceWithoutBindings,
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
		"service with bindings": {
			Input: serviceWithBindings,
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
		"service with bindings and generate_unbound_methods": {
			Input: serviceWithBindings,
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
