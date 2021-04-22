package protoc

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
)

// ProtoCompileLanguage implements a ProtoLanguage for that adds in
// *_proto_compile targets along with an associated test.  It is the default
// implementation chosen if no existing other implementaiton is registered.
type ProtoCompileLanguage struct {
	Name string
}

// GenerateRules implements the ProtoLanguage interface.
func (s *ProtoCompileLanguage) GenerateRules(
	rel string,
	c *ProtoPackageConfig,
	p *ProtoLanguageConfig,
	libs []ProtoLibrary,
) []RuleProvider {
	rules := make([]RuleProvider, 0)

	for _, lib := range libs {
		labels := make([]label.Label, 0)
		generatedSrcs := make([]string, 0)
		generatedOptions := make(map[string][]string)
		generatedOuts := make(map[string]string)

		for _, plugin := range p.Plugins {
			if !plugin.Implementation.ShouldApply(rel, c, lib) {
				continue
			}
			labels = append(labels, plugin.Label)
			srcs := plugin.Implementation.GeneratedSrcs(rel, c, lib)
			if len(srcs) > 0 {
				generatedSrcs = append(generatedSrcs, srcs...)
			}
			if provider, ok := plugin.Implementation.(PluginOptionsProvider); ok {
				options := provider.GeneratedOptions(rel, c, lib)
				if len(options) > 0 {
					generatedOptions[plugin.Label.String()] = options
				}
			}
			if provider, ok := plugin.Implementation.(PluginOutProvider); ok {
				out := provider.GeneratedOut(rel, c, lib)
				if out != "" {
					generatedOuts[plugin.Label.String()] = out
				}
			}
			// log.Printf("plugin %s generated_srcs: %v", plugin.Name, srcs)
		}

		if len(labels) > 0 {
			sort.Strings(generatedSrcs)
			rules = append(rules, NewProtoCompileRule(
				p.Name,
				lib,
				labels,
				generatedSrcs,
				generatedOptions,
				generatedOuts,
			))
		}
	}

	return rules
}
