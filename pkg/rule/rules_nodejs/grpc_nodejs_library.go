package rules_nodejs

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	grpcNodeJsLibraryRuleName   = "grpc_nodejs_library"
	grpcNodeJsLibraryRuleSuffix = "_grpc_nodejs_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:grpc_nodejs_library", &grpcNodeJsLibrary{})
}

// grpcNodeJsLibrary implements LanguageRule for the 'grpc_nodejs_library' rule from
// @rules_proto.
type grpcNodeJsLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *grpcNodeJsLibrary) Name() string {
	return grpcNodeJsLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *grpcNodeJsLibrary) KindInfo() rule.KindInfo {
	return jsLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *grpcNodeJsLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/nodejs:grpc_nodejs_library.bzl",
		Symbols: []string{grpcNodeJsLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *grpcNodeJsLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("grpc:grpc-node:protoc-gen-grpc-node")
	if len(outputs) == 0 {
		return nil
	}

	return &JsLibrary{
		KindName:       grpcNodeJsLibraryRuleName,
		RuleNameSuffix: grpcNodeJsLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver: func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
			deps := append(r.AttrStrings("deps"), ":"+pc.Library.BaseName()+ProtoNodeJsLibraryRuleSuffix)

			if len(deps) > 0 {
				r.SetAttr("deps", deps)
			}
		},
	}
}
