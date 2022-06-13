package rules_nodejs

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	grpcWebJsLibraryRuleName   = "grpc_web_js_library"
	grpcWebJsLibraryRuleSuffix = "_grpc_web_js_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:grpc_web_js_library", &grpcWebJsLibrary{})
}

// grpcWebJsLibrary implements LanguageRule for the 'grpc_web_js_library' rule from @build_stack_rules_proto.
// (which is essentially a wrapper for 'js_library' rule from @build_bazel_rules_nodejs)
type grpcWebJsLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *grpcWebJsLibrary) Name() string {
	return grpcWebJsLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *grpcWebJsLibrary) KindInfo() rule.KindInfo {
	return jsLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *grpcWebJsLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/nodejs:grpc_web_js_library.bzl",
		Symbols: []string{grpcWebJsLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *grpcWebJsLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("grpc:grpc-web:protoc-gen-grpc-web")
	if len(outputs) == 0 {
		return nil
	}

	return &JsLibrary{
		KindName:       grpcWebJsLibraryRuleName,
		RuleNameSuffix: grpcWebJsLibraryRuleSuffix,
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
