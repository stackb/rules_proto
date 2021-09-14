package builtin

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoGrpcCcLibraryRuleName   = "proto_grpc_cc_library"
	ProtoGrpcCcLibraryRuleSuffix = "_grpc_cc_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_grpc_cc_library", &protoGrpcCcLibrary{})
}

// protoGrpcCcLibrary implements LanguageRule for the 'proto_grpc_cc_library' rule from
// @rules_proto.
type protoGrpcCcLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoGrpcCcLibrary) Name() string {
	return ProtoGrpcCcLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoGrpcCcLibrary) KindInfo() rule.KindInfo {
	return ccLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoGrpcCcLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/cc:proto_grpc_cc_library.bzl",
		Symbols: []string{ProtoGrpcCcLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoGrpcCcLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("grpc:grpc:cpp")
	if len(outputs) == 0 {
		return nil
	}

	return &CcLibraryRule{
		KindName:       ProtoGrpcCcLibraryRuleName,
		RuleNameSuffix: ProtoGrpcCcLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver: func(impl DepsProvider, pc *protoc.ProtocConfiguration, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
			deps := impl.Deps()
			deps = append(deps, ":"+pc.Library.BaseName()+ProtoCcLibraryRuleSuffix)

			if len(deps) > 0 {
				r.SetAttr("deps", deps)
			}
		},
	}
}
