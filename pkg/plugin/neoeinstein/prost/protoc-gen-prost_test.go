package prost_test

import (
	"testing"

	"github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/prost"
	"github.com/stackb/rules_proto/v4/pkg/plugintest"
)

func TestProtocGenProstPlugin(t *testing.T) {
	plugintest.Cases(t, &prost.ProtocGenProstPlugin{}, map[string]plugintest.Case{
		"empty - no messages or enums": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost implementation neoeinstein:prost:protoc-gen-prost",
			),
			PluginName:      "protoc-gen-prost",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"simple message": {
			Input: "package example.v1;\nmessage Foo {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost implementation neoeinstein:prost:protoc-gen-prost",
			),
			PluginName: "protoc-gen-prost",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/prost:protoc-gen-prost"),
				plugintest.WithOutputs("example.v1.rs"),
			),
			SkipIntegration: true,
		},
		"simple enum": {
			Input: "package example.v1;\nenum Color { RED = 0; }",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost implementation neoeinstein:prost:protoc-gen-prost",
			),
			PluginName: "protoc-gen-prost",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/prost:protoc-gen-prost"),
				plugintest.WithOutputs("example.v1.rs"),
			),
			SkipIntegration: true,
		},
		"service only - emits stub .rs for tonic to insert into": {
			// Without this case the protoc plugin chain breaks: tonic
			// appends its client/server code at the @@protoc_insertion_point
			// in {package}.rs, and protoc errors with
			// "Tried to insert into file that doesn't exist" if prost
			// hasn't created that file first.
			Input: "package example.v1;\nservice Greeter { rpc Hello(Req) returns (Resp); }\nmessage Req {}\nmessage Resp {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost implementation neoeinstein:prost:protoc-gen-prost",
			),
			PluginName: "protoc-gen-prost",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/prost:protoc-gen-prost"),
				plugintest.WithOutputs("example.v1.rs"),
			),
			SkipIntegration: true,
		},
		"service only - no messages": {
			// Same as above but the proto has no messages/enums of its own
			// (could be e.g. a separate proto_library that imports its
			// request/response types from a sibling). Prost still has to
			// emit the .rs stub so tonic's append succeeds.
			Input: "package example.v1;\nimport \"google/protobuf/empty.proto\";\nservice Pinger { rpc Ping(google.protobuf.Empty) returns (google.protobuf.Empty); }",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost implementation neoeinstein:prost:protoc-gen-prost",
			),
			PluginName: "protoc-gen-prost",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/prost:protoc-gen-prost"),
				plugintest.WithOutputs("example.v1.rs"),
			),
			SkipIntegration: true,
		},
		"no package - skipped": {
			Input: "message Foo {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost implementation neoeinstein:prost:protoc-gen-prost",
			),
			PluginName:      "protoc-gen-prost",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"with options": {
			Input: "package example.v1;\nmessage Foo {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost implementation neoeinstein:prost:protoc-gen-prost",
				"proto_plugin", "protoc-gen-prost option type_attribute=.example.v1.Foo=#[derive(Eq)]",
			),
			PluginName: "protoc-gen-prost",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/prost:protoc-gen-prost"),
				plugintest.WithOutputs("example.v1.rs"),
				plugintest.WithOptions("type_attribute=.example.v1.Foo=#[derive(Eq)]"),
			),
			SkipIntegration: true,
		},
	})
}
