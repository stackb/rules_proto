package protobuf_javascript_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/protocolbuffers/protobuf-javascript"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestJsES6Plugin(t *testing.T) {
	plugintest.Cases(t, &protobuf_javascript.JsES6Plugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-js implementation protocolbuffers:protobuf-javascript:protoc-gen-js",
			),
			PluginName: "protoc-gen-js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.js"),
			),
			SkipIntegration: true,
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-js implementation protocolbuffers:protobuf-javascript:protoc-gen-js",
			),
			PluginName: "protoc-gen-js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.js"),
			),
			SkipIntegration: true,
		},
		"single message & enum": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-js implementation protocolbuffers:protobuf-javascript:protoc-gen-js",
			),
			PluginName: "protoc-gen-js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.js"),
			),
			SkipIntegration: true,
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-js implementation protocolbuffers:protobuf-javascript:protoc-gen-js",
			),
			PluginName: "protoc-gen-js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.js"),
			),
			SkipIntegration: true,
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-js implementation protocolbuffers:protobuf-javascript:protoc-gen-js",
			),
			PluginName: "protoc-gen-js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/test_pb.js"),
			),
			SkipIntegration: true,
		},
	})
}
