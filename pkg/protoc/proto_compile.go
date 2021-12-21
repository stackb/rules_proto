package protoc

import (
	"fmt"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// ProtoLibraryKey stores the ProtoLibrary implementation for a rule.
	ProtoLibraryKey = "_proto_library"
)

func init() {
	Rules().MustRegisterRule("stackb:rules_proto:proto_compile", &protoCompile{})
}

// protoCompile implements LanguageRule for the 'proto_compile' rule.
type protoCompile struct{}

// KindInfo implements part of the LanguageRule interface.
func (s *protoCompile) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		NonEmptyAttrs: map[string]bool{
			"outputs": true,
		},
		MergeableAttrs: map[string]bool{
			"outputs":         true,
			"plugins":         true,
			"output_mappings": true,
			"options":         true,
		},
		SubstituteAttrs: map[string]bool{
			"out": true,
		},
	}
}

// Name implements part of the LanguageRule interface.
func (s *protoCompile) Name() string {
	return "proto_compile"
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoCompile) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules:proto_compile.bzl",
		Symbols: []string{"proto_compile"},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoCompile) ProvideRule(cfg *LanguageRuleConfig, config *ProtocConfiguration) RuleProvider {
	return &protoCompileRule{
		kind:            "proto_compile",
		nameSuffix:      "compile",
		outputsAttrName: "outputs",
		config:          config,
		ruleConfig:      cfg,
	}
}

// protoCompile implements RuleProvider for the 'proto_compile' rule.
type protoCompileRule struct {
	kind            string
	nameSuffix      string
	outputsAttrName string
	config          *ProtocConfiguration
	ruleConfig      *LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *protoCompileRule) Kind() string {
	return s.kind
}

// Name implements part of the ruleProvider interface.
func (s *protoCompileRule) Name() string {
	return fmt.Sprintf("%s_%s_%s", s.config.Library.BaseName(), s.config.Prefix, s.nameSuffix)
}

// Visibility provides visibility labels.
func (s *protoCompileRule) Visibility() []string {
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
func (s *protoCompileRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	outputs := s.config.Outputs
	sort.Strings(outputs)

	newRule.SetAttr(s.outputsAttrName, outputs)
	newRule.SetAttr("plugins", GetPluginLabels(s.config.Plugins))
	newRule.SetAttr("proto", s.config.Library.Name())

	if s.config.LanguageConfig.Protoc != "" {
		newRule.SetAttr("protoc", s.config.LanguageConfig.Protoc)
	}

	if len(s.config.Mappings) > 0 {
		mappings := make([]string, len(s.config.Mappings))
		var i int
		for k, v := range s.config.Mappings {
			mappings[i] = k + "=" + v
			i++
		}
		newRule.SetAttr("output_mappings", mappings)
	}

	outs := GetPluginOuts(s.config.Plugins)
	if len(outs) > 0 {
		newRule.SetAttr("outs", MakeStringDict(outs))
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *protoCompileRule) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *protoCompileRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	options := GetPluginOptions(s.config.Plugins, r, from)
	if len(options) > 0 {
		r.SetAttr("options", MakeStringListDict(options))
	}
}
