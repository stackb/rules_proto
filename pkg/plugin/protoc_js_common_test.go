package plugin

import (
	"testing"
)

func TestProtocJsCommonPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocJsCommonPlugin{}, map[string]PluginTestCase{
		// --js_out with an empty proto generates an interesting result!
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
		// package statement does not affect output location
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
		// reldir affects output location via --js_out=REL.
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
