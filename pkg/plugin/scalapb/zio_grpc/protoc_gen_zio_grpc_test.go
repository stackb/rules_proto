package zio_grpc_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/scalapb/zio_grpc"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestProtoGenZioGrpcPlugin(t *testing.T) {
	plugintest.Cases(t, &zio_grpc.ProtocGenZioGrpcPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "zio implementation scalapb:zio-grpc:protoc-gen-zio-grpc",
			),
			PluginName:      "zio",
			SkipIntegration: true,
		},
		"only messages, no srcjar produced": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "zio implementation scalapb:zio-grpc:protoc-gen-zio-grpc",
			),
			PluginName:      "zio",
			SkipIntegration: true,
		},
		"with service": {
			Input: "package pkg;\n\nservice S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "zio implementation scalapb:zio-grpc:protoc-gen-zio-grpc",
			),
			PluginName: "zio",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/scalapb/zio-grpc:protoc-gen-zio-grpc"),
				plugintest.WithOutputs("test_zio_grpc.srcjar"),
			),
			SkipIntegration: true,
		},
	})
}
