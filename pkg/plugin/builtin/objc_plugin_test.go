package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestObjcPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.ObjcPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "objc implementation builtin:objc",
			),
			PluginName: "objc",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("Test.pbobjc.h", "Test.pbobjc.m"),
			),
		},
		"single enum message service": {
			Input: "enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "objc implementation builtin:objc",
			),
			PluginName: "objc",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("Test.pbobjc.h", "Test.pbobjc.m"),
			),
		},
		"package does not affect output location": {
			Input: "package p; enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "objc implementation builtin:objc",
			),
			PluginName: "objc",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("Test.pbobjc.h", "Test.pbobjc.m"),
			),
		},
		"objc_class_prefix does not affect output location": {
			Input: "package p; option objc_class_prefix=\"CGOOP\"; enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "objc implementation builtin:objc",
			),
			PluginName: "objc",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("Test.pbobjc.h", "Test.pbobjc.m"),
			),
		},
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "objc implementation builtin:objc",
			),
			PluginName: "objc",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/Test.pbobjc.h", "rel/Test.pbobjc.m"),
			),
		},
		"basename converted to capitalized camel case": {
			Basename: "foo_bar-baz",
			Input:    "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "objc implementation builtin:objc",
			),
			PluginName: "objc",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("FooBarBaz.pbobjc.h", "FooBarBaz.pbobjc.m"),
			),
		},
	})
}
