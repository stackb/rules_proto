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

// clone copies this config to a new one
func (c *LanguageRuleConfig) clone() *LanguageRuleConfig {
	clone := newLanguageRuleConfig(c.Config, c.Name)
	clone.Enabled = c.Enabled
	clone.Implementation = c.Implementation
	for k, v := range c.Deps {
		clone.Deps[k] = v
	}
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
