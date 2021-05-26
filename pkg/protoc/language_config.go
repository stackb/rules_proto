package protoc

import (
	"fmt"
	"strconv"
)

// LanguageConfig carries a set of configured Plugins and Rules that will
// contribute to a protoc invocation.
type LanguageConfig struct {
	Name    string
	Enabled bool
	Plugins map[string]*LanguagePluginConfig
	Rules   map[string]*LanguageRuleConfig
}

// newLanguageConfig constructs a new language configuration having the
// given name and list of default rule configurations.
func newLanguageConfig(name string, rules ...*LanguageRuleConfig) *LanguageConfig {
	c := &LanguageConfig{
		Name:    name,
		Rules:   make(map[string]*LanguageRuleConfig),
		Plugins: make(map[string]*LanguagePluginConfig),
		Enabled: true,
	}
	for _, rule := range rules {
		c.Rules[rule.Name] = rule.clone()
	}
	return c
}

func (c *LanguageConfig) clone() *LanguageConfig {
	clone := &LanguageConfig{
		Name:    c.Name,
		Enabled: c.Enabled,
		Plugins: make(map[string]*LanguagePluginConfig),
		Rules:   make(map[string]*LanguageRuleConfig),
	}
	for k, v := range c.Plugins {
		clone.Plugins[k] = v.clone()
	}
	for k, v := range c.Rules {
		clone.Rules[k] = v.clone()
	}
	return clone
}

// parseDirective parses the directive string or returns error.
func (c *LanguageConfig) parseDirective(cfg *PackageConfig, d, param, value string) error {
	intent := parseIntent(param)
	switch intent.Value {
	case "enabled", "enable":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("enable %s: %w", value, err)
		}
		c.Enabled = enabled
	case "plugin":
		plugin, err := cfg.getOrCreateLanguagePluginConfig(value)
		if err != nil {
			return fmt.Errorf("could not bind plugin %q to language: %w", value, err)
		}
		plugin = plugin.clone()
		if !intent.Want {
			plugin.Enabled = false
		}
		c.Plugins[value] = plugin
	case "rule":
		rule, err := cfg.getOrCreateLanguageRuleConfig(value)
		if err != nil {
			return fmt.Errorf("could not bind rule %q to language: %w", value, err)
		}
		rule = rule.clone()
		if !intent.Want {
			rule.Enabled = false
		}
		c.Rules[value] = rule
	default:
		return fmt.Errorf("unknown parameter %q", value)
	}

	return nil
}

// case "match":
// 	if _, err := path.Match(value, ""); err == path.ErrBadPattern {
// 		return fmt.Errorf("match glob: %w", err)
// 	}
// 	s.Pattern = value
// 	return nil
