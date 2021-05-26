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
	Plugins map[string]bool
	Rules   map[string]bool
}

// newLanguageConfig constructs a new language configuration having the
// given name and list of default rule configurations.
func newLanguageConfig(name string, rules ...*LanguageRuleConfig) *LanguageConfig {
	c := &LanguageConfig{
		Name:    name,
		Rules:   make(map[string]bool),
		Plugins: make(map[string]bool),
		Enabled: true,
	}
	for _, rule := range rules {
		c.Rules[rule.Name] = true
	}
	return c
}

func (c *LanguageConfig) clone() *LanguageConfig {
	clone := newLanguageConfig(c.Name)
	clone.Enabled = c.Enabled
	for k, v := range c.Plugins {
		clone.Plugins[k] = v
	}
	for k, v := range c.Rules {
		clone.Rules[k] = v
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
		c.Plugins[value] = intent.Want
	case "rule":
		c.Rules[value] = intent.Want
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
