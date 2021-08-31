package protoc

import (
	"fmt"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
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
			"srcs":       true,
			"plugins":    true,
			"visibility": true,
		},
		SubstituteAttrs: map[string]bool{
			"options":  true,
			"out":      true,
			"mappings": true,
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
	return &protoCompiledSourcesRule{ruleConfig: cfg, config: config}
}

// protoCompiledSources implements RuleProvider for the 'proto_compile' rule.
type protoCompiledSourcesRule struct {
	config     *ProtocConfiguration
	ruleConfig *LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *protoCompiledSourcesRule) Kind() string {
	return "proto_compiled_sources"
}

// Name implements part of the ruleProvider interface.
func (s *protoCompiledSourcesRule) Name() string {
	return fmt.Sprintf("%s_%s_compiled_sources", s.config.Library.BaseName(), s.config.Prefix)
}

// Visibility implements part of the ruleProvider interface.
func (s *protoCompiledSourcesRule) Visibility() []string {
	visibility := make([]string, 0)
	for k, want := range s.ruleConfig.Visibility {
		if !want {
			continue
		}
		visibility = append(visibility, k)
	}
	sort.Strings(visibility)
	return visibility
}

// Rule implements part of the ruleProvider interface.
func (s *protoCompiledSourcesRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	outputs := s.config.Outputs
	sort.Strings(outputs)

	newRule.SetAttr("srcs", outputs)
	newRule.SetAttr("plugins", GetPluginLabels(s.config.Plugins))
	newRule.SetAttr("proto", s.config.Library.Name())

	if s.config.LanguageConfig.Protoc != "" {
		newRule.SetAttr("protoc", s.config.LanguageConfig.Protoc)
	}

	if len(s.config.Mappings) > 0 {
		newRule.SetAttr("mappings", MakeStringDict(s.config.Mappings))
	}

	options := GetPluginOptions(s.config.Plugins)
	if len(options) > 0 {
		newRule.SetAttr("options", MakeStringListDict(options))
	}

	outs := GetPluginOuts(s.config.Plugins)
	if len(outs) > 0 {
		newRule.SetAttr("outs", MakeStringDict(outs))
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *protoCompiledSourcesRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
}
