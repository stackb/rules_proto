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
				WithSkip(true),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: WithDirectives(
				"proto_plugin", "java implementation protoc:java",
			),
			Configuration: WithConfiguration(
				WithName("java"),
				WithSkip(true),
			),
		},
		"message with no package": {
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "java implementation protoc:java",
			),
			Configuration: WithConfiguration(
				WithName("java"),
				WithSkip(false),
				WithOutputs("test.srcjar"),
			),
		},
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
	})
}
