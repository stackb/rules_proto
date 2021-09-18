package builtin_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/builtin"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestPhpPlugin(t *testing.T) {
	plugintest.Cases(t, &builtin.PhpPlugin{}, map[string]plugintest.Case{
		"empty file": {
			Input: "",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "php implementation builtin:php",
			),
			PluginName: "php",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("GPBMetadata/Test.php"),
			),
		},
		"single enum message service": {
			Input: "enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "php implementation builtin:php",
			),
			PluginName: "php",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("GPBMetadata/Test.php", "E.php", "M.php"),
			),
		},
		"package": {
			Input: "package p; enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "php implementation builtin:php",
			),
			PluginName: "php",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("GPBMetadata/Test.php", "P/E.php", "P/M.php"),
			),
		},
		"php_namespace": {
			Input: "package p; option php_namespace=\"foo\"; enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "php implementation builtin:php",
			),
			PluginName: "php",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("GPBMetadata/Test.php", "foo/E.php", "foo/M.php"),
			),
		},
		"php_metadata_namespace": {
			Input: "package p; option php_metadata_namespace=\"bar\"; enum E{U=0;} message M{} service S{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "php implementation builtin:php",
			),
			PluginName: "php",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("bar/Test.php", "P/E.php", "P/M.php"),
			),
		},
		"relative directory": {
			Rel:   "a/b/c",
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "php implementation builtin:php",
			),
			PluginName: "php",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithOutputs("a/b/c/GPBMetadata/A/B/C/Test.php", "a/b/c/M.php"),
			),
		},
	})
}
