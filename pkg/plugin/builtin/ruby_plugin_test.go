package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestRubyPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.RubyPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "ruby implementation builtin:ruby",
			),
			PluginName: "ruby",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.rb"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "ruby implementation builtin:ruby",
			),
			PluginName: "ruby",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.rb"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "ruby implementation builtin:ruby",
			),
			PluginName: "ruby",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_pb.rb"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "ruby implementation builtin:ruby",
			),
			PluginName: "ruby",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/test_pb.rb"),
			),
		},
	})
}
