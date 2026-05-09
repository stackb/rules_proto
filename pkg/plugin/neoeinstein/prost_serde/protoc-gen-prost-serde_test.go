package prost_serde_test

import (
	"testing"

	"github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/prost_serde"
	"github.com/stackb/rules_proto/v4/pkg/plugintest"
)

func TestProtocGenProstSerdePlugin(t *testing.T) {
	plugintest.Cases(t, &prost_serde.ProtocGenProstSerdePlugin{}, map[string]plugintest.Case{
		"empty - no messages or enums": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost-serde implementation neoeinstein:prost:protoc-gen-prost-serde",
			),
			PluginName:      "protoc-gen-prost-serde",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"simple message": {
			Input: "package example.v1;\nmessage Foo {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost-serde implementation neoeinstein:prost:protoc-gen-prost-serde",
			),
			PluginName: "protoc-gen-prost-serde",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/prost-serde:protoc-gen-prost-serde"),
				plugintest.WithOutputs("example.v1.serde.rs"),
			),
			SkipIntegration: true,
		},
		"simple enum": {
			Input: "package example.v1;\nenum Color { RED = 0; }",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost-serde implementation neoeinstein:prost:protoc-gen-prost-serde",
			),
			PluginName: "protoc-gen-prost-serde",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/prost-serde:protoc-gen-prost-serde"),
				plugintest.WithOutputs("example.v1.serde.rs"),
			),
			SkipIntegration: true,
		},
		"no package - skipped": {
			Input: "message Foo {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost-serde implementation neoeinstein:prost:protoc-gen-prost-serde",
			),
			PluginName:      "protoc-gen-prost-serde",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"with options": {
			Input: "package example.v1;\nmessage Foo {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-prost-serde implementation neoeinstein:prost:protoc-gen-prost-serde",
				"proto_plugin", "protoc-gen-prost-serde option compile_well_known_types=true",
			),
			PluginName: "protoc-gen-prost-serde",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/prost-serde:protoc-gen-prost-serde"),
				plugintest.WithOutputs("example.v1.serde.rs"),
				plugintest.WithOptions("compile_well_known_types=true"),
			),
			SkipIntegration: true,
		},
	})
}
