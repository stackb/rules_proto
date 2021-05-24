package plugin

import "testing"

func TestProtocPythonPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocPythonPlugin{}, map[string]PluginTestCase{
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "python implementation protoc:python",
			),
			Configuration: WithConfiguration(
				WithName("python"),
				WithOutputs("test_pb2.py"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: WithDirectives(
				"proto_plugin", "python implementation protoc:python",
			),
			Configuration: WithConfiguration(
				WithName("python"),
				WithOutputs("test_pb2.py"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "python implementation protoc:python",
			),
			Configuration: WithConfiguration(
				WithName("python"),
				WithOutputs("test_pb2.py"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "package a;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "python implementation protoc:python",
			),
			Configuration: WithConfiguration(
				WithName("python"),
				WithOutputs("rel/test_pb2.py"),
			),
		},
		"basename replacement": {
			Basename: "a-b*c+d=e|g!h#i",
			Input:    "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "python implementation protoc:python",
			),
			Configuration: WithConfiguration(
				WithName("python"),
				WithOutputs("a_b*c+d=e|g!h#i_pb2.py"),
			),
		},
	})
}
