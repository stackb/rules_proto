package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestJsCommonPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.JsCommonPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.js"),
				plugintest.WithOptions("import_style=commonjs"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.js"),
				plugintest.WithOptions("import_style=commonjs"),
			),
		},
		"single message & enum": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.js"),
				plugintest.WithOptions("import_style=commonjs"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.js"),
				plugintest.WithOptions("import_style=commonjs"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/test_pb.js"),
				plugintest.WithOptions("import_style=commonjs"),
			),
		},
	})
}
