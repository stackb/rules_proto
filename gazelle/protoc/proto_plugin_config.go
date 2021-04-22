package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/label"
)

// ProtoPluginConfig associates metadata with a plugin implementation.
type ProtoPluginConfig struct {
	Label          label.Label
	Name           string
	Implementation ProtoPlugin
}

func (c *ProtoPluginConfig) Clone() *ProtoPluginConfig {
	return &ProtoPluginConfig{
		Label:          c.Label,
		Name:           c.Name,
		Implementation: c.Implementation,
	}
}
