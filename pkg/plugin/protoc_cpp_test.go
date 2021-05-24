package plugin

import "testing"

func TestProtocCppPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocCppPlugin{}, map[string]PluginTestCase{
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "cpp implementation protoc:cpp",
			),
			Configuration: WithConfiguration(
				WithName("cpp"),
				WithOutputs("test.pb.cc", "test.pb.h"),
			),
		},
		"only services": {
			Input: "service S{}",
			Directives: WithDirectives(
				"proto_plugin", "cpp implementation protoc:cpp",
			),
			Configuration: WithConfiguration(
				WithName("cpp"),
				WithOutputs("test.pb.cc", "test.pb.h"),
			),
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "cpp implementation protoc:cpp",
			),
			Configuration: WithConfiguration(
				WithName("cpp"),
				WithOutputs("test.pb.cc", "test.pb.h"),
			),
		},
		"in a relative directory": {
			Rel:   "rel",
			Input: "package a;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "cpp implementation protoc:cpp",
			),
			Configuration: WithConfiguration(
				WithName("cpp"),
				WithOutputs("rel/test.pb.cc", "rel/test.pb.h"),
			),
		},
		"snake_case": {
			Basename: "snake_case",
			Input:    "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "cpp implementation protoc:cpp",
			),
			Configuration: WithConfiguration(
				WithName("cpp"),
				WithOutputs("snake_case.pb.cc", "snake_case.pb.h"),
			),
		},
		"PascalCase": {
			Basename: "PascalCase",
			Input:    "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "cpp implementation protoc:cpp",
			),
			Configuration: WithConfiguration(
				WithName("cpp"),
				WithOutputs("PascalCase.pb.cc", "PascalCase.pb.h"),
			),
		},
		"camelCase": {
			Basename: "camelCase",
			Input:    "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "cpp implementation protoc:cpp",
			),
			Configuration: WithConfiguration(
				WithName("cpp"),
				WithOutputs("camelCase.pb.cc", "camelCase.pb.h"),
			),
		},
	})
}
