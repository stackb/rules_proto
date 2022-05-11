package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestPyiPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.PyiPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "pyi implementation builtin:pyi",
			),
			PluginName: "pyi",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb2.pyi"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "pyi implementation builtin:pyi",
			),
			PluginName: "pyi",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb2.pyi"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "pyi implementation builtin:pyi",
			),
			PluginName: "pyi",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb2.pyi"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "package a;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "pyi implementation builtin:pyi",
			),
			PluginName: "pyi",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/test_pb2.pyi"),
			),
		},
		"basename replacement": {
			Basename: "a-b*c+d=e|g!h#i",
			Input:    "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "pyi implementation builtin:pyi",
			),
			PluginName: "pyi",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("a_b*c+d=e|g!h#i_pb2.pyi"),
			),
		},
	})
}
