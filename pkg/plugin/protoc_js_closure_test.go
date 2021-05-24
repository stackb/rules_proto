package plugin

import "testing"

func TestProtocJsClosurePlugin(t *testing.T) {
	PluginTestCases(t, &ProtocJsClosurePlugin{}, map[string]PluginTestCase{
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:closure",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("test.js"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:closure",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("test.js"),
			),
		},
		"single message & enum": {
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:closure",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("test.js"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:closure",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("test.js"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:closure",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("rel/test.js"),
			),
		},
	})
}
