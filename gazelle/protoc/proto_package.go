package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// ProtoPackage provides a specific set of rules for a given list of .proto
// files.
type ProtoPackage struct {
	files []*ProtoFile
	rules []RuleProvider
}

// NewProtoPackage constructs a ProtoPackage given a variadic list of
// ProtoFiles.
func NewProtoPackage(file *rule.File,
	rel string,
	cfg *protoPackageConfig,
	protoLibraries []*rule.Rule,
	files ...*ProtoFile) *ProtoPackage {

	prelim := make([]RuleProvider, 0)

	libs := make([]ProtoLibrary, 0)
	for _, rule := range protoLibraries {
		libs = append(libs, &OtherProtoLibrary{rule: rule, files: files})
	}

	for _, lang := range cfg.Languages() {
		rules := lang.GenerateRules(rel, cfg, libs)
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

	return &ProtoPackage{
		files: files,
		rules: rules,
	}
}

// Rules provides the aggregated rule list for the package.
func (s *ProtoPackage) Rules() []*rule.Rule {
	rules := make([]*rule.Rule, 0)
	for _, r := range s.rules {
		rules = append(rules, r.Rule())
	}
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
