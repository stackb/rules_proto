package protoc

import (
	"fmt"
	"sort"
	"strconv"
)

// ProtoRuleConfig carries metadata about a rule and its dependencies.
type ProtoRuleConfig struct {
	Deps           map[string]bool
	Enabled        bool
	Implementation ProtoRule
	Name           string
}

// newProtoRuleConfig returns a pointer to a new ProtoRule config with the
// 'Enabled' bit set to true.
func newProtoRuleConfig(name string) *ProtoRuleConfig {
	return &ProtoRuleConfig{
		Name:    name,
		Enabled: true,
		Deps:    make(map[string]bool),
	}
}

// Clone copies this config to a new one
func (c *ProtoRuleConfig) Clone() *ProtoRuleConfig {
	clone := &ProtoRuleConfig{
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

// GetDeps returns the sorted list of options
func (c *ProtoRuleConfig) GetDeps() []string {
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

// parseDirective parses the directive string or returns error.
func (c *ProtoRuleConfig) parseDirective(cfg *ProtoPackageConfig, d, param, value string) error {
	intent := parseIntent(param)
	switch intent.Value {
	case "dep":
		if intent.Negative {
			delete(c.Deps, value)
		} else {
			c.Deps[value] = true
		}
		return nil
	case "enabled":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("enabled %s: %w", value, err)
		}
		c.Enabled = enabled
		return nil
	default:
		return fmt.Errorf("unknown parameter %q", value)
	}
}
