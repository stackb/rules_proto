package plugin

import (
	"testing"
)

func TestProtocRubyPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocRubyPlugin{}, map[string]PluginTestCase{
		// --ruby_out always generates output files
		"empty file": {
			Input: "",
			Directives: WithDirectives(
				"proto_plugin", "ruby implementation protoc:ruby",
			),
			Configuration: WithConfiguration(
				WithName("ruby"),
				WithOutputs("test_pb.rb"),
			),
		},
		// it does not matter if it only has services
		"only services": {
			Input: "service S{}",
			Directives: WithDirectives(
				"proto_plugin", "ruby implementation protoc:ruby",
			),
			Configuration: WithConfiguration(
				WithName("ruby"),
				WithOutputs("test_pb.rb"),
			),
		},
		// package statement does not affect output location
		"with a package": {
			Input: "package pkg;\n\nmessage M{}",
			Directives: WithDirectives(
				"proto_plugin", "ruby implementation protoc:ruby",
			),
			Configuration: WithConfiguration(
				WithName("ruby"),
				WithOutputs("test_pb.rb"),
			),
		},
		// reldir influences output location via --ruby_out=REL.  However, since
		// we are expecting a relative output location by default (otherwise PluginConfiguration.Mappings would be populated)
		"relative directory": {
			Rel:   "rel",
			Input: "message M{}",
			Directives: WithDirectives(
				"proto_plugin", "ruby implementation protoc:ruby",
			),
			Configuration: WithConfiguration(
				WithName("ruby"),
				WithOutputs("rel/test_pb.rb"),
			),
		},
	})
}
