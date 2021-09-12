package builtin

import (
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

var ccLibraryKindInfo = rule.KindInfo{
	MergeableAttrs: map[string]bool{
		"srcs":       true,
		"hdrs":       true,
		"deps":       true,
		"visibility": true,
	},
}

type Resolver func(impl *CcLibraryRule, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label)

// CcLibraryRule implements RuleProvider for 'cc_library'-derived rules.
type CcLibraryRule struct {
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       Resolver
}

// Kind implements part of the ruleProvider interface.
func (s *CcLibraryRule) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *CcLibraryRule) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *CcLibraryRule) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.Outputs {
		if strings.HasSuffix(output, ".cc") {
			srcs = append(srcs, derel(s.Config.Rel, output))
		}
	}
	return srcs
}

// Hdrs computes the hdrs list for the rule.
func (s *CcLibraryRule) Hdrs() []string {
	hdrs := make([]string, 0)
	for _, output := range s.Outputs {
		if strings.HasSuffix(output, ".h") {
			hdrs = append(hdrs, derel(s.Config.Rel, output))
		}
	}
	return hdrs
}

// Deps computes the deps list for the rule.
func (s *CcLibraryRule) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility implements part of the ruleProvider interface.
func (s *CcLibraryRule) Visibility() []string {
	visibility := make([]string, 0)
	for k, want := range s.RuleConfig.Visibility {
		if !want {
			continue
		}
		visibility = append(visibility, k)
	}
	sort.Strings(visibility)
	return visibility
}

// Rule implements part of the ruleProvider interface.
func (s *CcLibraryRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())
	newRule.SetAttr("hdrs", s.Hdrs())

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *CcLibraryRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
	if s.Resolver == nil {
		return
	}
	s.Resolver(s, c, r, importsRaw, from)
}

func derel(rel string, filename string) string {
	if !strings.HasPrefix(filename, rel) {
		return filename
	}
	return filename[len(rel)+1:] // +1 for slash separator
}
