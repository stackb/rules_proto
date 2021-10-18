package protoc

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

// type PluginContextResolver func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label)
