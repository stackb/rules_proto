package protoc

import (
	"fmt"
	"strconv"
)

// ProtoLangConfig carries a set of configured Plugins and Rules that will
// contribute to a protoc invocation.
type ProtoLangConfig struct {
	Name    string
	Enabled bool
	Plugins map[string]*ProtoPluginConfig
	Rules   map[string]*ProtoRuleConfig
}

// newProtoLangConfig constructs a new language configuration having the
// given name and list of default rule configurations.
func newProtoLangConfig(name string, rules ...*ProtoRuleConfig) *ProtoLangConfig {
	c := &ProtoLangConfig{
		Name:    name,
		Rules:   make(map[string]*ProtoRuleConfig),
		Plugins: make(map[string]*ProtoPluginConfig),
		Enabled: true,
	}
	for _, rule := range rules {
		c.Rules[rule.Name] = rule.Clone()
	}
	return c
}

// Clone, well... it clones the config.
func (c *ProtoLangConfig) Clone() *ProtoLangConfig {
	clone := &ProtoLangConfig{
		Name:    c.Name,
		Enabled: c.Enabled,
		Plugins: make(map[string]*ProtoPluginConfig),
		Rules:   make(map[string]*ProtoRuleConfig),
	}
	for k, v := range c.Plugins {
		clone.Plugins[k] = v.Clone()
	}
	for k, v := range c.Rules {
		clone.Rules[k] = v.Clone()
	}
	return clone
}

// parseDirective parses the directive string or returns error.
func (c *ProtoLangConfig) parseDirective(cfg *ProtoPackageConfig, d, param, value string) error {
	intent := parseIntent(param)
	switch intent.Value {
	case "enabled", "enable":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("enable %s: %w", value, err)
		}
		c.Enabled = enabled
	case "plugin":
		plugin, ok := cfg.plugins[value]
		if intent.Negative {
			if ok {
				plugin.Enabled = false
			}
			return nil
		}
		if !ok {
			return fmt.Errorf("unknown plugin: %q", value)
		}
		c.Plugins[value] = plugin // TODO: clone here?
	case "rule":
		rule, ok := cfg.rules[value]
		if intent.Negative {
			if ok {
				rule.Enabled = false
			}
			return nil
		}
		if !ok {
			return fmt.Errorf("unknown rule: %q", value)
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
