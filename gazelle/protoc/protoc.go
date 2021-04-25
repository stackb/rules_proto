package protoc

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
)

// protoc is the default protoc language singleton.
var protoc = newProtocLanguage()

func newProtocLanguage() *protocLanguage {
	return &protocLanguage{
		rules:         make(map[string]ProtoRule),
		plugins:       make(map[string]ProtoPlugin),
		ruleProviders: make(map[label.Label]RuleProvider),
	}
}

// protocLanguage implements language.Language, ProtoRuleRegistry, and
// ProtoPluginRegisty.
type protocLanguage struct {
	rules   map[string]ProtoRule
	plugins map[string]ProtoPlugin
	// ruleProviders is a mapping from label -> the provider that produced the
	// rule. we save this in the config such that we can retrieve the
	// association later in the resolve step.
	ruleProviders map[label.Label]RuleProvider
}

// mustLookupProtoRule returns the given rule or panics.
func (p *protocLanguage) mustLookupProtoRule(name string) ProtoRule {
	rule, ok := p.rules[name]
	if !ok {
		panic("unknown rule: " + name)
	}
	return rule
}

// RegisterRuleProvider implements the RuleProviderRegistry interface.
func (c *protocLanguage) RegisterRuleProvider(l label.Label, provider RuleProvider) {
	c.ruleProviders[l] = provider
}

// RuleNames implements part of the ProtoRuleRegistry interface.
func (p *protocLanguage) RuleNames() []string {
	names := make([]string, 0)
	for name := range p.rules {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// MustRegisterProtoRule implements part of the ProtoRuleRegistry interface.
func (p *protocLanguage) MustRegisterProtoRule(name string, rule ProtoRule) {
	_, ok := p.rules[name]
	if ok {
		panic("duplicate proto_rule registration: " + name)
	}
	p.rules[name] = rule
}

// LookupProtoRule implements part of the ProtoRuleRegistry interface.
func (p *protocLanguage) LookupProtoRule(name string) (ProtoRule, error) {
	rule, ok := p.rules[name]
	if !ok {
		return nil, ErrUnknownRule
	}
	return rule, nil
}

// PluginNames implements part of the ProtoPluginRegistry interface.
func (p *protocLanguage) PluginNames() []string {
	names := make([]string, 0)
	for name := range p.plugins {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// MustRegisterProtoPlugin implements part of the ProtoPluginRegistry interface.
func (p *protocLanguage) MustRegisterProtoPlugin(name string, plugin ProtoPlugin) {
	_, ok := p.plugins[name]
	if ok {
		panic("duplicate proto_plugin registration: " + name)
	}
	p.plugins[name] = plugin
}

// LookupProtoPlugin implements part of the ProtoPluginRegistry interface.
func (p *protocLanguage) LookupProtoPlugin(name string) (ProtoPlugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, ErrUnknownPlugin
	}
	return plugin, nil
}
