package scalapb_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/scalapb/scalapb"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestProtoGenScalaPlugin(t *testing.T) {
	plugintest.Cases(t, &scalapb.ProtocGenScalaPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "scala implementation scalapb:scalapb:protoc-gen-scala",
			),
			PluginName: "scala",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_scala.srcjar"),
			),
			SkipIntegration: true,
		},
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "scala implementation scalapb:scalapb:protoc-gen-scala",
			),
			PluginName: "scala",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("test_scala.srcjar"),
			),
			SkipIntegration: true,
		},
		"in a relative directory": {
			Rel:   "rel",
			Input: "package a;\n\nmessage M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "scala implementation scalapb:scalapb:protoc-gen-scala",
			),
			PluginName: "scala",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("rel/test_scala.srcjar"),
			),
			SkipIntegration: true,
		},
	})
}
