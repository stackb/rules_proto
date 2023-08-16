package bufbuild_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/bufbuild"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestProtocGenTsProtoPlugin(t *testing.T) {
	plugintest.Cases(t, &bufbuild.ConnectProto{}, map[string]plugintest.Case{
		"simple": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "es implementation bufbuild:connect-es",
			),
			PluginName: "es",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/bufbuild:es"),
				plugintest.WithOutputs("test.pb.ts"),
			),
			SkipIntegration: true,
		},
		"includes only relevant M options": {
			Input: `
syntax = "proto3";

package corp.common;

import "google/type/datetime.proto";
import "google/protobuf/duration.proto";

message M {}
`,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "es implementation bufbuild:connect-es",
				"proto_plugin", "es option Mgoogle/protobuf/empty.proto=./external/protobufapis/google/protobuf/empty",
				"proto_plugin", "es option Mgoogle/protobuf/timestamp.proto=./external/protobufapis/google/protobuf/timestamp",
				"proto_plugin", "es option Mgoogle/protobuf/duration.proto=./external/protobufapis/google/protobuf/duration",
				"proto_plugin", "es option Mgoogle/type/timeofday.proto=./external/googleapis/google/type/timeofday",
				"proto_plugin", "es option Mgoogle/type/datetime.proto=./external/googleapis/google/type/datetime",
			),
			PluginName: "es",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/bufbuild:es"),
				plugintest.WithOutputs("test.pb.ts"),
				plugintest.WithOptions(
					"Mgoogle/protobuf/duration.proto=./external/protobufapis/google/protobuf/duration",
					"Mgoogle/type/datetime.proto=./external/googleapis/google/type/datetime",
				),
			),
			SkipIntegration: true,
		},
	})
}
