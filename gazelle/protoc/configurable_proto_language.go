package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/label"
)

// DirectableProtoLanguage implements a ProtoLanguage that can be configured via
// gazelle directives.
type DirectableProtoLanguage struct {
	name     string
	load     string
	kind     string
	disabled bool
	pattern  string
	plugins  []*protoPluginConfig
}

// GenerateRules implements the ProtoLanguage interface.
func (s *DirectableProtoLanguage) GenerateRules(
	rel string,
	cfg *protoPackageConfig,
	libs []ProtoLibrary,
) []RuleProvider {
	rules := make([]RuleProvider, 0)

	for _, lib := range libs {
		labels := make([]label.Label, 0)
		generatedSrcs := make([]string, 0)

		for _, p := range s.plugins {
			if !p.Plugin.ShouldApply(rel, cfg, lib) {
				continue
			}
			labels = append(labels, p.Label)
			srcs := p.Plugin.GeneratedSrcs(rel, cfg, lib)
			if len(srcs) > 0 {
				generatedSrcs = append(generatedSrcs, srcs...)
			}
		}

		if len(labels) > 0 {
			rules = append(rules, NewProtoCompileRule(lib, labels, generatedSrcs))
		}

	}

	return rules
}
