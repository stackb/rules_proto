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
				WithSkip(false),
				WithOutputs("test.pb.cc", "test.pb.h"),
			),
		},
		// "only services": {
		// 	Input: "service S{}",
		// 	Directives: WithDirectives(
		// 		"proto_plugin", "cpp implementation protoc:cpp",
		// 	),
		// 	Configuration: WithConfiguration(
		// 		WithName("cpp"),
		// 		WithSkip(true),
		// 	),
		// },
		// "message with no package": {
		// 	Input: "message M{}",
		// 	Directives: WithDirectives(
		// 		"proto_plugin", "cpp implementation protoc:cpp",
		// 	),
		// 	Configuration: WithConfiguration(
		// 		WithName("cpp"),
		// 		WithSkip(false),
		// 		WithOutputs("test.pb.cc", "test.pb.h"),
		// 	),
		// },
		// "message with a package": {
		// 	Input: "package a;\n\nmessage M{}",
		// 	Directives: WithDirectives(
		// 		"proto_plugin", "cpp implementation protoc:cpp",
		// 	),
		// 	Configuration: WithConfiguration(
		// 		WithName("cpp"),
		// 		WithSkip(false),
		// 		WithOutputs("a/test.pb.cc", "a/test.pb.h"),
		// 	),
		// },
	})
}
