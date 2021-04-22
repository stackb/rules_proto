package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// ProtoPackage provides a set of proto_library derived rules for the package.
type ProtoPackage struct {
	rules []RuleProvider
}

// NewProtoPackage constructs a ProtoPackage given a list of proto_library rules
// in the package.
func NewProtoPackage(
	file *rule.File,
	rel string,
	cfg *ProtoPackageConfig,
	libs []ProtoLibrary) *ProtoPackage {

	prelim := make([]RuleProvider, 0)

	for _, lang := range cfg.Languages() {
		if !lang.Enabled {
			continue
		}
		rules := lang.Implementation.GenerateRules(rel, cfg, lang, libs)
		for _, rule := range rules {
			cfg.RegisterRuleProvider(label.Label{
				Repo: "", // TODO: how to know if we are in an external repo?
				Pkg:  rel,
				Name: rule.Rule().Name(),
			}, rule)
		}
		prelim = append(prelim, rules...)
	}

	rules := make([]RuleProvider, 0)
	for _, rule := range prelim {
		// Remove blacklisted rules unless they are specifically whitelisted
		if cfg.IsRuleExcluded(rule.Kind()) {
			continue
		}
		if cfg.IsRuleExcluded(rule.Name()) {
			continue
		}
		rules = append(rules, rule)

		// if the rule implements FileVisitor, give it the chance to mutate the
		// File.
		var i interface{} = rule
		if visitor, ok := i.(FileVisitor); ok {
			visitor.VisitFile(file)
		}
	}

	return &ProtoPackage{rules}
}

// Rules provides the aggregated rule list for the package.
func (s *ProtoPackage) Rules() []*rule.Rule {
	rules := make([]*rule.Rule, 0)
	for _, r := range s.rules {
		rules = append(rules, r.Rule())
	}
	// log.Printf("%d rules generated", len(rules))
	return rules
}

// Imports provides the aggregated list of imports for the package.
func (s *ProtoPackage) Imports() []interface{} {
	imports := make([]interface{}, 0)
	for _, r := range s.rules {
		for _, v := range r.Imports() {
			imports = append(imports, v)
		}
	}
	return imports
}
