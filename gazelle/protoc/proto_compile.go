package protoc

import (
	"fmt"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const ProtoCompileName = "proto_compile"

func init() {
	MustRegisterProtoRule(ProtoCompileName, &ProtoCompile{})
}

// ProtoCompile implements ProtoRule for the 'proto_compile' rule.
type ProtoCompile struct {
}

// KindInfo implements part of the ProtoRule interface.
func (s *ProtoCompile) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		NonEmptyAttrs: map[string]bool{"genfiles": true},
		MergeableAttrs: map[string]bool{
			"genfiles": true,
			"plugins":  true,
		},
	}
}

// LoadInfo implements part of the ProtoRule interface.
func (s *ProtoCompile) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name: "@build_stack_rules_proto//rules:proto_compile.bzl",
		Symbols: []string{
			"proto_compile",
		},
	}
}

// GenerateRule implements part of the ProtoRule interface.
func (s *ProtoCompile) GenerateRule(cfg *ProtoRuleConfig, config *ProtocConfiguration) RuleProvider {
	return &protoCompileRule{config}
}

// ProtoCompile implements RuleProvider for the 'proto_compile' rule.
type protoCompileRule struct {
	config *ProtocConfiguration
}

// Kind implements part of the ruleProvider interface.
func (s *protoCompileRule) Kind() string {
	return ProtoCompileName
}

// Name implements part of the ruleProvider interface.
func (s *protoCompileRule) Name() string {
	return fmt.Sprintf("%s_%s_compile", s.config.Library.BaseName(), s.config.Prefix)
}

// Imports implements part of the ruleProvider interface.
func (s *protoCompileRule) Imports() []string {
	return []string{s.Kind()}
}

// Visibility implements part of the ruleProvider interface.
func (s *protoCompileRule) Visibility() []string {
	return nil // TODO: visibility feature?
}

// Rule implements part of the ruleProvider interface.
func (s *protoCompileRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("genfiles", s.config.GeneratedSrcs)
	newRule.SetAttr("plugins", GetPluginLabels(s.config.Plugins))
	newRule.SetAttr("proto", s.config.Library.Name())

	if len(s.config.GeneratedMappings) > 0 {
		newRule.SetAttr("mappings", makeStringDict(s.config.GeneratedMappings))
	}

	options := GetPluginOptions(s.config.Plugins)
	if len(options) > 0 {
		newRule.SetAttr("options", makeStringListDict(options))
	}

	outs := GetPluginOuts(s.config.Plugins)
	if len(outs) > 0 {
		newRule.SetAttr("outs", makeStringDict(outs))
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *protoCompileRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
}
