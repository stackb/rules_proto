package rules_nodejs

import (
	"log"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoTsLibraryRuleName   = "proto_ts_library"
	ProtoTsLibraryRuleSuffix = "_ts_proto"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_ts_library", &protoTsLibrary{})
}

// protoTsLibrary implements LanguageRule for the 'proto_ts_library' rule from
// @rules_proto.
type protoTsLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoTsLibrary) Name() string {
	return ProtoTsLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoTsLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs": true,
			"tsc":  true,
			"args": true,
		},
		ResolveAttrs: map[string]bool{
			"deps": true,
		},
	}

}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoTsLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/ts:proto_ts_library.bzl",
		Symbols: []string{ProtoTsLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoTsLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := make([]string, 0)
	for _, out := range pc.Outputs {
		if strings.HasSuffix(out, ".ts") {
			outputs = append(outputs, out)
		}
	}
	if len(outputs) == 0 {
		return nil
	}
	return &tsLibrary{
		KindName:       ProtoTsLibraryRuleName,
		RuleNameSuffix: ProtoTsLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver:       protoc.ResolveDepsAttr("deps", true),
	}
}

// tsLibrary implements RuleProvider for 'ts_library'-like rules.
type tsLibrary struct {
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       protoc.DepsResolver
}

// Kind implements part of the ruleProvider interface.
func (s *tsLibrary) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *tsLibrary) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *tsLibrary) Srcs() []string {
	return s.Outputs
}

// Deps computes the deps list for the rule.
func (s *tsLibrary) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility provides visibility labels.
func (s *tsLibrary) Visibility() []string {
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
func (s *tsLibrary) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())

	deps := s.Deps()
	if len(deps) > 0 {
		newRule.SetAttr("deps", deps)
	}

	tsc := s.RuleConfig.GetAttr("tsc")
	if len(tsc) > 0 {
		if len(tsc) > 1 {
			log.Printf("warning (%s): found multiple entries for 'tsc', choosing last one: %v", s.Kind(), tsc)
		}
		newRule.SetAttr("tsc", tsc[len(tsc)-1])
	}

	args := s.RuleConfig.GetAttr("args")
	if len(args) > 0 {
		newRule.SetAttr("args", args)
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *tsLibrary) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	if lib, ok := r.PrivateAttr(protoc.ProtoLibraryKey).(protoc.ProtoLibrary); ok {
		return protoc.ProtoLibraryImportSpecsForKind(r.Kind(), lib)
	}
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *tsLibrary) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	if s.Resolver == nil {
		return
	}
	s.Resolver(c, ix, r, imports, from)
}
