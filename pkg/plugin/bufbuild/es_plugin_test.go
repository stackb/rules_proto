package bufbuild_test

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/plugin/bufbuild"
	"github.com/stackb/rules_proto/pkg/plugintest"
)

func TestProtocGenTsProtoPlugin(t *testing.T) {
	plugintest.Cases(t, &bufbuild.EsProto{}, map[string]plugintest.Case{
		"simple": {
			Input: "message M{}",
			Directives: plugintest.WithDirectives(
				"proto_plugin", "es implementation bufbuild:connect-es",
			),
			PluginName: "es",
			Configuration: plugintest.WithConfiguration(
				plugintest.WithLabel(t, "@build_stack_rules_proto//plugin/bufbuild:es"),
				plugintest.WithOptions("keep_empty_files=true", "target=ts"),
				plugintest.WithOutputs("test_pb.ts"),
			),
			SkipIntegration: true,
		},
	})
}
