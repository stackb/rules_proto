package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestCsharpPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.CsharpPlugin{}, map[string]plugintest.Case{
		// --csharp_out always generates output files
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "csharp implementation builtin:csharp",
			),
			PluginName: "csharp",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("Test.cs"),
			),
		},
		// per-object generated files
		"single enum message service": {
			Input: "enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "csharp implementation builtin:csharp",
			),
			PluginName: "csharp",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("Test.cs"),
			),
		},
		// package does not affect output location
		"package enum message service": {
			Input: "package p; enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "csharp implementation builtin:csharp",
			),
			PluginName: "csharp",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("Test.cs"),
			),
		},
		"csharp_namespace does not on affect output location": {
			Input: "package p; option csharp_namespace=\"Aa.Bb.Cc\"; enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "csharp implementation builtin:csharp",
			),
			PluginName: "csharp",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("Test.cs"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "csharp implementation builtin:csharp",
			),
			PluginName: "csharp",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/Test.cs"),
			),
		},
		"basename converted to pascal": {
			Basename: "foo_bar-baz",
			Input:    "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "csharp implementation builtin:csharp",
			),
			PluginName: "csharp",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("FooBarBaz.cs"),
			),
		},
	})
}
