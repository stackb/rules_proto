package plugin

import (
	"testing"
)

func TestProtocPythonPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocPythonPlugin{}, map[string]PluginTestCase{
		// --python_out always generates output files
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
		// it does not matter if it only has services
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
		// package statement does not affect output location
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
		// reldir affects output location via --python_out=REL.
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
