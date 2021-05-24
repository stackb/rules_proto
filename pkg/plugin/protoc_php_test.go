package plugin

import (
	"testing"
)

func TestProtocPhpPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocPhpPlugin{}, map[string]PluginTestCase{
		// --php_out always generates output files
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
		// per-object generated files
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
		// package affects output location
		"package enum message service": {
			Input: "package p; enum E{U=0;} message M{} service S{}",
			Directives: WithDirectives(
				"proto_plugin", "php implementation protoc:php",
			),
			Configuration: WithConfiguration(
				WithName("php"),
				WithOutputs("GPBMetadata/Test.php", "P/E.php", "P/M.php"),
			),
		},
		// php_namespace
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
		// php_metadata_namespace
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
		// // reldir influences output location via --php_out=REL.  However, since
		// // we are expecting a relative output location by default (otherwise PluginConfiguration.Mappings would be populated)
		// "relative directory": {
		// 	Rel:   "rel",
		// 	Input: "message M{}",
		// 	Directives: WithDirectives(
		// 		"proto_plugin", "php implementation protoc:php",
		// 	),
		// 	Configuration: WithConfiguration(
		// 		WithName("php"),
		// 		WithOutputs("rel/test_pb.php"),
		// 	),
		// },
	})
}
