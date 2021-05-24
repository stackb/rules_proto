package plugin

import "testing"

func TestProtocPhpPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocPhpPlugin{}, map[string]PluginTestCase{
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "php implementation protoc:php",
			),
			Configuration: WithConfiguration(
				WithName("php"),
				WithOutputs("GPBMetadata/Test.php"),
			),
		},
		"single enum message service": {
			Input: "enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "php implementation protoc:php",
			),
			Configuration: WithConfiguration(
				WithName("php"),
				WithOutputs("GPBMetadata/Test.php", "E.php", "M.php"),
			),
		},
		"package": {
			Input: "package p; enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "php implementation protoc:php",
			),
			Configuration: WithConfiguration(
				WithName("php"),
				WithOutputs("GPBMetadata/Test.php", "P/E.php", "P/M.php"),
			),
		},
		"php_namespace": {
			Input: "package p; option php_namespace=\"foo\"; enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "php implementation protoc:php",
			),
			Configuration: WithConfiguration(
				WithName("php"),
				WithOutputs("GPBMetadata/Test.php", "foo/E.php", "foo/M.php"),
			),
		},
		"php_metadata_namespace": {
			Input: "package p; option php_metadata_namespace=\"bar\"; enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "php implementation protoc:php",
			),
			Configuration: WithConfiguration(
				WithName("php"),
				WithOutputs("bar/Test.php", "P/E.php", "P/M.php"),
			),
		},
		"relative directory": {
			Rel:   "a/b/c",
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "php implementation protoc:php",
			),
			Configuration: WithConfiguration(
				WithName("php"),
				WithOutputs("a/b/c/GPBMetadata/A/B/C/Test.php", "a/b/c/M.php"),
			),
		},
	})
}
