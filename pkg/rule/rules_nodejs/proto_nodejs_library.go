package rules_nodejs

import (
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoNodeJsLibraryRuleName   = "proto_nodejs_library"
	ProtoNodeJsLibraryRuleSuffix = "_nodejs_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_nodejs_library", &protoNodeJsLibrary{})
}

// protoNodeJsLibrary implements LanguageRule for the 'proto_nodejs_library' rule from
// @rules_proto.
type protoNodeJsLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoNodeJsLibrary) Name() string {
	return ProtoNodeJsLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoNodeJsLibrary) KindInfo() rule.KindInfo {
	return jsLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoNodeJsLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/nodejs:proto_nodejs_library.bzl",
		Symbols: []string{ProtoNodeJsLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoNodeJsLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("builtin:js:common")
	if len(outputs) == 0 {
		return nil
	}
	return &JsLibrary{
		KindName:       ProtoNodeJsLibraryRuleName,
		RuleNameSuffix: ProtoNodeJsLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver:       protoc.ResolveDepsAttr("deps", true),
	}
}
