package rules_scala

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Rules().MustRegisterRule("bazelbuild:rules_scala:scala_proto_library",
		&scalaProtoLibrary{})
}

// scalaProtoLibrary implements LanguageRule for the 'scala_proto_library' rule
// from @rules_scala.
type scalaProtoLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *scalaProtoLibrary) Name() string {
	return "scala_proto_library"
}

// KindInfo implements part of the LanguageRule interface.
func (s *scalaProtoLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs":       true,
			"deps":       true,
			"visibility": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *scalaProtoLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/scala:scala_proto_library.bzl",
		Symbols: []string{s.Name()},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *scalaProtoLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	return &scalaProtoLibraryRule{
		kindName:   s.Name(),
		ruleConfig: cfg,
		config:     pc,
	}
}

// scalaProtoLibraryRule implements RuleProvider for 'scala_library'-derived rules.
type scalaProtoLibraryRule struct {
	kindName   string
	config     *protoc.ProtocConfiguration
	ruleConfig *protoc.LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *scalaProtoLibraryRule) Kind() string {
	return s.kindName
}

// Name implements part of the ruleProvider interface.
func (s *scalaProtoLibraryRule) Name() string {
	return s.config.Library.BaseName() + "_scala_proto"
}

// Visibility provides visibility labels.
func (s *scalaProtoLibraryRule) Visibility() []string {
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
func (s *scalaProtoLibraryRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *scalaProtoLibraryRule) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *scalaProtoLibraryRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	r.SetAttr("deps", []string{":" + s.config.Library.Name()})
}
