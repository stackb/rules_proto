package builtin

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_cc_library", &protoCcLibrary{})
}

// protoCcLibrary implements LanguageRule for the 'proto_cc_library' rule from
// @rules_proto.
type protoCcLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoCcLibrary) Name() string {
	return "proto_cc_library"
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoCcLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs":       true,
			"hdrs":       true,
			"deps":       true,
			"visibility": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoCcLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/cc:proto_cc_library.bzl",
		Symbols: []string{"proto_cc_library"},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoCcLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, config *protoc.ProtocConfiguration) protoc.RuleProvider {
	return &protoCcLibraryRule{ruleConfig: cfg, config: config}
}

// protoCcLibrary implements RuleProvider for the 'proto_compile' rule.
type protoCcLibraryRule struct {
	config     *protoc.ProtocConfiguration
	ruleConfig *protoc.LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *protoCcLibraryRule) Kind() string {
	return "proto_cc_library"
}

// Name implements part of the ruleProvider interface.
func (s *protoCcLibraryRule) Name() string {
	return fmt.Sprintf("%s_cc_library", s.config.Library.BaseName())
}

// Srcs computes the srcs list for the rule.
func (s *protoCcLibraryRule) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.config.Outputs {
		if strings.HasSuffix(output, ".pb.cc") {
			srcs = append(srcs, output)
		}
	}
	return srcs
}

// Hdrs computes the hdrs list for the rule.
func (s *protoCcLibraryRule) Hdrs() []string {
	hdrs := make([]string, 0)
	for _, output := range s.config.Outputs {
		if strings.HasSuffix(output, ".pb.h") {
			hdrs = append(hdrs, output)
		}
	}
	return hdrs
}

// Deps computes the deps list for the rule.
func (s *protoCcLibraryRule) Deps() []string {
	deps := s.ruleConfig.GetDeps()
	// for _, output := range s.config.Outputs {
	// 	if strings.HasSuffix(output, ".pb.cc") {
	// 		deps = append(deps, output)
	// 	}
	// }
	return deps
}

// Visibility implements part of the ruleProvider interface.
func (s *protoCcLibraryRule) Visibility() []string {
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
func (s *protoCcLibraryRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())
	newRule.SetAttr("hdrs", s.Hdrs())

	deps := s.Deps()
	if len(deps) > 0 {
		newRule.SetAttr("deps", s.Deps())
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *protoCcLibraryRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
}
