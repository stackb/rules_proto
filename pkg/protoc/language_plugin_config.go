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
	// Deps is a set of dep labels.  For example, consider a plugin config for
	// 'protoc-gen-go'.  That plugin produces .go files.  The 'Deps' field can
	// be used by downstream Rule implmentations to gather necessary
	// dependencies for the .go files produced by that plugin.
	Deps map[string]bool
	// Enabled flag
	Enabled bool
}

func newLanguagePluginConfig(name string) *LanguagePluginConfig {
	return &LanguagePluginConfig{
		Name:    name,
		Options: make(map[string]bool),
		Deps:    make(map[string]bool),
		Enabled: true,
	}
}

// GetOptions returns the sorted list of options with positive intent.
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

// GetDeps returns the sorted list of deps with positive intent.
func (c *LanguagePluginConfig) GetDeps() []string {
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

func (c *LanguagePluginConfig) clone() *LanguagePluginConfig {
	clone := newLanguagePluginConfig(c.Name)
	clone.Label = c.Label
	clone.Implementation = c.Implementation
	clone.Enabled = c.Enabled
	for k, v := range c.Options {
		clone.Options[k] = v
	}
	for k, v := range c.Deps {
		clone.Deps[k] = v
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
		c.Options[value] = intent.Want
	case "deps", "dep":
		c.Deps[value] = intent.Want
	default:
		return fmt.Errorf("unknown parameter %q", intent.Value)
	}

	return nil
}

// fromYAML loads configuration from the yaml plugin confug.
func (c *LanguagePluginConfig) fromYAML(y *YPlugin) error {
	if c.Name != y.Name {
		return fmt.Errorf("yaml plugin mismatch: want %q got %q", c.Name, y.Name)
	}
	c.Implementation = y.Implementation
	// only true intent is supported via yaml
	for _, option := range y.Option {
		c.Options[option] = true
	}
	for _, dep := range y.Dep {
		c.Deps[dep] = true
	}
	if y.Label != "" {
		l, err := label.Parse(y.Label)
		if err != nil {
			return fmt.Errorf("%s label parse error %w", y.Name, err)
		}
		c.Label = l
	}
	if y.Enabled != nil {
		c.Enabled = *y.Enabled
	} else {
		c.Enabled = true
	}
	return nil
}
