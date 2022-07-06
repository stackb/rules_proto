package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestJavaPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.JavaPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "java implementation builtin:java",
			),
			PluginName: "java",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/builtin:java"),
				plugintest.WithOutputs("test.srcjar"),
				plugintest.WithOut("test.srcjar"),
			),
		},
		"message with a package": {
			Input: "package a;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "java implementation builtin:java",
			),
			PluginName: "java",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/builtin:java"),
				plugintest.WithOutputs("test.srcjar"),
				plugintest.WithOut("test.srcjar"),
			),
		},
		"relative package location": {
			Rel:   "src/main/java/foo",
			Input: "package a;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "java implementation builtin:java",
			),
			PluginName: "java",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/builtin:java"),
				plugintest.WithOutputs("src/main/java/foo/test.srcjar"),
				plugintest.WithOut("src/main/java/foo/test.srcjar"),
			),
		},
	})
}
