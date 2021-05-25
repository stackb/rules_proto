package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestJsCommonPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.JsCommonPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("test_pb.js"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("test_pb.js"),
			),
		},
		"single message & enum": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("test_pb.js"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("test_pb.js"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "js implementation builtin:js:common",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("js"),
				plugintest.WithOutputs("rel/test_pb.js"),
			),
		},
	})
}
