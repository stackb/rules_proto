package rules_python

import (
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

var pyLibraryKindInfo = rule.KindInfo{
	MergeableAttrs: map[string]bool{
		"srcs":       true,
		"deps":       true,
		"visibility": true,
	},
	ResolveAttrs: map[string]bool{"deps": true},
}

// PyLibrary implements RuleProvider for 'py_library'-derived rules.
type PyLibrary struct {
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       protoc.DepsResolver
}

// Kind implements part of the ruleProvider interface.
func (s *PyLibrary) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *PyLibrary) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *PyLibrary) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.Outputs {
		if strings.HasSuffix(output, ".py") {
			srcs = append(srcs, protoc.StripRel(s.Config.Rel, output))
		}
	}
	return srcs
}

// Deps computes the deps list for the rule.
func (s *PyLibrary) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility provides visibility labels.
func (s *PyLibrary) Visibility() []string {
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
func (s *PyLibrary) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())

	deps := s.Deps()
	if len(deps) > 0 {
		newRule.SetAttr("deps", deps)
	}
	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *PyLibrary) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	if lib, ok := r.PrivateAttr(protoc.ProtoLibraryKey).(protoc.ProtoLibrary); ok {
		return protoc.ProtoLibraryImportSpecsForKind(r.Kind(), lib)
	}
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *PyLibrary) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	s.Resolver(c, ix, r, imports, from)
}
