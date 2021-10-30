package protoc

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
)

// LanguageRuleConfig carries metadata about a rule and its dependencies.
type LanguageRuleConfig struct {
	// Config is the parent gazelle Config
	Config *config.Config
	// Deps is a mapping from label to +/- intent.
	Deps map[string]bool
	// Attr is a mapping from string to intent map.
	Attrs map[string]map[string]bool
	// Options is a generic key -> value string mapping.  Various rule
	// implementations are free to document/interpret options in an
	// implementation-dependenct manner.
	Options map[string]bool
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

// NewLanguageRuleConfig returns a pointer to a new LanguageRule config with the
// 'Enabled' bit set to true.
func NewLanguageRuleConfig(config *config.Config, name string) *LanguageRuleConfig {
	return &LanguageRuleConfig{
		Config:     config,
		Name:       name,
		Enabled:    true,
		Attrs:      make(map[string]map[string]bool),
		Deps:       make(map[string]bool),
		Options:    make(map[string]bool),
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

// GetOptions returns the rule options.
func (c *LanguageRuleConfig) GetOptions() []string {
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

// GetAttr returns the positive-intent attr values under the given key.
func (c *LanguageRuleConfig) GetAttr(name string) []string {
	vals := make([]string, 0)
	for val, want := range c.Attrs[name] {
		if !want {
			continue
		}
		vals = append(vals, val)
	}
	sort.Strings(vals)
	return vals
}

// GetRewrites returns a copy of the resolve mappings
func (c *LanguageRuleConfig) GetRewrites() []Rewrite {
	return c.Resolves[:]
}

// clone copies this config to a new one
func (c *LanguageRuleConfig) clone() *LanguageRuleConfig {
	clone := NewLanguageRuleConfig(c.Config, c.Name)
	clone.Enabled = c.Enabled
	clone.Implementation = c.Implementation
	for name, vals := range c.Attrs {
		clone.Attrs[name] = make(map[string]bool)
		for k, v := range vals {
			clone.Attrs[name][k] = v
		}
	}
	for k, v := range c.Deps {
		clone.Deps[k] = v
	}
	for k, v := range c.Options {
		clone.Options[k] = v
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
	case "option":
		if intent.Want {
			c.Options[value] = true
		} else {
			delete(c.Options, value)
		}
	case "attr":
		kv := strings.Fields(value)
		if len(kv) == 0 {
			return fmt.Errorf("malformed attr (missing attr name and value) %q: expected form is 'gazelle:proto_rule {RULE_NAME} attr {ATTR_NAME} [+/-]{VALUE}'", value)
		}
		key := parseIntent(kv[0])

		if len(kv) == 1 {
			if intent.Want {
				return fmt.Errorf("malformed attr %q (missing named attr value): expected form is 'gazelle:proto_rule {RULE_NAME} attr {ATTR_NAME} [+/-]{VALUE}'", value)
			} else {
				delete(c.Attrs, key.Value)
				return nil
			}
		}

		val := strings.Join(kv[1:], " ")

		if intent.Want {
			values, ok := c.Attrs[key.Value]
			if !ok {
				values = make(map[string]bool)
				c.Attrs[key.Value] = values
			}
			values[val] = key.Want
		} else {
			delete(c.Attrs, key.Value)
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

func (c *LanguageRuleConfig) addRewrite(r Rewrite) {
	c.Resolves = append([]Rewrite{r}, c.Resolves...)
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
	for _, opt := range y.Option {
		c.Options[opt] = true
	}
	for _, resolve := range y.Resolves {
		rw, err := ParseRewrite(resolve)
		if err != nil {
			return fmt.Errorf("invalid resolve rewrite %s: %w", resolve, err)
		}
		c.addRewrite(*rw)
	}
	for _, v := range y.Visibility {
		c.Visibility[v] = true
	}
	if y.Enabled != nil {
		c.Enabled = *y.Enabled
	}
	return nil
}
