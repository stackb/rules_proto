package protoc

import (
	"fmt"
	"sort"
	"strconv"
)

// LanguageRuleConfig carries metadata about a rule and its dependencies.
type LanguageRuleConfig struct {
	Deps           map[string]bool
	Enabled        bool
	Implementation LanguageRule
	Name           string
}

// newLanguageRuleConfig returns a pointer to a new LanguageRule config with the
// 'Enabled' bit set to true.
func newLanguageRuleConfig(name string) *LanguageRuleConfig {
	return &LanguageRuleConfig{
		Name:    name,
		Enabled: true,
		Deps:    make(map[string]bool),
	}
}

// GetDeps returns the sorted list of options
func (c *LanguageRuleConfig) GetDeps() []string {
	deps := make([]string, 0)
	for dep, want := range c.Deps {
		if !want {
			continue
		}
		deps = append(deps, dep)
	}
	sort.Strings(deps)
	return deps
}

// clone copies this config to a new one
func (c *LanguageRuleConfig) clone() *LanguageRuleConfig {
	clone := &LanguageRuleConfig{
		Deps:           c.Deps,
		Enabled:        c.Enabled,
		Implementation: c.Implementation,
		Name:           c.Name,
	}
	for k, v := range c.Deps {
		clone.Deps[k] = v
	}
	return clone
}

// parseDirective parses the directive string or returns error.
func (c *LanguageRuleConfig) parseDirective(cfg *PackageConfig, d, param, value string) error {
	intent := parseIntent(param)
	switch intent.Value {
	case "dep":
		if intent.Negative {
			delete(c.Deps, value)
		} else {
			c.Deps[value] = true
		}
	case "enabled":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("enabled %s: %w", value, err)
		}
		c.Enabled = enabled
	default:
		return fmt.Errorf("unknown parameter %q", value)
	}

	return nil
}
