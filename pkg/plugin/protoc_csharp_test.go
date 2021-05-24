package plugin

import (
	"testing"
)

func TestProtocCsharpPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocCsharpPlugin{}, map[string]PluginTestCase{
		// --csharp_out always generates output files
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "csharp implementation protoc:csharp",
			),
			Configuration: WithConfiguration(
				WithName("csharp"),
				WithOutputs("Test.cs"),
			),
		},
		// per-object generated files
		"single enum message service": {
			Input: "enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "csharp implementation protoc:csharp",
			),
			Configuration: WithConfiguration(
				WithName("csharp"),
				WithOutputs("Test.cs"),
			),
		},
		// package does not affect output location
		"package enum message service": {
			Input: "package p; enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "csharp implementation protoc:csharp",
			),
			Configuration: WithConfiguration(
				WithName("csharp"),
				WithOutputs("Test.cs"),
			),
		},
		"csharp_namespace does not on affect output location": {
			Input: "package p; option csharp_namespace=\"Aa.Bb.Cc\"; enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "csharp implementation protoc:csharp",
			),
			Configuration: WithConfiguration(
				WithName("csharp"),
				WithOutputs("Test.cs"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "csharp implementation protoc:csharp",
			),
			Configuration: WithConfiguration(
				WithName("csharp"),
				WithOutputs("rel/Test.cs"),
			),
		},
		"basename converted to pascal": {
			Basename: "foo_bar-baz",
			Input:    "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "csharp implementation protoc:csharp",
			),
			Configuration: WithConfiguration(
				WithName("csharp"),
				WithOutputs("FooBarBaz.cs"),
			),
		},
	})
}
