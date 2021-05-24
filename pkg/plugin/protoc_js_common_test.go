package plugin

import "testing"

func TestProtocJsCommonPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocJsCommonPlugin{}, map[string]PluginTestCase{
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:common",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("test_pb.js"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:common",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("test_pb.js"),
			),
		},
		"single message & enum": {
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:common",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("test_pb.js"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:common",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("test_pb.js"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "js implementation protoc:js:common",
			),
			Configuration: WithConfiguration(
				WithName("js"),
				WithOutputs("rel/test_pb.js"),
			),
		},
	})
}
