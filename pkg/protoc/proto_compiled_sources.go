package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func init() {
	Rules().MustRegisterRule("stackb:rules_proto:proto_compiled_sources", &protoCompiledSources{})
}

// protoCompiledSources implements LanguageRule for the 'proto_compiled_sources' rule.
type protoCompiledSources struct{}

// KindInfo implements part of the LanguageRule interface.
func (s *protoCompiledSources) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		NonEmptyAttrs: map[string]bool{
			"srcs": true,
		},
		MergeableAttrs: map[string]bool{
			"srcs":            true,
			"plugins":         true,
			"output_mappings": true,
			"options":         true,
		},
		SubstituteAttrs: map[string]bool{
			"out": true,
		},
	}
}

// Name implements part of the LanguageRule interface.
func (s *protoCompiledSources) Name() string {
	return "proto_compiled_sources"
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoCompiledSources) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules:proto_compiled_sources.bzl",
		Symbols: []string{"proto_compiled_sources"},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoCompiledSources) ProvideRule(cfg *LanguageRuleConfig, config *ProtocConfiguration) RuleProvider {
	return &protoCompileRule{
		kind:            "proto_compiled_sources",
		nameSuffix:      "compiled_sources",
		outputsAttrName: "srcs",
		config:          config,
		ruleConfig:      cfg,
	}
}
