package tonic_test

import (
	"testing"

	"github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/tonic"
	"github.com/stackb/rules_proto/v4/pkg/plugintest"
)

func TestProtocGenTonicPlugin(t *testing.T) {
	plugintest.Cases(t, &tonic.ProtocGenTonicPlugin{}, map[string]plugintest.Case{
		"empty - no services": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-tonic implementation neoeinstein:prost:protoc-gen-tonic",
			),
			PluginName:      "protoc-gen-tonic",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"only messages - no output": {
			Input: "package example.v1;\nmessage Foo {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-tonic implementation neoeinstein:prost:protoc-gen-tonic",
			),
			PluginName:      "protoc-gen-tonic",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"only enums - no output": {
			Input: "package example.v1;\nenum Color { RED = 0; }",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-tonic implementation neoeinstein:prost:protoc-gen-tonic",
			),
			PluginName:      "protoc-gen-tonic",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"simple service": {
			Input: "package example.v1;\nservice Greeter {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-tonic implementation neoeinstein:prost:protoc-gen-tonic",
			),
			PluginName: "protoc-gen-tonic",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/tonic:protoc-gen-tonic"),
				plugintest.WithOutputs("example.v1.tonic.rs"),
			),
			SkipIntegration: true,
		},
		"no package - skipped": {
			Input: "service Greeter {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-tonic implementation neoeinstein:prost:protoc-gen-tonic",
			),
			PluginName:      "protoc-gen-tonic",
			Configuration:   nil,
			SkipIntegration: true,
		},
		"with options": {
			Input: "package example.v1;\nservice Greeter {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-tonic implementation neoeinstein:prost:protoc-gen-tonic",
				"proto_plugin", "protoc-gen-tonic option compile_well_known_types=true",
			),
			PluginName: "protoc-gen-tonic",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/tonic:protoc-gen-tonic"),
				plugintest.WithOutputs("example.v1.tonic.rs"),
				plugintest.WithOptions("compile_well_known_types=true"),
			),
			SkipIntegration: true,
		},
		"message and service": {
			Input: "package example.v1;\nmessage Foo {}\nservice Greeter {}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "protoc-gen-tonic implementation neoeinstein:prost:protoc-gen-tonic",
			),
			PluginName: "protoc-gen-tonic",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/neoeinstein/tonic:protoc-gen-tonic"),
				plugintest.WithOutputs("example.v1.tonic.rs"),
			),
			SkipIntegration: true,
		},
	})
}
