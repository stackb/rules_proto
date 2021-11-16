package protoc

import (
	"sort"
)

// Plugins returns a reference to the global PluginRegistry
func Plugins() PluginRegistry {
	return globalRegistry
}

// Rules returns a reference to the global RuleRegistry
func Rules() RuleRegistry {
	return globalRegistry
}

// registry is the default registry singleton.
var globalRegistry = &registry{
	rules:   make(map[string]LanguageRule),
	plugins: make(map[string]Plugin),
}

// registry implements RuleRegistry and PluginRegisty.
type registry struct {
	rules   map[string]LanguageRule
	plugins map[string]Plugin
}

// RuleNames implements part of the RuleRegistry interface.
func (p *registry) RuleNames() []string {
	names := make([]string, 0)
	for name := range p.rules {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// MustRegisterRule implements part of the RuleRegistry interface.
func (p *registry) MustRegisterRule(name string, rule LanguageRule) RuleRegistry {
	_, ok := p.rules[name]
	if ok {
		panic("duplicate proto_rule registration: " + name)
	}
	p.rules[name] = rule
	return p
}

// LookupRule implements part of the RuleRegistry interface.
func (p *registry) LookupRule(name string) (LanguageRule, error) {
	rule, ok := p.rules[name]
	if !ok {
		return nil, ErrUnknownRule
	}
	return rule, nil
}

// PluginNames implements part of the PluginRegistry interface.
func (p *registry) PluginNames() []string {
	names := make([]string, 0)
	for name := range p.plugins {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// MustRegisterPlugin implements part of the PluginRegistry interface.
func (p *registry) MustRegisterPlugin(plugin Plugin) PluginRegistry {
	_, ok := p.plugins[plugin.Name()]
	if ok {
		panic("duplicate proto_plugin registration: " + plugin.Name())
	}
	p.plugins[plugin.Name()] = plugin
	return p
}

// LookupPlugin implements part of the PluginRegistry interface.
func (p *registry) LookupPlugin(name string) (Plugin, error) {
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, ErrUnknownPlugin
	}
	return plugin, nil
}
