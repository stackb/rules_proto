package protoc

import (
	"path"
	"sort"
	"strings"
)

// ProtocConfiguration represents the complete configuration and source
// mappings.
type ProtocConfiguration struct {
	// The config for the p
	LanguageConfig *LanguageConfig
	// the workspace relative path of the BUILD file where this rule is being
	// generated.
	Rel string
	// the prefix for the rule (e.g. 'java')
	Prefix string
	// the library thar holds the proto files
	Library ProtoLibrary
	// the configuration for the plugins
	Plugins []*PluginConfiguration
	// The merged set of Source files for the compilations
	Outputs []string
	// The merged set of imports for the compilations
	Imports []string
	// The generated source mappings
	Mappings map[string]string
}

func newProtocConfiguration(lc *LanguageConfig, rel, prefix string, lib ProtoLibrary, plugins []*PluginConfiguration) *ProtocConfiguration {
	srcs, mappings := mergeSources(rel, plugins)
	imports := mergeImports(plugins)

	return &ProtocConfiguration{
		LanguageConfig: lc,
		Rel:            rel,
		Prefix:         prefix,
		Library:        lib,
		Plugins:        plugins,
		Outputs:        srcs,
		Imports:        imports,
		Mappings:       mappings,
	}
}

func (c *ProtocConfiguration) GetPluginOutputs(name string) []string {
	for _, plugin := range c.Plugins {
		if plugin.Name == name {
			return plugin.Outputs
		}
	}
	return nil
}

// mergeSources computes the source files that are generated by the rule and any
// necessary mappings.
func mergeSources(rel string, plugins []*PluginConfiguration) ([]string, map[string]string) {
	srcs := make([]string, 0)
	mappings := make(map[string]string)

	for _, plugin := range plugins {
		// if plugin provided mappings for us, use those preferentially
		if len(plugin.Mappings) > 0 {
			srcs = append(srcs, plugin.Outputs...)

			for k, v := range plugin.Mappings {
				mappings[k] = v
			}
			continue
		}

		// otherwise, fallback to baseline method
		for _, filename := range plugin.Outputs {
			dir := path.Dir(filename)
			if dir == "." && rel == "" {
				dir = rel
			}
			base := path.Base(filename)
			if dir == rel {
				// no mapping required, just add to the srcs list
				srcs = append(srcs, strings.TrimPrefix(filename, rel+"/"))
			} else {
				// add the basename only to the srcs list and add a mapping.
				srcs = append(srcs, base)
				mappings[base] = filename
			}
		}
	}

	return srcs, mappings
}

// mergeImports computes the merged list of imports for the list of plugins.
func mergeImports(plugins []*PluginConfiguration) []string {
	imports := make([]string, 0)

	for _, plugin := range plugins {
		imports = append(imports, plugin.Imports...)
	}

	sort.Strings(imports)

	return imports
}
