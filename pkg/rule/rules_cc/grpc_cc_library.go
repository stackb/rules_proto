package rules_cc

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	grpcCcLibraryRuleName   = "grpc_cc_library"
	grpcCcLibraryRuleSuffix = "_grpc_cc_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:grpc_cc_library", &grpcCcLibrary{})
}

// grpcCcLibrary implements LanguageRule for the 'grpc_cc_library' rule from
// @rules_proto.
type grpcCcLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *grpcCcLibrary) Name() string {
	return grpcCcLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *grpcCcLibrary) KindInfo() rule.KindInfo {
	return ccLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *grpcCcLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/cc:grpc_cc_library.bzl",
		Symbols: []string{grpcCcLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *grpcCcLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("grpc:grpc:cpp")
	if len(outputs) == 0 {
		return nil
	}

	return &CcLibrary{
		KindName:       grpcCcLibraryRuleName,
		RuleNameSuffix: grpcCcLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver: func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
			deps := append(r.AttrStrings("deps"), ":"+pc.Library.BaseName()+ProtoCcLibraryRuleSuffix)

			if len(deps) > 0 {
				r.SetAttr("deps", deps)
			}
		},
	}
}
