package protoc

import (
	"log"

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

		for _, plugin := range p.Plugins {
			log.Printf("processing plugin %s", plugin.Name)
			if !plugin.Implementation.ShouldApply(rel, c, lib) {
				log.Printf("skipping plugin %s (should not apply)", plugin.Name)
				continue
			}
			labels = append(labels, plugin.Label)
			srcs := plugin.Implementation.GeneratedSrcs(rel, c, lib)
			if len(srcs) > 0 {
				generatedSrcs = append(generatedSrcs, srcs...)
			}
			log.Printf("plugin %s generated_srcs: %v", plugin.Name, srcs)
		}

		if len(labels) > 0 {
			rules = append(rules, NewProtoCompileRule(p.Name, lib, labels, generatedSrcs))
		}
	}

	return rules
}
