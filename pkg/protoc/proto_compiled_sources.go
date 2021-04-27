package protoc

import (
	"fmt"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// ProtoCompiledSourcesName the name the protoCompiledSourcesdSources rule registers under.
const ProtoCompiledSourcesName = "proto_compiled_sources"

func init() {
	Rules().MustRegisterRule(ProtoCompiledSourcesName, &protoCompiledSources{})
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
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoCompiledSources) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules:proto_compiled_sources.bzl",
		Symbols: []string{ProtoCompiledSourcesName},
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
	return ProtoCompiledSourcesName
}

// Name implements part of the ruleProvider interface.
func (s *protoCompiledSourcesRule) Name() string {
	return fmt.Sprintf("%s_%s_compiled_sources", s.config.Library.BaseName(), s.config.Prefix)
}

// Imports implements part of the ruleProvider interface.
func (s *protoCompiledSourcesRule) Imports() []string {
	return []string{s.Kind()}
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

	newRule.SetAttr("srcs", s.config.Outputs)
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

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *protoCompiledSourcesRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
}
