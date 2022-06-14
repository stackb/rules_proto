package ts_proto_test

import (
	"testing"

	ts_proto "github.com/stackb/rules_proto/pkg/plugin/stephenh/ts-proto"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestProtocGenTsProtoPlugin(t *testing.T) {
	plugintest.Cases(t, &ts_proto.ProtocGenTsProto{}, map[string]plugintest.Case{
		"simple": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-ts-proto implementation stephenh:ts-proto:protoc-gen-ts-proto",
			),
			PluginName: "protoc-gen-ts-proto",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test.ts"),
			),
			SkipIntegration: true,
		},
		"flag --exclude_output": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-ts-proto implementation stephenh:ts-proto:protoc-gen-ts-proto",
				"proto_plugin", "protoc-gen-ts-proto flag --exclude_output=test.ts",
			),
			PluginName:    "protoc-gen-ts-proto",
			Configuration: plugintest.WithConfiguration(
			// no outputs
			),
			SkipIntegration: true,
		},
	})
}
