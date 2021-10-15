package protoc

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// PluginConfiguration represents the configuration of a protoc plugin and the
// sources & source mappings that are expected to be produced.
type PluginConfiguration struct {
	// Config is the associated plugin configuration.
	Config *LanguagePluginConfig
	// Label is the bazel label for the corresponding proto_plugin rule.
	Label label.Label
	// Mappings is a dictionary that maps filenames listed in Outputs to
	// 'Out'-relative filepaths.  This is used when the plugin writes to a
	// location outside the bazel package and needs to be relocated (copied) to
	// the Output location.
	Mappings map[string]string
	// Options is the list of options that the plugin expects
	Options []string
	// Out is the output directory the plugin is predicted to write to
	Out string
	// Outputs is the list of output files the plugin generates
	Outputs []string
	// Plugin the Plugin implementation that created the configuration.
	Plugin Plugin
}

// GetPluginLabels returns the list of labels strings for a list of plugins.
func GetPluginLabels(plugins []*PluginConfiguration) []string {
	labels := make([]string, len(plugins))
	for i, plugin := range plugins {
		labels[i] = plugin.Label.String()
	}
	sort.Strings(labels)
	return labels
}

// GetPluginOptions returns the list of options by plugin.
func GetPluginOptions(plugins []*PluginConfiguration, r *rule.Rule, from label.Label) map[string][]string {
	options := make(map[string][]string)
	for _, cfg := range plugins {
		opts := cfg.Options
		if resolver, ok := cfg.Plugin.(PluginOptionsResolver); ok {
			opts = resolver.ResolvePluginOptions(cfg, r, from)
		}
		if len(opts) == 0 {
			continue
		}
		sort.Strings(opts)
		options[cfg.Label.String()] = opts
	}
	return options
}

// GetPluginOuts returns the output location by plugin.
func GetPluginOuts(plugins []*PluginConfiguration) map[string]string {
	outs := make(map[string]string)
	for _, plugin := range plugins {
		if plugin.Out == "" {
			continue
		}
		outs[plugin.Label.String()] = plugin.Out
	}
	return outs
}
