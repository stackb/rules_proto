package plugin

import "testing"

func TestProtocRubyPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocRubyPlugin{}, map[string]PluginTestCase{
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
