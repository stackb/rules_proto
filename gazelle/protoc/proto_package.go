package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// protoPackage provides a set of proto_library derived rules for the package.
type protoPackage struct {
	registry RuleProviderRegistry
	// the build file currently being visited
	file *rule.File
	rel  string
	cfg  *ProtoPackageConfig
	libs []ProtoLibrary

	// cached/computed providers
	gen, empty []RuleProvider
}

// newProtoPackage constructs a protoPackage given a list of proto_library rules
// in the package.
func newProtoPackage(
	registry RuleProviderRegistry,
	file *rule.File,
	rel string,
	cfg *ProtoPackageConfig,
	libs []ProtoLibrary) *protoPackage {
	return &protoPackage{
		registry: registry,
		file:     file,
		rel:      rel,
		cfg:      cfg,
		libs:     libs,
	}
}

// generateRules constructs a list of rules based on the configured set of
// languages.
func (s *protoPackage) generateRules(enabled bool) []RuleProvider {
	if enabled && s.gen != nil {
		return s.gen
	}
	if !enabled && s.empty != nil {
		return s.empty
	}

	rules := make([]RuleProvider, 0)

	for _, lang := range s.cfg.configuredLangs() {
		if enabled != lang.Enabled {
			continue
		}
		for _, lib := range s.libs {
			rules := s.libraryRules(s.rel, s.cfg, lang, lib)
			for _, rule := range rules {
				s.registry.RegisterRuleProvider(label.Label{
					Repo: "", // TODO: how to know if we are in an external repo?
					Pkg:  s.rel,
					Name: rule.Rule().Name(),
				}, rule)
			}
			rules = append(rules, rules...)
		}
	}

	if enabled {
		s.gen = rules
	} else {
		s.empty = rules
	}

	return rules
}

func (s *protoPackage) libraryRules(
	rel string,
	c *ProtoPackageConfig,
	p *ProtoLangConfig,
	lib ProtoLibrary,
) []RuleProvider {
	// list of plugin configurations that apply to this proto_library
	configs := make([]*PluginConfiguration, 0)

	for _, plugin := range p.Plugins {
		if !plugin.Implementation.ShouldApply(rel, c, lib) {
			continue
		}
		config := &PluginConfiguration{
			Label: plugin.Label,
			Srcs:  plugin.Implementation.GeneratedSrcs(rel, c, lib),
		}
		configs = append(configs, config)

		if provider, ok := plugin.Implementation.(PluginOptionsProvider); ok {
			config.Options = provider.GeneratedOptions(rel, c, lib)
		}
		if provider, ok := plugin.Implementation.(PluginMappingsProvider); ok {
			config.Mappings = provider.GeneratedMappings(rel, c, lib)
		}
		if provider, ok := plugin.Implementation.(PluginOutProvider); ok {
			config.Out = provider.GeneratedOut(rel, c, lib)
		}
	}
	if len(configs) == 0 {
		return nil
	}

	pc := newProtocConfiguration(rel, p.Name, lib, configs)

	rules := make([]RuleProvider, 0)
	for _, cfg := range p.Rules {
		rules = append(rules, cfg.Implementation.GenerateRule(cfg, pc))
	}
	return rules
}

// Rules provides the aggregated rule list for the package.
func (s *protoPackage) Rules() []*rule.Rule {
	return getProvidedRules(s.generateRules(true))
}

// Empty names the rules that can be deleted.
func (s *protoPackage) Empty() []*rule.Rule {
	rules := getProvidedRules(s.generateRules(false))

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
	return getProvidedImports(s.generateRules(true))
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
