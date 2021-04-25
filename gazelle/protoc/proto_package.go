package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// protoPackage provides a set of proto_library derived rules for the package.
type protoPackage struct {
	// relative path of build file
	rel string
	// the config for this package
	cfg *ProtoPackageConfig
	// list of proto_library targets in the package
	libs []ProtoLibrary
	// computed providers
	gen, empty []RuleProvider
}

// newProtoPackage constructs a protoPackage given a list of proto_library rules
// in the package.
func newProtoPackage(rel string, cfg *ProtoPackageConfig, libs ...ProtoLibrary) *protoPackage {
	s := &protoPackage{
		rel:  rel,
		cfg:  cfg,
		libs: libs,
	}
	s.gen = s.generateRules(true)
	s.empty = s.generateRules(false)
	return s
}

// generateRules constructs a list of rules based on the configured set of
// languages.
func (s *protoPackage) generateRules(enabled bool) []RuleProvider {

	rules := make([]RuleProvider, 0)
	langs := s.cfg.configuredLangs()

	for _, lang := range langs {
		if enabled != lang.Enabled {
			continue
		}
		for _, lib := range s.libs {
			rules = append(rules, s.libraryRules(lang, lib)...)
		}
	}
	return rules
}

func (s *protoPackage) libraryRules(p *ProtoLangConfig, lib ProtoLibrary) []RuleProvider {
	// list of plugin configurations that apply to this proto_library
	configs := make([]*PluginConfiguration, 0)

	for _, plugin := range p.Plugins {
		if !plugin.Implementation.ShouldApply(s.rel, s.cfg, lib) {
			continue
		}
		config := &PluginConfiguration{
			Label: plugin.Label,
			Srcs:  plugin.Implementation.GeneratedSrcs(s.rel, s.cfg, lib),
		}
		if provider, ok := plugin.Implementation.(PluginOptionsProvider); ok {
			config.Options = provider.GeneratedOptions(s.rel, s.cfg, lib)
		}
		if provider, ok := plugin.Implementation.(PluginMappingsProvider); ok {
			config.Mappings = provider.GeneratedMappings(s.rel, s.cfg, lib)
		}
		if provider, ok := plugin.Implementation.(PluginOutProvider); ok {
			config.Out = provider.GeneratedOut(s.rel, s.cfg, lib)
		}
		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil
	}

	rules := make([]RuleProvider, 0)

	pc := newProtocConfiguration(s.rel, p.Name, lib, configs)
	for _, cfg := range p.Rules {
		rules = append(rules, cfg.Implementation.GenerateRule(cfg, pc))
	}

	return rules
}

// RuleProviders returns the list of generated rules.
func (s *protoPackage) RuleProviders() []RuleProvider {
	return s.gen
}

// Rules provides the aggregated rule list for the package.
func (s *protoPackage) Rules() []*rule.Rule {
	return getProvidedRules(s.gen)
}

// Empty names the rules that can be deleted.
func (s *protoPackage) Empty() []*rule.Rule {
	rules := getProvidedRules(s.empty)

	// it's a bit sad that we construct the full rules only for their kind and
	// name, but that's how it is right now.
	empty := make([]*rule.Rule, len(rules))
	for i, r := range rules {
		empty[i] = rule.NewRule(r.Kind(), r.Name())
	}

	return empty
}

// Imports provides the aggregated list of imports for the package.
func (s *protoPackage) Imports() []interface{} {
	return getProvidedImports(s.gen)
}

func getProvidedRules(providers []RuleProvider) []*rule.Rule {
	rules := make([]*rule.Rule, len(providers))
	for i, p := range providers {
		rules[i] = p.Rule()
	}
	return rules
}

func getProvidedImports(providers []RuleProvider) []interface{} {
	imports := make([]interface{}, 0)
	for _, r := range providers {
		for _, v := range r.Imports() {
			imports = append(imports, v)
		}
	}
	return imports
}
