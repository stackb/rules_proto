package protoc

import (
	"log"

	"github.com/bazelbuild/bazel-gazelle/rule"
)

// Package provides a set of proto_library derived rules for the package.
type Package struct {
	// relative path of build file
	rel string
	// the config for this package
	cfg *PackageConfig
	// list of proto_library targets in the package
	libs []ProtoLibrary
	// computed providers
	gen, empty []RuleProvider
}

// NewPackage constructs a Package given a list of proto_library rules
// in the package.
func NewPackage(rel string, cfg *PackageConfig, libs ...ProtoLibrary) *Package {
	s := &Package{
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
func (s *Package) generateRules(enabled bool) []RuleProvider {

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

func (s *Package) libraryRules(p *LanguageConfig, lib ProtoLibrary) []RuleProvider {
	// list of plugin configurations that apply to this proto_library
	configs := make([]*PluginConfiguration, 0)

	for _, plugin := range p.Plugins {
		if !plugin.Enabled {
			continue
		}

		ctx := &PluginContext{
			Rel:           s.rel,
			ProtoLibrary:  lib,
			PackageConfig: *s.cfg,
			PluginConfig:  *plugin,
		}

		impl, err := globalRegistry.LookupPlugin(plugin.Implementation)
		if err == ErrUnknownPlugin {
			log.Fatalf(
				"plugin not registered: %q (available: %v)", plugin.Implementation,
				globalRegistry.PluginNames())
		}

		config := &PluginConfiguration{}

		// Delegate to the implementation for configuration
		impl.Configure(ctx, config)

		// if implementation says "skip", abort now.
		if config.Skip {
			continue
		}

		// plugin.Label overrides the default value from the implementation
		if plugin.Label.Name != "" {
			config.Label = plugin.Label
		}

		configs = append(configs, config)
	}
	if len(configs) == 0 {
		return nil
	}

	rules := make([]RuleProvider, 0)

	pc := newProtocConfiguration(s.rel, p.Name, lib, configs)
	for _, cfg := range p.Rules {
		rules = append(rules, cfg.Implementation.ProvideRule(cfg, pc))
	}

	return rules
}

// RuleProviders returns the list of generated rules.
func (s *Package) RuleProviders() []RuleProvider {
	return s.gen
}

// Rules provides the aggregated rule list for the package.
func (s *Package) Rules() []*rule.Rule {
	return getProvidedRules(s.gen)
}

// Empty names the rules that can be deleted.
func (s *Package) Empty() []*rule.Rule {
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
func (s *Package) Imports() []interface{} {
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
