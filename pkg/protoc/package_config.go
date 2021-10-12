package protoc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	// RuleDirective is the directive for toggling rule generation.
	RuleDirective = "proto_rule"
	// LanguageDirective tells gazelle which languages a package should
	// produce and how it is configured.
	LanguageDirective = "proto_language"
	// PluginDirective created an association between proto_lang
	// and the label of a proto_plugin.
	PluginDirective = "proto_plugin"
	// importpathPrefixDirective is the same as 'gazelle:prefix'
	importpathPrefixDirective = "prefix"
)

// PackageConfig represents the config extension for the protobuf language.
type PackageConfig struct {
	// config is the parent gazelle config.
	config *config.Config
	// the gazelle:prefix for golang
	importpathPrefix string
	// configured languages for this package
	langs map[string]*LanguageConfig
	// exclude patterns for rules that should be skipped for this package.
	plugins map[string]*LanguagePluginConfig
	// exclude patterns for rules that should be skipped for this package.
	rules map[string]*LanguageRuleConfig
	// IMPORTANT! Adding new fields here?  Don't forget to copy it in the Clone
	// method!
}

// GetPackageConfig returns the associated package config.
func GetPackageConfig(config *config.Config) *PackageConfig {
	if cfg, ok := config.Exts["protobuf"].(*PackageConfig); ok {
		return cfg
	}
	return nil
}

// NewPackageConfig initializes a new PackageConfig.
func NewPackageConfig(config *config.Config) *PackageConfig {
	return &PackageConfig{
		config:  config,
		langs:   make(map[string]*LanguageConfig),
		plugins: make(map[string]*LanguagePluginConfig),
		rules:   make(map[string]*LanguageRuleConfig),
	}
}

// Plugin returns a readonly copy of the plugin configuration having the given
// name. If the plugin is not known the bool return arg is false.
func (c *PackageConfig) Plugin(name string) (LanguagePluginConfig, bool) {
	if c.plugins == nil {
		return LanguagePluginConfig{}, false
	}
	if plugin, ok := c.plugins[name]; ok {
		return *plugin, true
	} else {
		return LanguagePluginConfig{}, false
	}
}

// Clone copies this config to a new one.
func (c *PackageConfig) Clone() *PackageConfig {
	clone := NewPackageConfig(c.config)
	clone.importpathPrefix = c.importpathPrefix

	for k, v := range c.rules {
		clone.rules[k] = v.clone()
	}
	for k, v := range c.langs {
		clone.langs[k] = v.clone()
	}
	for k, v := range c.plugins {
		clone.plugins[k] = v.clone()
	}

	return clone
}

// ParseDirectives is called in each directory visited by gazelle.  The relative
// directory name is given by 'rel' and the list of directives in the BUILD file
// are specified by 'directives'.
func (c *PackageConfig) ParseDirectives(rel string, directives []rule.Directive) (err error) {
	for _, d := range directives {
		switch d.Key {
		case importpathPrefixDirective:
			err = c.parsePrefixDirective(d)
		case PluginDirective:
			err = c.parsePluginDirective(d)
		case RuleDirective:
			err = c.parseRuleDirective(d)
		case LanguageDirective:
			err = c.parseLanguageDirective(d)
		}
		if err != nil {
			return fmt.Errorf("parse %v: %w", d, err)
		}
	}
	return
}

func (c *PackageConfig) parsePrefixDirective(d rule.Directive) error {
	c.importpathPrefix = strings.TrimSpace(d.Value)
	return nil
}

func (c *PackageConfig) parseLanguageDirective(d rule.Directive) error {
	fields := strings.Fields(d.Value)
	if len(fields) != 3 {
		return fmt.Errorf("invalid directive %v: expected three fields, got %d", d, len(fields))
	}
	name, param, value := fields[0], fields[1], fields[2]
	lang, ok := c.langs[name]
	if !ok {
		lang = newLanguageConfig(name)
		c.langs[name] = lang
	}
	return lang.parseDirective(c, name, param, value)
}

func (c *PackageConfig) parsePluginDirective(d rule.Directive) error {
	fields := strings.Fields(d.Value)
	if len(fields) != 3 {
		return fmt.Errorf("invalid directive %v: expected three fields, got %d", d, len(fields))
	}
	name, param, value := fields[0], fields[1], fields[2]
	plugin, err := c.getOrCreateLanguagePluginConfig(name)
	if err != nil {
		return fmt.Errorf("invalid proto_plugin directive %+v: %w", d, err)
	}
	return plugin.parseDirective(c, name, param, value)
}

func (c *PackageConfig) parseRuleDirective(d rule.Directive) error {
	fields := strings.Fields(d.Value)
	if len(fields) < 3 {
		return fmt.Errorf("invalid directive %v: expected three or more fields, got %d", d, len(fields))
	}
	name, param, value := fields[0], fields[1], strings.Join(fields[2:], " ")
	r, err := c.getOrCreateLanguageRuleConfig(c.config, name)
	if err != nil {
		return fmt.Errorf("invalid proto_rule directive %+v: %w", d, err)
	}
	return r.parseDirective(c, name, param, value)
}

func (c *PackageConfig) getOrCreateLanguagePluginConfig(name string) (*LanguagePluginConfig, error) {
	plugin, ok := c.plugins[name]
	if !ok {
		plugin = newLanguagePluginConfig(name)
		c.plugins[name] = plugin
	}
	return plugin, nil
}

func (c *PackageConfig) getOrCreateLanguageRuleConfig(config *config.Config, name string) (*LanguageRuleConfig, error) {
	r, ok := c.rules[name]
	if !ok {
		r = NewLanguageRuleConfig(config, name)
		r.Implementation = name
		c.rules[name] = r
	}
	return r, nil
}

// configuredLangs returns a determinstic ordered list of configured
// langs
func (c *PackageConfig) configuredLangs() []*LanguageConfig {
	names := make([]string, 0)
	for name := range c.langs {
		names = append(names, name)
	}
	sort.Strings(names)
	langs := make([]*LanguageConfig, 0)
	for _, name := range names {
		langs = append(langs, c.langs[name])
	}
	return langs
}

func (c *PackageConfig) LoadYConfig(y *YConfig) error {
	for _, plugin := range y.Plugin {
		if err := c.loadYPlugin(plugin); err != nil {
			return err
		}
	}
	for _, rule := range y.Rule {
		if err := c.loadYRule(rule); err != nil {
			return err
		}
	}
	for _, lang := range y.Language {
		if err := c.loadYLanguage(lang); err != nil {
			return err
		}
	}
	return nil
}

func (c *PackageConfig) loadYPlugin(y *YPlugin) error {
	if y.Name == "" {
		return fmt.Errorf("yaml plugin name missing in: %+v", y)
	}
	plugin, err := c.getOrCreateLanguagePluginConfig(y.Name)
	if err != nil {
		return err
	}
	return plugin.fromYAML(y)
}

func (c *PackageConfig) loadYRule(y *YRule) error {
	if y.Name == "" {
		return fmt.Errorf("yaml rule name missing in: %+v", y)
	}
	rule, err := c.getOrCreateLanguageRuleConfig(c.config, y.Name)
	if err != nil {
		return err
	}
	return rule.fromYAML(y)
}

func (c *PackageConfig) loadYLanguage(y *YLanguage) error {
	if y.Name == "" {
		return fmt.Errorf("yaml language name missing in: %+v", y)
	}
	lang, ok := c.langs[y.Name]
	if !ok {
		lang = newLanguageConfig(y.Name)
		c.langs[y.Name] = lang
	}
	return lang.fromYAML(y)
}
