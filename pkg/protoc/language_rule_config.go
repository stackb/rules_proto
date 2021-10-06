package protoc

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/bazelbuild/bazel-gazelle/config"
)

// LanguageRuleConfig carries metadata about a rule and its dependencies.
type LanguageRuleConfig struct {
	// Config is the parent gazelle Config
	Config *config.Config
	// Deps is a mapping from label to +/- intent.
	Deps map[string]bool
	// Resolves is a mapping from resolve mapping spec to rewrite.  Negative
	// intent is represented by the empty rewrite value.
	Resolves []Rewrite
	// Enabled is a flag that marks language generation as enabled or not
	Enabled bool
	// Implementation is the registry identifier for the Rule
	Implementation string
	// Impl is the actual implementation
	Impl LanguageRule
	// Name is the name of the Rule config
	Name string
	// Deps is a mapping from visibility label to +/- intent.
	Visibility map[string]bool
}

// newLanguageRuleConfig returns a pointer to a new LanguageRule config with the
// 'Enabled' bit set to true.
func newLanguageRuleConfig(config *config.Config, name string) *LanguageRuleConfig {
	return &LanguageRuleConfig{
		Config:     config,
		Name:       name,
		Enabled:    true,
		Deps:       make(map[string]bool),
		Resolves:   make([]Rewrite, 0),
		Visibility: make(map[string]bool),
	}
}

// GetDeps returns the sorted list of dependencies
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

// GetRewrites returns a copy of the resolve mappings
func (c *LanguageRuleConfig) GetRewrites() []Rewrite {
	return c.Resolves[:]
}

// clone copies this config to a new one
func (c *LanguageRuleConfig) clone() *LanguageRuleConfig {
	clone := newLanguageRuleConfig(c.Config, c.Name)
	clone.Enabled = c.Enabled
	clone.Implementation = c.Implementation
	for k, v := range c.Deps {
		clone.Deps[k] = v
	}
	clone.Resolves = c.Resolves[:]
	for k, v := range c.Visibility {
		clone.Visibility[k] = v
	}
	return clone
}

// parseDirective parses the directive string or returns error.
func (c *LanguageRuleConfig) parseDirective(cfg *PackageConfig, d, param, value string) error {
	intent := parseIntent(param)
	switch intent.Value {
	case "dep", "deps":
		if intent.Want {
			c.Deps[value] = true
		} else {
			delete(c.Deps, value)
		}
	case "resolve":
		rw, err := ParseRewrite(value)
		if err != nil {
			return fmt.Errorf("invalid resolve rewrite %s: %w", value, err)
		}
		c.Resolves = append(c.Resolves, *rw)
	case "visibility":
		c.Visibility[value] = intent.Want
	case "implementation":
		c.Implementation = value
	case "enabled":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("enabled %s: %w", value, err)
		}
		c.Enabled = enabled
	default:
		return fmt.Errorf("unknown parameter %q", intent.Value)
	}

	return nil
}

// fromYAML loads configuration from the yaml rule confug.
func (c *LanguageRuleConfig) fromYAML(y *YRule) error {
	if c.Name != y.Name {
		return fmt.Errorf("yaml rule mismatch: want %q got %q", c.Name, y.Name)
	}
	c.Implementation = y.Implementation
	for _, dep := range y.Deps {
		c.Deps[dep] = true
	}
	for _, resolve := range y.Resolves {
		rw, err := ParseRewrite(resolve)
		if err != nil {
			return fmt.Errorf("invalid resolve rewrite %s: %w", resolve, err)
		}
		c.Resolves = append(c.Resolves, *rw)
	}
	for _, v := range y.Visibility {
		c.Visibility[v] = true
	}
	if y.Enabled != nil {
		c.Enabled = *y.Enabled
	} else {
		c.Enabled = true
	}
	return nil
}
