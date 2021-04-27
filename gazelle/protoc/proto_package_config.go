package protoc

import (
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// protoRuleDirective is the directive for toggling rule generation.
	protoRuleDirective = "proto_rule"
	// protoLangDirective tells gazelle which languages a package should
	// produce and how it is configured.
	protoLangDirective = "proto_lang"
	// protoPluginDirective created an association between proto_lang
	// and the label of a proto_plugin.
	protoPluginDirective = "proto_plugin"
	// prefixDirective is the same as 'gazelle:prefix'
	importpathPrefixDirective = "prefix"
)

// ProtoPackageConfig represents the config extension for the rosetta language.
type ProtoPackageConfig struct {
	// the gazelle:prefix for golang
	importpathPrefix string
	// configured languages for this package
	langs map[string]*ProtoLangConfig
	// exclude patterns for rules that should be skipped for this package.
	plugins map[string]*ProtoPluginConfig
	// exclude patterns for rules that should be skipped for this package.
	rules map[string]*ProtoRuleConfig
	// IMPORTANT! Adding new fields here?  Don't forget to copy it in the Clone
	// method!
}

// NewProtoPackageConfig initializes a new ProtoPackageConfig.
func NewProtoPackageConfig() *ProtoPackageConfig {
	return &ProtoPackageConfig{
		langs:   make(map[string]*ProtoLangConfig),
		plugins: make(map[string]*ProtoPluginConfig),
		rules:   make(map[string]*ProtoRuleConfig),
	}
}

// Clone copies this config to a new one
func (c *ProtoPackageConfig) Clone() *ProtoPackageConfig {
	clone := NewProtoPackageConfig()
	clone.importpathPrefix = c.importpathPrefix

	for k, v := range c.rules {
		clone.rules[k] = v.Clone()
	}
	for k, v := range c.langs {
		clone.langs[k] = v.Clone()
	}
	for k, v := range c.plugins {
		clone.plugins[k] = v.Clone()
	}

	return clone
}

// ParseDirectives is called in each directory visited by gazelle.  The relative
// directory name is given by 'rel' and the list of directives in the BUILD file
// are specified by 'directives'.
func (c *ProtoPackageConfig) ParseDirectives(rel string, directives []rule.Directive) (err error) {
	for _, d := range directives {
		switch d.Key {
		case importpathPrefixDirective:
			err = c.parsePrefixDirective(d)
		case protoPluginDirective:
			err = c.parseProtoPluginDirective(d)
		case protoRuleDirective:
			err = c.parseProtoRuleDirective(d)
		case protoLangDirective:
			err = c.parseProtoLangDirective(d)
		}
		if err != nil {
			return fmt.Errorf("parse %v: %w", d, err)
		}
	}
	return
}

func (c *ProtoPackageConfig) parsePrefixDirective(d rule.Directive) error {
	c.importpathPrefix = strings.TrimSpace(d.Value)
	return nil
}

func (c *ProtoPackageConfig) parseProtoPluginDirective(d rule.Directive) error {
	fields := strings.Fields(d.Value)
	if len(fields) != 3 {
		return fmt.Errorf("invalid directive %v: expected three fields, got %d", d, len(fields))
	}
	name, param, value := fields[0], fields[1], fields[2]
	plugin, ok := c.plugins[name]
	if !ok {
		plugin = newProtoPluginConfig(name)
		impl, err := protoc.LookupProtoPlugin(name)
		if err == ErrUnknownPlugin {
			return fmt.Errorf("invalid proto_plugin directive: plugin not registered: %s", name)
		}
		plugin.Implementation = impl
		c.plugins[name] = plugin
	}
	return plugin.parseDirective(c, name, param, value)
}

func (c *ProtoPackageConfig) parseProtoRuleDirective(d rule.Directive) error {
	fields := strings.Fields(d.Value)
	if len(fields) != 3 {
		return fmt.Errorf("invalid directive %v: expected three fields, got %d", d, len(fields))
	}
	name, param, value := fields[0], fields[1], fields[2]
	r, ok := c.rules[name]
	if !ok {
		r = newProtoRuleConfig(name)
		impl, err := protoc.LookupProtoRule(name)
		if err == ErrUnknownRule {
			return fmt.Errorf("invalid proto_rule directive: rule not registered: %s", name)
		}
		r.Implementation = impl
		c.rules[name] = r
	}
	return r.parseDirective(c, name, param, value)
}

func (c *ProtoPackageConfig) parseProtoLangDirective(d rule.Directive) error {
	fields := strings.Fields(d.Value)
	if len(fields) != 3 {
		return fmt.Errorf("invalid directive %v: expected three fields, got %d", d, len(fields))
	}
	name, param, value := fields[0], fields[1], fields[2]
	lang, ok := c.langs[name]
	if !ok {
		// All language configurations get the proto_compile rule by default
		lang = newProtoLangConfig(name, c.mustGetOrCreateProtoRuleConfig(ProtoCompileName))
		c.langs[name] = lang
	}
	return lang.parseDirective(c, name, param, value)
}

func (c *ProtoPackageConfig) mustGetOrCreateProtoRuleConfig(name string) *ProtoRuleConfig {
	r, ok := c.rules[name]
	if !ok {
		r = newProtoRuleConfig(name)
		impl, err := protoc.LookupProtoRule(name)
		if err == ErrUnknownRule {
			log.Fatalf("required proto_rule not registered: %s", name)
		}
		r.Implementation = impl
		c.rules[name] = r
	}
	return r
}

// configuredLangs returns a determinstic ordered list of configured
// langs
func (c *ProtoPackageConfig) configuredLangs() []*ProtoLangConfig {
	names := make([]string, 0)
	for name := range c.langs {
		names = append(names, name)
	}
	sort.Strings(names)
	langs := make([]*ProtoLangConfig, 0)
	for _, name := range names {
		langs = append(langs, c.langs[name])
	}
	return langs
}

// // protocLoads provides the aggregated list of LoadInfo for the package.
// func protocLoads() []rule.LoadInfo {
// 	return []rule.LoadInfo{
// 		{
// 			Name: "@build_stack_rules_proto//rules:proto_compile.bzl",
// 			Symbols: []string{
// 				"proto_compile",
// 			},
// 		},
// 		{
// 			Name: "@build_stack_rules_proto//rules/proto:proto_compile_test.bzl",
// 			Symbols: []string{
// 				"proto_compile_test",
// 			},
// 		},
// 	}
// }
