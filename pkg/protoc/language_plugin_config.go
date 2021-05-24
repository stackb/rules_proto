package protoc

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/bazelbuild/bazel-gazelle/label"
)

// LanguagePluginConfig associates metadata with a plugin implementation.
type LanguagePluginConfig struct {
	// Name is the identifier for the configuration object
	Name string
	// Implementation is the identifier for the implementation
	Implementation string
	// Label is the bazel label of the PluginInfo provider
	Label label.Label
	// Options is a set of option strings.
	Options map[string]bool
	// Enabled flag
	Enabled bool
}

func newLanguagePluginConfig(name string) *LanguagePluginConfig {
	return &LanguagePluginConfig{
		Name:    name,
		Options: make(map[string]bool),
		Enabled: true,
	}
}

// GetOptions returns the sorted list of options
func (c *LanguagePluginConfig) GetOptions() []string {
	opts := make([]string, 0)
	for opt, want := range c.Options {
		if !want {
			continue
		}
		opts = append(opts, opt)
	}
	sort.Strings(opts)
	return opts
}

func (c *LanguagePluginConfig) clone() *LanguagePluginConfig {
	clone := newLanguagePluginConfig(c.Name)
	clone.Label = c.Label
	clone.Implementation = c.Implementation
	clone.Enabled = c.Enabled
	for k, v := range c.Options {
		clone.Options[k] = v
	}
	return clone
}

// parseDirective parses the directive string or returns error.
func (c *LanguagePluginConfig) parseDirective(cfg *PackageConfig, d, param, value string) error {
	intent := parseIntent(param)
	switch intent.Value {
	case "enabled", "enable":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("enabled %s: %w", value, err)
		}
		c.Enabled = enabled
	case "label":
		l, err := label.Parse(value)
		if err != nil {
			return fmt.Errorf("label %q: %w", value, err)
		}
		c.Label = l
	case "implementation":
		c.Implementation = value
	case "option":
		c.Options[value] = !intent.Negative
	default:
		return fmt.Errorf("unknown parameter %q", intent.Value)
	}

	return nil
}
