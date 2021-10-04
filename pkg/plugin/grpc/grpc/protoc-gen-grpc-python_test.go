package grpc_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/grpc/grpc"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestProtocGenGrpcPython(t *testing.T) {
	plugintest.Cases(t, &grpc.ProtocGenGrpcPython{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation grpc:grpc:protoc-gen-grpc-python",
			),
			PluginName:      "python",
			SkipIntegration: true,
		},
		"only messages": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation grpc:grpc:protoc-gen-grpc-python",
			),
			PluginName:      "python",
			SkipIntegration: true,
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation grpc:grpc:protoc-gen-grpc-python",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb2_grpc.py"),
			),
			PluginName:      "python",
			SkipIntegration: true,
		},
		"with a package": {
			Input: "package pkg;\n\nservice S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation grpc:grpc:protoc-gen-grpc-python",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb2_grpc.py"),
			),
			PluginName:      "python",
			SkipIntegration: true,
		},
	})
}
