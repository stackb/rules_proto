package builtin

import (
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoJavaLibraryRuleName   = "proto_java_library"
	ProtoJavaLibraryRuleSuffix = "_java_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_java_library", &protoJavaLibrary{})
}

// protoJavaLibrary implements LanguageRule for the 'proto_java_library' rule from
// @rules_proto.
type protoJavaLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoJavaLibrary) Name() string {
	return ProtoJavaLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoJavaLibrary) KindInfo() rule.KindInfo {
	return ccLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoJavaLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/java:proto_java_library.bzl",
		Symbols: []string{ProtoJavaLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoJavaLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("builtin:java")
	if len(outputs) == 0 {
		return nil
	}
	return &JavaLibraryRule{
		KindName:       ProtoJavaLibraryRuleName,
		RuleNameSuffix: ProtoJavaLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver:       ResolveWithSuffix(ProtoJavaLibraryRuleSuffix),
	}
}
