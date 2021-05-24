package plugin

import (
	"testing"
)

func TestProtocObjcPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocObjcPlugin{}, map[string]PluginTestCase{
		// --objc_out always generates output files
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "objc implementation protoc:objc",
			),
			Configuration: WithConfiguration(
				WithName("objc"),
				WithOutputs("Test.pbobjc.h", "Test.pbobjc.m"),
			),
		},
		// per-object generated files
		"single enum message service": {
			Input: "enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "objc implementation protoc:objc",
			),
			Configuration: WithConfiguration(
				WithName("objc"),
				WithOutputs("Test.pbobjc.h", "Test.pbobjc.m"),
			),
		},
		// package does not affect output location
		"package enum message service": {
			Input: "package p; enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "objc implementation protoc:objc",
			),
			Configuration: WithConfiguration(
				WithName("objc"),
				WithOutputs("Test.pbobjc.h", "Test.pbobjc.m"),
			),
		},
		"objc_class_prefix does not on affect output location": {
			Input: "package p; option objc_class_prefix=\"CGOOP\"; enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "objc implementation protoc:objc",
			),
			Configuration: WithConfiguration(
				WithName("objc"),
				WithOutputs("Test.pbobjc.h", "Test.pbobjc.m"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "objc implementation protoc:objc",
			),
			Configuration: WithConfiguration(
				WithName("objc"),
				WithOutputs("rel/Test.pbobjc.h", "rel/Test.pbobjc.m"),
			),
		},
		"basename converted to pascal": {
			Basename: "foo_bar-baz",
			Input:    "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "objc implementation protoc:objc",
			),
			Configuration: WithConfiguration(
				WithName("objc"),
				WithOutputs("FooBarBaz.pbobjc.h", "FooBarBaz.pbobjc.m"),
			),
		},
	})
}
