package akka_grpc_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/akka/akka_grpc"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestProtoGenAkkaGrpcPlugin(t *testing.T) {
	plugintest.Cases(t, &akka_grpc.ProtocGenAkkaGrpcPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "akka implementation akka:akka-grpc:protoc-gen-akka-grpc",
			),
			PluginName:      "akka",
			SkipIntegration: true,
		},
		"only messages, no srcjar produced": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "akka implementation akka:akka-grpc:protoc-gen-akka-grpc",
			),
			PluginName:      "akka",
			SkipIntegration: true,
		},
		"with service": {
			Input: "package pkg;\n\nservice S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "akka implementation akka:akka-grpc:protoc-gen-akka-grpc",
			),
			PluginName: "akka",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_akka_grpc.srcjar"),
			),
			SkipIntegration: true,
		},
	})
}
