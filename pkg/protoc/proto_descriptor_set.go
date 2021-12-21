package protoc

import (
	"fmt"
	"path"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func init() {
	Rules().MustRegisterRule("stackb:rules_proto:proto_descriptor_set", &protoDescriptorSetRule{})
	Plugins().MustRegisterPlugin(&protoDescriptorSetPlugin{})
}

// protoDescriptorSetRule implements LanguageRule for the 'proto_descriptor_set'
// rule from @rules_proto.
type protoDescriptorSetRule struct{}

// Name implements part of the LanguageRule interface.
func (s *protoDescriptorSetRule) Name() string {
	return "rules_proto_descriptor_set"
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoDescriptorSetRule) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"deps": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoDescriptorSetRule) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules:proto_descriptor_set.bzl",
		Symbols: []string{"rules_proto_descriptor_set"},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoDescriptorSetRule) ProvideRule(cfg *LanguageRuleConfig, config *ProtocConfiguration) RuleProvider {
	return &protoDescriptorSetRuleRule{ruleConfig: cfg, config: config}
}

// protoDescriptorSetRule implements RuleProvider for the 'proto_compile' rule.
type protoDescriptorSetRuleRule struct {
	config     *ProtocConfiguration
	ruleConfig *LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *protoDescriptorSetRuleRule) Kind() string {
	return "rules_proto_descriptor_set"
}

// Name implements part of the ruleProvider interface.
func (s *protoDescriptorSetRuleRule) Name() string {
	return fmt.Sprintf("%s_descriptor", s.config.Library.BaseName())
}

// Visibility provides visibility labels.
func (s *protoDescriptorSetRuleRule) Visibility() []string {
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
func (s *protoDescriptorSetRuleRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("deps", []string{s.config.Library.Name()})

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *protoDescriptorSetRuleRule) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *protoDescriptorSetRuleRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
}

type protoDescriptorSetPlugin struct{}

// Name implements part of the Plugin interface.
func (p *protoDescriptorSetPlugin) Name() string {
	return "bazelbuild:rules_proto:proto_descriptor_set"
}

// Configure implements part of the Plugin interface.
func (p *protoDescriptorSetPlugin) Configure(ctx *PluginContext) *PluginConfiguration {
	descriptorSetOut := path.Join(ctx.Rel, fmt.Sprintf("%s_descriptor.pb", ctx.ProtoLibrary.BaseName()))

	return &PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "bazelbuild/rules_proto", "proto_descriptor_set"),
		Outputs: []string{descriptorSetOut},
		Out:     ctx.Rel,
		Options: ctx.PluginConfig.GetOptions(),
	}
}
