package protoc

import (
	gproto "github.com/bazelbuild/bazel-gazelle/language/proto"
)

// PluginContext represents the environment available to the plugin when
// invoked.
type PluginContext struct {
	// Rel is the relative path of the package.
	Rel string
	// ProtoLibrary is the proto_library under observation.
	ProtoLibrary ProtoLibrary
	// PackageConfig is the configuration for the package.
	PackageConfig PackageConfig
	// PluginConfig is the configuration object associated with the plugin.
	PluginConfig LanguagePluginConfig
	// Plugin is a reference to the plugin implementation
	Plugin Plugin
}

// IsProtoFileMode reports whether the surrounding bazel package is using
// the standard gazelle proto language's `file` mode (`# gazelle:proto file`),
// which generates one proto_library per .proto file. Plugins use this to
// switch between per-package and per-file output schemes.
func (ctx *PluginContext) IsProtoFileMode() bool {
	if ctx == nil || ctx.PackageConfig.Config == nil {
		return false
	}
	cfg := gproto.GetProtoConfig(ctx.PackageConfig.Config)
	if cfg == nil {
		return false
	}
	return cfg.Mode == gproto.FileMode
}

// type PluginContextResolver func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label)
