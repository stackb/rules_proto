package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestJavaPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.JavaPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "java implementation protoc:java",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("java"),
				plugintest.WithOutputs("test.srcjar"),
			),
		},
		"message with a package": {
			Input: "package a;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "java implementation protoc:java",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("java"),
				plugintest.WithSkip(false),
				plugintest.WithOutputs("test.srcjar"),
			),
		},
		"relative package location": {
			Rel:   "src/main/java/foo",
			Input: "package a;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "java implementation protoc:java",
			),
			Configuration: plugintest.WithConfiguration(
				plugintest.WithName("java"),
				plugintest.WithSkip(false),
				plugintest.WithOutputs("src/main/java/foo/test.srcjar"),
			),
		},
	})
}
