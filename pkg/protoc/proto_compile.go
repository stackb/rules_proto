package protoc

import (
	"fmt"
	"log"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// ProtoLibraryKey stores the ProtoLibrary implementation for a rule.
	ProtoLibraryKey = "_proto_library"
)

func init() {
	Rules().MustRegisterRule("stackb:rules_proto:proto_compile", &protoCompile{})
}

// protoCompile implements LanguageRule for the 'proto_compile' rule.
type protoCompile struct{}

// KindInfo implements part of the LanguageRule interface.
func (s *protoCompile) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		NonEmptyAttrs: map[string]bool{
			"outputs": true,
		},
		MergeableAttrs: map[string]bool{
			"outputs":         true,
			"plugins":         true,
			"output_mappings": true,
			"options":         true,
			"protos":          true,
		},
		SubstituteAttrs: map[string]bool{
			"out": true,
		},
	}
}

// Name implements part of the LanguageRule interface.
func (s *protoCompile) Name() string {
	return "proto_compile"
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoCompile) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules:proto_compile.bzl",
		Symbols: []string{"proto_compile"},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoCompile) ProvideRule(cfg *LanguageRuleConfig, config *ProtocConfiguration) RuleProvider {
	if len(config.Outputs) == 0 {
		return nil
	}
	return &protoCompileRule{
		kind:            "proto_compile",
		nameSuffix:      "compile",
		outputsAttrName: "outputs",
		config:          config,
		ruleConfig:      cfg,
	}
}

// protoCompile implements RuleProvider for the 'proto_compile' rule.
type protoCompileRule struct {
	kind            string
	nameSuffix      string
	outputsAttrName string
	config          *ProtocConfiguration
	ruleConfig      *LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *protoCompileRule) Kind() string {
	return s.kind
}

// Name implements part of the ruleProvider interface.
func (s *protoCompileRule) Name() string {
	return fmt.Sprintf("%s_%s_%s", s.config.Library.BaseName(), s.config.Prefix, s.nameSuffix)
}

// Visibility provides visibility labels.
func (s *protoCompileRule) Visibility() []string {
	return s.ruleConfig.GetVisibility()
}

func (s *protoCompileRule) Outputs() []string {
	outputs := s.config.Outputs
	sort.Strings(outputs)
	return outputs
}

// Rule implements part of the ruleProvider interface.
func (s *protoCompileRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	outputs := s.Outputs()

	// Check for output overlap with existing proto_compile rules of the same
	// kind. When a package-level plugin (e.g. protoc-gen-prost) produces the
	// same output file from multiple proto_library rules, merge them into a
	// single proto_compile rule using the "protos" attribute.
	for _, other := range otherGen {
		if other.Kind() != s.Kind() {
			continue
		}
		otherOutputs := other.AttrStrings(s.outputsAttrName)
		if !hasOverlap(outputs, otherOutputs) {
			continue
		}

		// Merge outputs
		other.SetAttr(s.outputsAttrName, DeduplicateAndSort(append(otherOutputs, outputs...)))

		// Convert singular "proto" to list "protos" if needed, then append
		existingProtos := other.AttrStrings("protos")
		if len(existingProtos) == 0 {
			if p := other.AttrString("proto"); p != "" {
				existingProtos = []string{p}
				other.DelAttr("proto")
			}
		}
		existingProtos = append(existingProtos, s.config.Library.Name())
		other.SetAttr("protos", DeduplicateAndSort(existingProtos))

		// Merge plugins
		otherPlugins := other.AttrStrings("plugins")
		otherPlugins = append(otherPlugins, GetPluginLabels(s.config.Plugins)...)
		other.SetAttr("plugins", DeduplicateAndSort(otherPlugins))

		// Merge output_mappings
		if len(s.config.Mappings) > 0 {
			existing := other.AttrStrings("output_mappings")
			for k, v := range s.config.Mappings {
				existing = append(existing, k+"="+v)
			}
			other.SetAttr("output_mappings", DeduplicateAndSort(existing))
		}

		return other
	}

	// No overlap found — create new rule
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr(s.outputsAttrName, outputs)
	newRule.SetAttr("plugins", GetPluginLabels(s.config.Plugins))
	newRule.SetAttr("proto", s.config.Library.Name())

	if s.config.LanguageConfig.Protoc != "" {
		newRule.SetAttr("protoc", s.config.LanguageConfig.Protoc)
	}

	if len(s.config.Mappings) > 0 {
		mappings := make([]string, len(s.config.Mappings))
		var i int
		for k, v := range s.config.Mappings {
			mappings[i] = k + "=" + v
			i++
		}
		sort.Strings(mappings)
		newRule.SetAttr("output_mappings", mappings)
	}

	outs := GetPluginOuts(s.config.Plugins)
	if len(outs) > 0 {
		newRule.SetAttr("outs", MakeStringDict(outs))
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	for _, name := range s.ruleConfig.GetAttrNames() {
		vals := s.ruleConfig.GetAttr(name)
		if len(vals) == 0 {
			continue
		}
		switch name {
		case "verbose":
			val := vals[0]
			switch val {
			case "True", "true":
				newRule.SetAttr("verbose", true)
			case "False", "false":
				newRule.SetAttr("verbose", false)
			default:
				log.Printf("bad attr 'verbose' value: %q", val)
			}
		default:
			newRule.SetAttr(name, vals)
		}
	}

	return newRule
}

// hasOverlap returns true if two string slices share any common element.
func hasOverlap(a, b []string) bool {
	set := make(map[string]bool, len(a))
	for _, s := range a {
		set[s] = true
	}
	for _, s := range b {
		if set[s] {
			return true
		}
	}
	return false
}

// Imports implements part of the RuleProvider interface.
func (s *protoCompileRule) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *protoCompileRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	options := GetPluginOptions(s.config.Plugins, r, from)
	if len(options) > 0 {
		r.SetAttr("options", MakeStringListDict(options))
	}
}
