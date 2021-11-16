package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestPythonPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.PythonPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation builtin:python",
			),
			PluginName: "python",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb2.py"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation builtin:python",
			),
			PluginName: "python",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb2.py"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation builtin:python",
			),
			PluginName: "python",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb2.py"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "package a;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation builtin:python",
			),
			PluginName: "python",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/test_pb2.py"),
			),
		},
		"basename replacement": {
			Basename: "a-b*c+d=e|g!h#i",
			Input:    "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "python implementation builtin:python",
			),
			PluginName: "python",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("a_b*c+d=e|g!h#i_pb2.py"),
			),
		},
	})
}
