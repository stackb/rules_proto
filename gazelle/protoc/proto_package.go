package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// ProtoPackage provides a set of proto_library derived rules for the package.
type ProtoPackage struct {
	// the build file currently being visited
	file *rule.File
	rel  string
	cfg  *ProtoPackageConfig
	libs []ProtoLibrary

	// cached/computed providers
	gen, empty []RuleProvider
}

// NewProtoPackage constructs a ProtoPackage given a list of proto_library rules
// in the package.
func NewProtoPackage(
	file *rule.File,
	rel string,
	cfg *ProtoPackageConfig,
	libs []ProtoLibrary) *ProtoPackage {
	return &ProtoPackage{
		file: file,
		rel:  rel,
		cfg:  cfg,
		libs: libs,
	}
}

// generateRules constructs a list of rules based on the configured set of
// languages.
func (s *ProtoPackage) generateRules(enabled bool) []RuleProvider {
	if enabled && s.gen != nil {
		return s.gen
	}
	if !enabled && s.empty != nil {
		return s.empty
	}

	// log.Printf("visiting package %s", s.rel)

	prelim := make([]RuleProvider, 0)

	for _, lang := range s.cfg.Languages() {
		if enabled != lang.Enabled {
			continue
		}
		rules := lang.Implementation.GenerateRules(s.rel, s.cfg, lang, s.libs)
		for _, rule := range rules {
			s.cfg.RegisterRuleProvider(label.Label{
				Repo: "", // TODO: how to know if we are in an external repo?
				Pkg:  s.rel,
				Name: rule.Rule().Name(),
			}, rule)
		}
		prelim = append(prelim, rules...)
	}

	rules := make([]RuleProvider, 0)
	for _, rule := range prelim {
		// Remove blacklisted rules unless they are specifically whitelisted
		if s.cfg.IsRuleExcluded(rule.Kind()) {
			continue
		}
		if s.cfg.IsRuleExcluded(rule.Name()) {
			continue
		}
		rules = append(rules, rule)
	}

	if enabled {
		s.gen = rules
	} else {
		s.empty = rules
	}

	return rules
}

// Rules provides the aggregated rule list for the package.
func (s *ProtoPackage) Rules() []*rule.Rule {
	return getProvidedRules(s.generateRules(true))
}

// Empty names the rules that can be deleted.
func (s *ProtoPackage) Empty() []*rule.Rule {
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
func (s *ProtoPackage) Imports() []interface{} {
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
