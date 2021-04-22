package protoc

import (
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// protoRuleDirective is the directive for toggling rule generation.
	protoRuleDirective = "proto_rule"
	// protoLanguageDirective tells gazelle which languages a package should
	// produce
	protoLanguageDirective = "proto_language"
	// protoPluginDirective created an association between proto_language
	// and the label of a proto_plugin.
	protoPluginDirective = "proto_plugin"
	// prefixDirective is the same as 'gazelle:prefix'
	importpathPrefixDirective = "prefix"
)

// protoPackageConfig represents the config extension for the rosetta language.
type protoPackageConfig struct {
	// ruleProviders is a mapping from label -> the provider that produced
	// the rule. we save this in the config such that we can retrieve the
	// association later in the resolve step.
	ruleProviders map[label.Label]RuleProvider
	// exclude patterns for rules that should be skipped for this package.
	rules map[string]*protoRuleConfig
	// exclude patterns for rules that should be skipped for this package.
	plugins map[string]*protoPluginConfig
	// configured languages
	languages map[string]ProtoLanguage
	// the gazelle:prefix for golang
	importpathPrefix string
	// IMPORTANT! Adding new fields here?  Don't forget to copy it in the Clone
	// method!
}

// newProtoPackageConfig initializes a new protoPackageConfig.
func newProtoPackageConfig() *protoPackageConfig {
	return &protoPackageConfig{
		ruleProviders: make(map[label.Label]RuleProvider),
		rules:         make(map[string]*protoRuleConfig),
		languages:     make(map[string]ProtoLanguage),
		plugins:       make(map[string]*protoPluginConfig),
	}
}

// parseDirectives is called in each directory visited by gazelle.  The relative
// directory name is given by 'rel' and the list of directives in the BUILD file
// are specified by 'directives'.
func (c *protoPackageConfig) parseDirectives(rel string, directives []rule.Directive) {
	for _, d := range directives {
		switch d.Key {
		case importpathPrefixDirective:
			c.importpathPrefix = strings.TrimSpace(d.Value)
		case protoPluginDirective:
			fields := strings.Fields(d.Value)
			if len(fields) != 2 {
				panic(fmt.Sprintf(
					"invalid directive %s: expected {PLUGIN_NAME} {LABEL}, got %q",
					protoPluginDirective, d.Value))
			}
			pluginName := fields[0]
			pluginLabel := fields[1]
			l, err := label.Parse(pluginLabel)
			if err != nil {
				panic(fmt.Sprintf(
					"invalid plugin label: %s", pluginLabel))
			}
			plugin, err := LookupProtoPlugin(pluginName)
			if err != nil {
				panic(fmt.Sprintf("error while processing gazelle directive %q: %v", pluginName, err))
			}
			c.plugins[pluginName] = &protoPluginConfig{Label: l, Name: pluginName, Plugin: plugin}
			log.Printf("Added proto_plugin: %s -> %v", pluginName, l)
		case protoRuleDirective:
			pattern := strings.TrimSpace(d.Value)
			negative := strings.HasPrefix(pattern, "-")
			positive := strings.HasPrefix(pattern, "+")
			if negative || positive {
				pattern = pattern[1:]
			}
			if _, err := path.Match(pattern, ""); err == path.ErrBadPattern {
				panic(fmt.Sprintf("invalid directive %s: bad match pattern %q", protoRuleDirective, pattern))
			}
			cfg := c.rules[pattern]
			if cfg == nil {
				cfg = &protoRuleConfig{pattern: pattern}
				c.rules[pattern] = cfg
			}
			cfg.exclude = negative
		case protoLanguageDirective:
			name := strings.TrimSpace(d.Value)
			negative := strings.HasPrefix(name, "-")
			positive := strings.HasPrefix(name, "+")
			if negative || positive {
				name = name[1:]
			}
			lang, err := LookupProtoLanguage(name)
			if err != nil {
				panic(fmt.Sprintf("invalid or unknown plugin %q: %v", name, err))
			}
			if negative {
				delete(c.languages, name)
			} else {
				c.languages[name] = lang
			}
		}
	}
}

func (c *protoPackageConfig) IsRuleExcluded(name string) bool {
	if c.IsRuleIncluded(name) {
		return false
	}
	if len(c.rules) == 0 {
		return false
	}
	for _, cfg := range c.rules {
		if cfg.IsRuleExcluded(name) {
			return true
		}
	}
	return false
}

func (c *protoPackageConfig) IsRuleIncluded(name string) bool {
	if len(c.rules) == 0 {
		return false
	}
	for _, cfg := range c.rules {
		if cfg.IsRuleIncluded(name) {
			return true
		}
	}
	return false
}

// Languages returns a determinstic ordered list of configured languages
func (c *protoPackageConfig) Languages() []ProtoLanguage {
	names := make([]string, 0)
	for name := range c.languages {
		names = append(names, name)
	}
	sort.Strings(names)
	langs := make([]ProtoLanguage, 0)
	for _, name := range names {
		langs = append(langs, c.languages[name])
	}
	return langs
}

// Clone copies this config to a new one
func (c *protoPackageConfig) Clone() *protoPackageConfig {
	clone := newProtoPackageConfig()
	clone.importpathPrefix = c.importpathPrefix

	for k, v := range c.ruleProviders {
		clone.ruleProviders[k] = v
	}
	for k, v := range c.rules {
		clone.rules[k] = v
	}
	for k, v := range c.languages {
		clone.languages[k] = v
	}
	for k, v := range c.plugins {
		clone.plugins[k] = v
	}

	return clone
}

func (c *protoPackageConfig) RegisterRuleProvider(l label.Label, provider RuleProvider) {
	c.ruleProviders[l] = provider
}

func (c *protoPackageConfig) LookupRuleProvider(l label.Label) RuleProvider {
	return c.ruleProviders[l]
}

// mustGetProtoPlugin hardcodes the mapping between plugin names and
// implementations.  If the plugin is not recognized the function panics.
func (c *protoPackageConfig) mustGetProtoPlugin(rel, name string) ProtoPlugin {
	switch name {
	case "py_proto":
		return &PyProtoPlugin{}
	}
	panic(fmt.Sprintf(
		"%s/BUILD.bazel: invalid directive %s: unknown or unsupported proto_plugin %q",
		rel,
		protoPluginDirective,
		name))
}

// getExtensionConfig either inserts a new config into the map under the rosetta
// language name or replaces it with a clone.
func getExtensionConfig(exts map[string]interface{}) *protoPackageConfig {
	var cfg *protoPackageConfig
	if existingExt, ok := exts[languageName]; ok {
		cfg = existingExt.(*protoPackageConfig).Clone()
	} else {
		cfg = newProtoPackageConfig()
	}
	exts[languageName] = cfg
	return cfg
}

// protocKinds provides the aggregated list of KindInfo for the package.
func protocKinds() map[string]rule.KindInfo {
	// build of a list of all possible rules that we can see; delegate to the
	// rule implementations for the KindInfo.
	file := NewProtoFile("", "example.proto")
	protoLibraryRule := rule.NewRule("proto_library", "example")
	protoLibrary := &OtherProtoLibrary{rule: protoLibraryRule, files: []*ProtoFile{file}}
	pyProtoCompile := NewProtoRule(protoLibrary, "py", "proto", "compile")
	protoCompileTest := &ProtoCompileTest{pyProtoCompile}

	rules := []RuleProvider{
		// TODO(pcj): proto_library can apparently not be claimed as a kind by
		// two separate extensions. We will either have to take over this
		// responsibility or work with the proto_library targets that get
		// generated as it stands currently.
		pyProtoCompile,
		protoCompileTest,
		NewProtoDescriptorSet(protoLibrary),
		NewProtoRule(protoLibrary, "gogo", "proto", "compile"),
		NewProtoRule(protoLibrary, "gogofast", "proto", "compile"),
		NewProtoRule(protoLibrary, "gogofaster", "proto", "compile"),
		NewProtoRule(protoLibrary, "py", "proto", "library"),
		NewProtoRule(protoLibrary, "py", "grpc", "compile"),
		NewProtoRule(protoLibrary, "py_abc", "proto", "compile"),
		NewProtoRule(protoLibrary, "py_enum_choices", "proto", "compile"),
		NewProtoRule(protoLibrary, "py_rgrpc", "proto", "compile"),
		&PyLibrary{Lib: protoLibrary},
	}

	kinds := make(map[string]rule.KindInfo)
	for _, r := range rules {
		kinds[r.Kind()] = r.KindInfo()
	}

	return kinds
}

// protocLoads provides the aggregated list of LoadInfo for the package.
func protocLoads() []rule.LoadInfo {
	return []rule.LoadInfo{
		{
			Name: "@build_stack_rules_proto//rules:protoc.bzl",
			Symbols: []string{
				"protoc",
			},
		},
		{
			Name: "@build_stack_rules_proto//rules/proto:proto_compile_test.bzl",
			Symbols: []string{
				"proto_compile_test",
			},
		},
	}
}
