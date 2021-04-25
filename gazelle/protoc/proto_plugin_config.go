package protoc

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/bazelbuild/bazel-gazelle/label"
)

// ProtoPluginConfig associates metadata with a plugin implementation.
type ProtoPluginConfig struct {
	// Name of the config, for the sake of configuration
	Name string
	// Label is the bazel label of the ProtoPluginInfo provider
	Label label.Label
	// Tool is the bazel label for the binary tool
	Tool label.Label
	// Options is a set of option strings.
	Options map[string]bool
	// Implementation is the ProtoPlugin implementation registered.
	Implementation ProtoPlugin
	// Enabled flag
	Enabled bool
}

func newProtoPluginConfig(name string) *ProtoPluginConfig {
	return &ProtoPluginConfig{
		Name:    name,
		Options: make(map[string]bool),
		Enabled: true,
	}
}

func (c *ProtoPluginConfig) Clone() *ProtoPluginConfig {
	// if c.Label.Name == "" {
	// 	panic(fmt.Sprintf("proto_plugin label is mandatory, please declare it (e.g. '#gazelle:proto_plugin label @foo//bar'): %+v", c))
	// }
	clone := &ProtoPluginConfig{
		Label:          c.Label,
		Name:           c.Name,
		Tool:           c.Tool,
		Implementation: c.Implementation,
		Enabled:        c.Enabled,
	}
	for k, v := range c.Options {
		clone.Options[k] = v
	}
	return clone
}

// GetOptions returns the sorted list of options
func (c *ProtoPluginConfig) GetOptions() []string {
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

// parseDirective parses the directive string or returns error.
func (c *ProtoPluginConfig) parseDirective(cfg *ProtoPackageConfig, d, param, value string) error {
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
	case "tool":
		l, err := label.Parse(value)
		if err != nil {
			return fmt.Errorf("tool %q: %w", value, err)
		}
		c.Tool = l
	case "option":
		if intent.Negative {
			delete(c.Options, value)
		} else {
			c.Options[value] = true
		}
	default:
		return fmt.Errorf("invalid directive %q: unknown parameter %q", d, value)
	}

	return nil
}
