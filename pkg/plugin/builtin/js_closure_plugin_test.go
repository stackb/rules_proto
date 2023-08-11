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
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/builtin:closurejs"),
				plugintest.WithOutputs("test_closure.js"),
				plugintest.WithOptions("import_style=closure", "library=test_closure"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/builtin:closurejs"),
				plugintest.WithOutputs("test_closure.js"),
				plugintest.WithOptions("import_style=closure", "library=test_closure"),
			),
		},
		"single message & enum": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/builtin:closurejs"),
				plugintest.WithOutputs("test_closure.js"),
				plugintest.WithOptions("import_style=closure", "library=test_closure"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			PluginName: "js",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/builtin:closurejs"),
				plugintest.WithOutputs("test_closure.js"),
				plugintest.WithOptions("import_style=closure", "library=test_closure"),
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
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/builtin:closurejs"),
				plugintest.WithOutputs("rel/test_closure.js"),
				plugintest.WithOptions("import_style=closure", "library=rel/test_closure"),
			),
		},
	})
}
