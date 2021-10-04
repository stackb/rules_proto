package protoc

import (
	"fmt"
	"sort"
	"strconv"
)

// LanguageConfig carries a set of configured Plugins and Rules that will
// contribute to a protoc invocation.
type LanguageConfig struct {
	Name    string
	Protoc  string
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

// GetEnabledRules filters the list of enabled Rules in a sorted manner.
func (c *LanguageConfig) GetRulesByIntent(intent bool) []string {
	names := make([]string, 0)
	for name, want := range c.Rules {
		if want != intent {
			continue
		}
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (c *LanguageConfig) clone() *LanguageConfig {
	clone := newLanguageConfig(c.Name)
	clone.Enabled = c.Enabled
	clone.Protoc = c.Protoc
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
	case "protoc":
		c.Protoc = value
	case "rule":
		c.Rules[value] = intent.Want
	default:
		return fmt.Errorf("unknown parameter %q", value)
	}

	return nil
}

// fromYAML loads configuration from the yaml rule confug.
func (c *LanguageConfig) fromYAML(y *YLanguage) error {
	if c.Name != y.Name {
		return fmt.Errorf("yaml language mismatch: want %q got %q", c.Name, y.Name)
	}
	for _, plugin := range y.Plugin {
		c.Plugins[plugin] = true
	}
	for _, rule := range y.Rule {
		c.Rules[rule] = true
	}
	if y.Enabled != nil {
		c.Enabled = *y.Enabled
	} else {
		c.Enabled = true
	}
	return nil
}
