package protoc

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func init() {
	Rules().MustRegisterRule("stackb:rules_proto:fileset:proto_library", &filesetRule{
		collectLibraryRule: true,
	})
	Rules().MustRegisterRule("stackb:rules_proto:fileset:package", &filesetRule{
		collectPackageRules: true,
	})
	Rules().MustRegisterRule("stackb:rules_proto:fileset:all", &filesetRule{
		collectAllRules: true,
	})
}

// allFilesetRules represents all fileset rules created in all packages.
var allFilesetRules = make([]label.Label, 0)

// filesetRule implements LanguageRule for the 'fileset' rule.
type filesetRule struct {
	collectLibraryRule  bool
	collectPackageRules bool
	collectAllRules     bool
}

// Name implements part of the LanguageRule interface.
func (s *filesetRule) Name() string {
	return "fileset"
}

// KindInfo implements part of the LanguageRule interface.
func (s *filesetRule) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *filesetRule) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules:fileset.bzl",
		Symbols: []string{"fileset"},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *filesetRule) ProvideRule(cfg *LanguageRuleConfig, config *ProtocConfiguration) RuleProvider {
	// collectAllRules is only effective for the root package.
	collectAllRules := s.collectAllRules && config.Rel == ""

	return &filesetRuleRule{
		ruleConfig:          cfg,
		config:              config,
		collectLibraryRule:  s.collectLibraryRule,
		collectPackageRules: s.collectPackageRules,
		collectAllRules:     collectAllRules,
	}
}

// filesetRule implements RuleProvider for the 'proto_compile' rule.
type filesetRuleRule struct {
	// file is the rule.File captured during the Imports function.
	file       *rule.File
	config     *ProtocConfiguration
	ruleConfig *LanguageRuleConfig

	// collectLibrary is true if the rule should the parent proto_library rule.
	collectLibraryRule bool
	// collectPackageRules is true if the rule should collect package rules.
	collectPackageRules bool
	// collectAllRules is true if the rule should collect All rules.
	collectAllRules bool
}

// Kind implements part of the ruleProvider interface.
func (s *filesetRuleRule) Kind() string {
	return "fileset"
}

// Name implements part of the ruleProvider interface.
func (s *filesetRuleRule) Name() string {
	base := s.config.Library.BaseName()
	if s.collectPackageRules {
		base = filepath.Base(s.config.Rel)
	}
	if s.collectAllRules {
		base = "all"
	}
	return fmt.Sprintf("%s_files", base)
}

// Visibility provides visibility labels.
func (s *filesetRuleRule) Visibility() []string {
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
func (s *filesetRuleRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	srcs := make([]string, 0)
	for _, file := range s.config.Library.Files() {
		srcs = append(srcs, file.Basename)
	}
	newRule.SetAttr("srcs", srcs)

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	// populate the allFilesetRules with this label, unless it represents the root fileset itself.
	if !s.collectAllRules {
		allFilesetRules = append(allFilesetRules, label.New("", s.config.Rel, newRule.Name()))
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *filesetRuleRule) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	s.file = file
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *filesetRuleRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {

	srcs := r.AttrStrings("srcs")

	if s.collectLibraryRule {
		srcs = append(srcs, s.config.Library.Rule().Name())
	}

	if s.collectPackageRules && s.file != nil {
		for _, other := range s.file.Rules {
			if other.Name() == r.Name() {
				continue
			}
			switch other.Kind() {
			case "proto_compile":
				srcs = append(srcs, other.AttrStrings("outputs")...)
			case "proto_descriptor_set":
				srcs = append(srcs, other.Name())
			case "proto_library":
				srcs = append(srcs, other.AttrStrings("srcs")...)
				srcs = append(srcs, other.Name())
			}
		}
	}

	if s.collectAllRules {
		for _, other := range allFilesetRules {
			srcs = append(srcs, other.String())
		}
	}

	r.SetAttr("srcs", DeduplicateAndSort(srcs))
}
