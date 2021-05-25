package protoc

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
)

// PluginConfiguration represents the configuration of a protoc plugin
// and the sources & source mappings that are expected to be produced.
type PluginConfiguration struct {
	Skip     bool
	Name     string
	Label    label.Label
	Mappings map[string]string
	Options  []string
	Out      string
	Outputs  []string
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
func GetPluginOptions(plugins []*PluginConfiguration) map[string][]string {
	options := make(map[string][]string)
	for _, plugin := range plugins {
		if len(plugin.Options) == 0 {
			continue
		}
		opts := plugin.Options
		sort.Strings(opts)
		options[plugin.Label.String()] = opts
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
