package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestJsClosurePlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.JsClosurePlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("test.js"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("test.js"),
			),
		},
		"single message & enum": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("test.js"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("test.js"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:closure",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("rel/test.js"),
			),
		},
	})
}
