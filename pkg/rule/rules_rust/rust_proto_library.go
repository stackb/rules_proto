package rules_rust

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	rustProtoLibraryRuleName   = "rust_proto_library"
	rustGrpcLibraryRuleName    = "rust_grpc_library"
	rustProtoLibraryRuleSuffix = "_rust_proto"
	rustGrpcLibraryRuleSuffix  = "_rust_grpc"
)

var rustProtoLibraryKindInfo = rule.KindInfo{
	MergeableAttrs: map[string]bool{
		"deps": true,
	},
	ResolveAttrs: map[string]bool{
		"deps": true,
	},
}

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:rust_proto_library", &rustProtoLibrary{
		name: rustProtoLibraryRuleName,
	})
	protoc.Rules().MustRegisterRule("stackb:rules_proto:rust_grpc_library", &rustProtoLibrary{
		name: rustGrpcLibraryRuleName,
	})
}

// rustProtoLibrary implements LanguageRule for the 'rust_proto_library' rule.
type rustProtoLibrary struct {
	name string
}

// Name implements part of the LanguageRule interface.
func (s *rustProtoLibrary) Name() string {
	return s.name
}

// KindInfo implements part of the LanguageRule interface.
func (s *rustProtoLibrary) KindInfo() rule.KindInfo {
	return rustProtoLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *rustProtoLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@rules_rust//proto:proto.bzl",
		Symbols: []string{s.name},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *rustProtoLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	if protoc.HasServices(pc.Library.Files()...) {
		if s.name == rustProtoLibraryRuleName {
			return nil
		}
	} else {
		if s.name != rustProtoLibraryRuleName {
			return nil
		}
	}
	return &rustProtoLibraryRule{
		KindName:       s.name,
		RuleNameSuffix: rustProtoLibraryRuleSuffix,
		RuleConfig:     cfg,
		Config:         pc,
	}
}

// rustProtoLibraryRule implements RuleProvider for 'rust_proto_library'-derived rules.
type rustProtoLibraryRule struct {
	KindName       string
	RuleNameSuffix string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       protoc.DepsResolver
}

// Kind implements part of the ruleProvider interface.
func (s *rustProtoLibraryRule) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *rustProtoLibraryRule) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Deps computes the deps list for the rule.
func (s *rustProtoLibraryRule) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility provides visibility labels.
func (s *rustProtoLibraryRule) Visibility() []string {
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
func (s *rustProtoLibraryRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("deps", []string{s.Config.Library.Name()})

	deps := s.Deps()
	if len(deps) > 0 {
		newRule.SetAttr("rust_deps", deps)
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *rustProtoLibraryRule) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *rustProtoLibraryRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
}
