package protoc

import (
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// ProtoCompileName the name the protoCompile rule registers under.
const ProtoCompileName = "proto_compile"

func init() {
	Rules().MustRegisterRule(ProtoCompileName, &protoCompile{})
}

// protoCompile implements LanguageRule for the 'proto_compile' rule.
type protoCompile struct{}

// KindInfo implements part of the LanguageRule interface.
func (s *protoCompile) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		NonEmptyAttrs: map[string]bool{
			"outputs": true,
			"srcs":    true,
		},
		MergeableAttrs: map[string]bool{
			"outputs": true,
			"srcs":    true,
			"plugins": true,
		},
		SubstituteAttrs: map[string]bool{
			"options":  true,
			"out":      true,
			"mappings": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoCompile) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules:proto_compile.bzl",
		Symbols: []string{ProtoCompileName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoCompile) ProvideRule(cfg *LanguageRuleConfig, config *ProtocConfiguration) RuleProvider {
	return &protoCompileRule{config}
}

// protoCompile implements RuleProvider for the 'proto_compile' rule.
type protoCompileRule struct {
	config *ProtocConfiguration
}

// Kind implements part of the ruleProvider interface.
func (s *protoCompileRule) Kind() string {
	return ProtoCompileName
}

// Name implements part of the ruleProvider interface.
func (s *protoCompileRule) Name() string {
	return fmt.Sprintf("%s_%s_compile", s.config.Library.BaseName(), s.config.Prefix)
}

// Imports implements part of the ruleProvider interface.
func (s *protoCompileRule) Imports() []string {
	return []string{s.Kind()}
}

// Visibility implements part of the ruleProvider interface.
func (s *protoCompileRule) Visibility() []string {
	return nil // TODO: visibility feature?
}

// Rule implements part of the ruleProvider interface.
func (s *protoCompileRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("outputs", s.config.Outputs)
	newRule.SetAttr("plugins", GetPluginLabels(s.config.Plugins))
	newRule.SetAttr("proto", s.config.Library.Name())

	if len(s.config.Mappings) > 0 {
		newRule.SetAttr("mappings", makeStringDict(s.config.Mappings))
	}

	options := GetPluginOptions(s.config.Plugins)
	if len(options) > 0 {
		newRule.SetAttr("options", makeStringListDict(options))
	}

	outs := GetPluginOuts(s.config.Plugins)
	if len(outs) > 0 {
		newRule.SetAttr("outs", makeStringDict(outs))
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *protoCompileRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
}
