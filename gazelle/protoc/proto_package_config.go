package protoc

import (
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
	// // protoLanguageConfigDirective configures a proto_language.
	// protoLanguageConfigDirective = "proto_language_config"
	// protoPluginDirective created an association between proto_language
	// and the label of a proto_plugin.
	protoPluginDirective = "proto_plugin"
	// prefixDirective is the same as 'gazelle:prefix'
	importpathPrefixDirective = "prefix"
)

// protoPackageConfig represents the config extension for the rosetta language.
type protoPackageConfig struct {
	// the gazelle:prefix for golang
	importpathPrefix string
	// configured languages for this package
	languages map[string]*ProtoLanguageConfig
	// exclude patterns for rules that should be skipped for this package.
	plugins map[string]*ProtoPluginConfig
	// ruleProviders is a mapping from label -> the provider that produced
	// the rule. we save this in the config such that we can retrieve the
	// association later in the resolve step.
	ruleProviders map[label.Label]RuleProvider
	// exclude patterns for rules that should be skipped for this package.
	rules map[string]*protoRuleConfig
	// IMPORTANT! Adding new fields here?  Don't forget to copy it in the Clone
	// method!
}

// newProtoPackageConfig initializes a new protoPackageConfig.
func newProtoPackageConfig() *protoPackageConfig {
	return &protoPackageConfig{
		languages:     make(map[string]*ProtoLanguageConfig),
		plugins:       make(map[string]*ProtoPluginConfig),
		ruleProviders: make(map[label.Label]RuleProvider),
		rules:         make(map[string]*protoRuleConfig),
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
				log.Fatalf(
					"invalid directive %s: expected {PLUGIN_NAME} {LABEL}, got %q",
					d.Key, d.Value)
			}
			pluginName := fields[0]
			pluginLabel := fields[1]
			l, err := label.Parse(pluginLabel)
			if err != nil {
				log.Fatalf(
					"invalid plugin label: %s", pluginLabel)
			}
			plugin, err := LookupProtoPlugin(pluginName)
			if err != nil {
				log.Fatalf("error while processing gazelle directive %q: %v", pluginName, err)
			}
			c.plugins[pluginName] = &ProtoPluginConfig{Label: l, Name: pluginName, Implementation: plugin}
			log.Printf("Added proto_plugin: %s -> %v", pluginName, l)
		case protoRuleDirective:
			pattern := strings.TrimSpace(d.Value)
			negative := strings.HasPrefix(pattern, "-")
			positive := strings.HasPrefix(pattern, "+")
			if negative || positive {
				pattern = pattern[1:]
			}
			if _, err := path.Match(pattern, ""); err == path.ErrBadPattern {
				log.Fatalf("invalid directive %s: bad match pattern %q", d.Key, pattern)
			}
			cfg := c.rules[pattern]
			if cfg == nil {
				cfg = &protoRuleConfig{pattern: pattern}
				c.rules[pattern] = cfg
			}
			cfg.exclude = negative
		// case protoLanguageDirective:
		// 	name := strings.TrimSpace(d.Value)
		// 	negative := strings.HasPrefix(name, "-")
		// 	positive := strings.HasPrefix(name, "+")
		// 	if negative || positive {
		// 		name = name[1:]
		// 	}
		// 	lang, ok := c.languages[name]
		// 	if !ok {
		// 		log.Fatalf("invalid or unknown proto_language %q", name)
		// 	}
		// 	if negative {
		// 		lang.Enabled = false
		// 	} else {
		// 		lang.Enabled = true
		// 	}
		case protoLanguageDirective:
			fields := strings.Fields(d.Value)
			if len(fields) != 3 {
				log.Fatalf("bad directive %v: expected three fields, got %d", d, len(fields))
			}
			name, param, value := fields[0], fields[1], fields[2]
			lang, ok := c.languages[name]
			if !ok {
				lang = &ProtoLanguageConfig{Name: name}
				lang.Implementation = MustLookupProtoLanguage(name)
			}
			lang.MustParseDirective(c, d.Key, param, value)
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
func (c *protoPackageConfig) Languages() []*ProtoLanguageConfig {
	names := make([]string, 0)
	for name := range c.languages {
		names = append(names, name)
	}
	sort.Strings(names)
	langs := make([]*ProtoLanguageConfig, 0)
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
	// file := NewProtoFile("", "example.proto")
	// protoLibraryRule := rule.NewRule("proto_library", "example")
	// protoLibrary := &OtherProtoLibrary{rule: protoLibraryRule, files: []*ProtoFile{file}}
	// pyProtoCompile := NewProtoRule(protoLibrary, "py", "proto", "compile")
	// protoCompileTest := &ProtoCompileTest{pyProtoCompile}

	rules := []RuleProvider{
		// TODO(pcj): proto_library can apparently not be claimed as a kind by
		// two separate extensions. We will either have to take over this
		// responsibility or work with the proto_library targets that get
		// generated as it stands currently.
		// pyProtoCompile,
		// protoCompileTest,
		// NewProtoDescriptorSet(protoLibrary),
		// NewProtoRule(protoLibrary, "gogo", "proto", "compile"),
		// NewProtoRule(protoLibrary, "gogofast", "proto", "compile"),
		// NewProtoRule(protoLibrary, "gogofaster", "proto", "compile"),
		// NewProtoRule(protoLibrary, "py", "proto", "library"),
		// NewProtoRule(protoLibrary, "py", "grpc", "compile"),
		// NewProtoRule(protoLibrary, "py_abc", "proto", "compile"),
		// NewProtoRule(protoLibrary, "py_enum_choices", "proto", "compile"),
		// NewProtoRule(protoLibrary, "py_rgrpc", "proto", "compile"),
		// &PyLibrary{Lib: protoLibrary},
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
