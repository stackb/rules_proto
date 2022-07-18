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
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto"),
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
			PluginName: "protoc-gen-ts-proto",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto"),
				plugintest.WithOutputs(),
			),
			SkipIntegration: true,
		},
		"includes only relevant M= options": {
			Input: `
syntax = "proto3";

package corp.common;

import "google/type/datetime.proto";
import "google/protobuf/duration.proto";

message M {}
`,
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-ts-proto implementation stephenh:ts-proto:protoc-gen-ts-proto",
				"proto_plugin", "protoc-gen-ts-proto option M=google/protobuf/empty.proto=./external/protobufapis/google/protobuf/empty",
				"proto_plugin", "protoc-gen-ts-proto option M=google/protobuf/timestamp.proto=./external/protobufapis/google/protobuf/timestamp",
				"proto_plugin", "protoc-gen-ts-proto option M=google/protobuf/duration.proto=./external/protobufapis/google/protobuf/duration",
				"proto_plugin", "protoc-gen-ts-proto option M=google/type/timeofday.proto=./external/googleapis/google/type/timeofday",
				"proto_plugin", "protoc-gen-ts-proto option M=google/type/datetime.proto=./external/googleapis/google/type/datetime",
			),
			PluginName: "protoc-gen-ts-proto",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto"),
				plugintest.WithOutputs("test.ts"),
				plugintest.WithOptions(
					"M=google/protobuf/duration.proto=./external/protobufapis/google/protobuf/duration",
					"M=google/type/datetime.proto=./external/googleapis/google/type/datetime",
				),
			),
			SkipIntegration: true,
		},
	})
}
