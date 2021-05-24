package plugin

import (
	"testing"
)

func TestProtocCppPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocCppPlugin{}, map[string]PluginTestCase{
		// --cpp_out always generates output files
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
		// it does not matter if it only has services
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
		// package statement does not affect output location
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
		// reldir influences output location via --cpp_out=REL.  However, since
		// we are expecting a relative output location by default (otherwise PluginConfiguration.Mappings would be populated)
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
