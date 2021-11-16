package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestJsClosurePlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.JsClosurePlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test.js"),
				plugintest.WithOptions("import_style=closure", "library=test"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test.js"),
				plugintest.WithOptions("import_style=closure", "library=test"),
			),
		},
		"single message & enum": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test.js"),
				plugintest.WithOptions("import_style=closure", "library=test"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test.js"),
				plugintest.WithOptions("import_style=closure", "library=test"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/test.js"),
				plugintest.WithOptions("import_style=closure", "library=rel/test"),
			),
		},
	})
}
