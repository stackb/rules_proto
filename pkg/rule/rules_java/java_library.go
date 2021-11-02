package rules_java

import (
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

var javaLibraryKindInfo = rule.KindInfo{
	MergeableAttrs: map[string]bool{
		"srcs":    true,
		"exports": true,
	},
	ResolveAttrs: map[string]bool{"deps": true},
}

// JavaLibrary implements RuleProvider for 'java_library'-derived rules.
type JavaLibrary struct {
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       protoc.DepsResolver
}

// Kind implements part of the ruleProvider interface.
func (s *JavaLibrary) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *JavaLibrary) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *JavaLibrary) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.Outputs {
		if strings.HasSuffix(output, ".srcjar") {
			srcs = append(srcs, protoc.StripRel(s.Config.Rel, output))
		}
	}
	return srcs
}

// Deps computes the deps list for the rule.
func (s *JavaLibrary) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility provides visibility labels.
func (s *JavaLibrary) Visibility() []string {
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
func (s *JavaLibrary) Rule(otherGen ...*rule.Rule) *rule.Rule {
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
func (s *JavaLibrary) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	if lib, ok := r.PrivateAttr(protoc.ProtoLibraryKey).(protoc.ProtoLibrary); ok {
		return protoc.ProtoLibraryImportSpecsForKind(r.Kind(), lib)
	}
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *JavaLibrary) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	if s.Resolver == nil {
		return
	}
	s.Resolver(c, ix, r, imports, from)
}
