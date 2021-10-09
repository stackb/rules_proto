package protobuf_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/golang/protobuf"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestProtocGenGoPlugin(t *testing.T) {
	plugintest.Cases(t, &protobuf.ProtocGenGoPlugin{}, map[string]plugintest.Case{
		"simple": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-go implementation golang:protobuf:protoc-gen-go",
			),
			PluginName: "protoc-gen-go",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test.pb.go"),
			),
			SkipIntegration: true,
		},
		"option go_package": {
			Input: "option go_package=\"github.com/example.com/foo\";\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-go implementation golang:protobuf:protoc-gen-go",
			),
			PluginName: "protoc-gen-go",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("github.com/example.com/foo/test.pb.go"),
			),
			SkipIntegration: true,
		},
		"impport mapping": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-go implementation golang:protobuf:protoc-gen-go",
				"proto_plugin", "protoc-gen-go option Mtest.proto=github.com/example.com/foo",
			),
			PluginName: "protoc-gen-go",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("github.com/example.com/foo/test.pb.go"),
			),
			SkipIntegration: true,
		},
	})
}
