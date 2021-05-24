package plugin

import (
	"testing"
)

func TestProtocJavaPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocJavaPlugin{}, map[string]PluginTestCase{
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "java implementation protoc:java",
			),
			Configuration: WithConfiguration(
				WithName("java"),
				WithOutputs("test.srcjar"),
			),
		},
		// package statement does not affect srcjar location
		"message with a package": {
			Input: "package a;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "java implementation protoc:java",
			),
			Configuration: WithConfiguration(
				WithName("java"),
				WithSkip(false),
				WithOutputs("test.srcjar"),
			),
		},
		// Rel affects srcjar location
		"relative package location": {
			Rel:   "src/main/java/foo",
			Input: "package a;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "java implementation protoc:java",
			),
			Configuration: WithConfiguration(
				WithName("java"),
				WithSkip(false),
				WithOutputs("src/main/java/foo/test.srcjar"),
			),
		},
	})
}
