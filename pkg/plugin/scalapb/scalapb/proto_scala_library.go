package scalapb

import (
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoScalaLibraryRuleName   = "proto_scala_library"
	ProtoScalaLibraryRuleSuffix = "_scala_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_scala_library", &protoScalaLibrary{})
}

// protoScalaLibrary implements LanguageRule for the 'proto_scala_library' rule from
// @rules_proto.
type protoScalaLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoScalaLibrary) Name() string {
	return ProtoScalaLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoScalaLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs":       true,
			"deps":       true,
			"visibility": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoScalaLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/scala:proto_scala_library.bzl",
		Symbols: []string{ProtoScalaLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoScalaLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("scalapb:scalapb:protoc-gen-scala")
	if len(outputs) == 0 {
		return nil
	}
	return &ScalaLibraryRule{
		KindName:       ProtoScalaLibraryRuleName,
		RuleNameSuffix: ProtoScalaLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver:       protoc.ResolveDepsWithSuffix(ProtoScalaLibraryRuleSuffix),
	}
}

// ScalaLibraryRule implements RuleProvider for 'cc_library'-derived rules.
type ScalaLibraryRule struct {
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       protoc.DepsResolver
}

// Kind implements part of the ruleProvider interface.
func (s *ScalaLibraryRule) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *ScalaLibraryRule) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *ScalaLibraryRule) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.Outputs {
		if strings.HasSuffix(output, ".srcjar") {
			srcs = append(srcs, protoc.StripRel(s.Config.Rel, output))
		}
	}
	return srcs
}

// Deps computes the deps list for the rule.
func (s *ScalaLibraryRule) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility implements part of the ruleProvider interface.
func (s *ScalaLibraryRule) Visibility() []string {
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
func (s *ScalaLibraryRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *ScalaLibraryRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
	if s.Resolver == nil {
		return
	}
	s.Resolver(s, s.Config, c, r, importsRaw, from)
}
