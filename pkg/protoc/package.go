package protoc

import (
	"log"
	"path"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	ruleProviderKey = "_x_rule_provider"
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
	// ruleLibs records the ProtoLibrary a RuleProvider was built on.
	ruleLibs map[RuleProvider]ProtoLibrary
}

// NewPackage constructs a Package given a list of proto_library rules
// in the package.
func NewPackage(rel string, cfg *PackageConfig, libs ...ProtoLibrary) *Package {
	s := &Package{
		rel:      rel,
		cfg:      cfg,
		libs:     libs,
		ruleLibs: make(map[RuleProvider]ProtoLibrary),
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

	for name, want := range p.Plugins {
		if !want {
			continue
		}
		plugin, ok := s.cfg.plugins[name]
		if !ok {
			log.Fatalf("plugin not configured: %q", name)
		}
		if !plugin.Enabled {
			continue
		}

		ctx := &PluginContext{
			Rel:           s.rel,
			ProtoLibrary:  lib,
			PackageConfig: *s.cfg,
			PluginConfig:  *plugin,
		}

		if plugin.Implementation == "" {
			plugin.Implementation = plugin.Name
		}
		impl, err := globalRegistry.LookupPlugin(plugin.Implementation)
		if err == ErrUnknownPlugin {
			log.Fatalf(
				"%s: plugin not registered: %q (available: %v) [%+v]",
				s.rel,
				plugin.Implementation,
				globalRegistry.PluginNames(),
				plugin,
			)
		}
		ctx.Plugin = impl

		// Delegate to the implementation for configuration
		config := impl.Configure(ctx)
		if config == nil {
			continue
		}
		config.Plugin = impl
		config.Config = plugin.clone()
		config.Options = DeduplicateAndSort(config.Options)

		// plugin.Label overrides the default value from the implementation
		if plugin.Label.Name != "" {
			config.Label = plugin.Label
		}

		configs = append(configs, config)
	}

	if len(configs) == 0 {
		return nil
	}

	imports := make([]string, len(lib.Files()))
	for i, file := range lib.Files() {
		imports[i] = path.Join(file.Dir, file.Basename)
	}

	rules := make([]RuleProvider, 0)

	pc := newProtocConfiguration(p, s.cfg.config.WorkDir, s.rel, p.Name, lib, configs)
	for _, name := range p.GetRulesByIntent(true) {
		ruleConfig, ok := s.cfg.rules[name]
		if !ok {
			names := make([]string, 0)
			for name := range s.cfg.rules {
				names = append(names, name)
			}
			log.Fatalf("proto_rule %q is not configured (available: %v)", name, names)
		}
		if !ruleConfig.Enabled {
			continue
		}
		impl, err := globalRegistry.LookupRule(ruleConfig.Implementation)
		if err == ErrUnknownRule {
			log.Fatalf(
				"%s: rule not registered: %q (available: %v)",
				s.rel,
				ruleConfig.Implementation,
				globalRegistry.RuleNames(),
			)
		}
		ruleConfig.Impl = impl

		rule := impl.ProvideRule(ruleConfig, pc)
		if rule == nil {
			continue
		}

		s.ruleLibs[rule] = lib

		rules = append(rules, rule)
	}

	return rules
}

// RuleProvider returns the provider of a rule or nil if not known.
func (s *Package) RuleProvider(r *rule.Rule) RuleProvider {
	if provider, ok := r.PrivateAttr(ruleProviderKey).(RuleProvider); ok {
		return provider
	}
	return nil
}

// Rules provides the aggregated rule list for the package.
func (s *Package) Rules() []*rule.Rule {
	return s.getProvidedRules(s.gen, true)
}

// Empty names the rules that can be deleted.
func (s *Package) Empty() []*rule.Rule {
	// it's a bit sad that we construct the full rules only for their kind and
	// name, but that's how it is right now.
	rules := s.getProvidedRules(s.empty, false)

	empty := make([]*rule.Rule, len(rules))
	for i, r := range rules {
		empty[i] = rule.NewRule(r.Kind(), r.Name())
	}

	return empty
}

func (s *Package) getProvidedRules(providers []RuleProvider, shouldResolve bool) []*rule.Rule {
	rules := make([]*rule.Rule, 0)
	for _, p := range providers {
		r := p.Rule(rules...)
		if r == nil {
			continue
		}

		if shouldResolve {
			// record the association of the rule provider here for the resolver.
			r.SetPrivateAttr(ruleProviderKey, p)

			imports := r.PrivateAttr(config.GazelleImportsKey)
			if imports == nil {
				lib := s.ruleLibs[p]
				r.SetPrivateAttr(ProtoLibraryKey, lib)
				r.SetPrivateAttr(config.GazelleImportsKey, lib.Imports())
			}

			// NOTE: this is a bit of a hack: it would be preferable to populate
			// the global resolver with import specs during the .Imports()
			// function.  One would think that the RuleProvider could be set as
			// a PrivateAttr to be retrieved in the Imports() function. However,
			// the rule ref seems to have changed by that time, the PrivateAttr
			// is removed.  Maybe this is due to rule merges?  Very difficult to
			// track down bug that cost me days.
			from := label.New(s.cfg.config.RepoName, s.rel, r.Name())
			file := rule.EmptyFile("", s.rel)
			provideResolverImportSpecs(s.cfg.config, p, r, file, from)
		}

		rules = append(rules, r)
	}
	return rules
}

func provideResolverImportSpecs(c *config.Config, provider RuleProvider, r *rule.Rule, f *rule.File, from label.Label) {
	for _, imp := range provider.Imports(c, r, f) {
		GlobalResolver().Provide(
			"protobuf",
			imp.Lang,
			imp.Imp,
			from,
		)
	}
}

// DeduplicateAndSort removes duplicate entries and sorts the list
func DeduplicateAndSort(in []string) (out []string) {
	seen := make(map[string]bool)
	for _, v := range in {
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	sort.Strings(out)
	return
}
