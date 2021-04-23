package protoc

import (
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
		gensrcs := make(map[string][]string)
		options := make(map[string][]string)
		mappings := make(map[string]map[string]string)
		outs := make(map[string]string)

		for _, plugin := range p.Plugins {
			if !plugin.Implementation.ShouldApply(rel, c, lib) {
				continue
			}
			labels = append(labels, plugin.Label)

			gensrcs[plugin.Label.String()] =
				plugin.Implementation.GeneratedSrcs(rel, c, lib)

			if provider, ok := plugin.Implementation.(PluginOptionsProvider); ok {
				opts := provider.GeneratedOptions(rel, c, lib)
				if len(opts) > 0 {
					options[plugin.Label.String()] = opts
				}
			}
			if provider, ok := plugin.Implementation.(PluginMappingsProvider); ok {
				m := provider.GeneratedMappings(rel, c, lib)
				if len(m) > 0 {
					mappings[plugin.Label.String()] = m
				}
			}
			if provider, ok := plugin.Implementation.(PluginOutProvider); ok {
				out := provider.GeneratedOut(rel, c, lib)
				if out != "" {
					outs[plugin.Label.String()] = out
				}
			}
		}

		if len(labels) > 0 {
			rules = append(rules, NewProtoCompileRule(
				rel,
				p.Name,
				lib,
				labels,
				gensrcs,
				options,
				mappings,
				outs,
			))
		}
	}

	return rules
}
