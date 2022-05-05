package protoc

import (
	"fmt"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func init() {
	Rules().MustRegisterRule("stackb:rules_proto:golden_filegroup", &goldenFilegroupRule{})
}

// goldenFilegroupRule implements LanguageRule for the 'golden_filegroup' rule.
type goldenFilegroupRule struct{}

// Name implements part of the LanguageRule interface.
func (s *goldenFilegroupRule) Name() string {
	return "rules_golden_filegroup"
}

// KindInfo implements part of the LanguageRule interface.
func (s *goldenFilegroupRule) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *goldenFilegroupRule) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules:golden_filegroup.bzl",
		Symbols: []string{"golden_filegroup"},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *goldenFilegroupRule) ProvideRule(cfg *LanguageRuleConfig, config *ProtocConfiguration) RuleProvider {
	return &goldenFilegroupRuleRule{ruleConfig: cfg, config: config}
}

// goldenFilegroupRule implements RuleProvider for the 'proto_compile' rule.
type goldenFilegroupRuleRule struct {
	config     *ProtocConfiguration
	ruleConfig *LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *goldenFilegroupRuleRule) Kind() string {
	return "golden_filegroup"
}

// Name implements part of the ruleProvider interface.
func (s *goldenFilegroupRuleRule) Name() string {
	return fmt.Sprintf("%s_goldens", s.config.Library.BaseName())
}

// Visibility provides visibility labels.
func (s *goldenFilegroupRuleRule) Visibility() []string {
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
func (s *goldenFilegroupRuleRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	srcs := make([]string, 0)
	for _, file := range s.config.Library.Files() {
		srcs = append(srcs, file.Name)
	}
	newRule.SetAttr("srcs", srcs)

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *goldenFilegroupRuleRule) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *goldenFilegroupRuleRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
}
