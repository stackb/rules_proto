package rules_closure

import (
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoClosurejsLibraryRuleName   = "proto_closure_js_library"
	ProtoClosureJsLibraryRuleSuffix = "_closure_js_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_closure_js_library", &protoClosureJsLibrary{})
}

// protoClosureJsLibrary implements LanguageRule for the
// 'proto_closure_js_library' rule.
type protoClosureJsLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoClosureJsLibrary) Name() string {
	return ProtoClosurejsLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoClosureJsLibrary) KindInfo() rule.KindInfo {
	return closureJsLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoClosureJsLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/closure:proto_closure_js_library.bzl",
		Symbols: []string{ProtoClosurejsLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoClosureJsLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("builtin:js:closure")
	if len(outputs) == 0 {
		return nil
	}
	return &ClosureJsLibrary{
		KindName:       ProtoClosurejsLibraryRuleName,
		RuleNameSuffix: ProtoClosureJsLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
	}
}
