package protoc

import (
	"log"

	"github.com/bazelbuild/bazel-gazelle/label"
)

// ProtoCompileLanguageName is the name of the proto_compile language
// implementation.
const ProtoCompileLanguageName = "proto_compile"

func init() {
	MustRegisterProtoLanguage(ProtoCompileLanguageName, &ProtoCompileLanguage{})
}

// ProtoCompileLanguage implements a ProtoLanguage for that adds in
// *_proto_compile targets along with an associated test.
type ProtoCompileLanguage struct {
	Name string
	Cfg  *protoPackageConfig
}

// GenerateRules implements the ProtoLanguage interface.
func (s *ProtoCompileLanguage) GenerateRules(rel string, cfg *protoPackageConfig, libs []ProtoLibrary) []RuleProvider {
	rules := make([]RuleProvider, 0)

	log.Printf(ProtoCompileLanguageName+": Generate Rules: %+v", cfg)

	for _, lib := range libs {
		labels := make([]label.Label, 0)
		generatedSrcs := make([]string, 0)

		log.Printf(ProtoCompileLanguageName+": Iterate plugins %d", len(cfg.plugins))

		for name, p := range cfg.plugins {
			log.Printf(name + ": apply plugin?")

			if !p.Implementation.ShouldApply(rel, cfg, lib) {
				continue
			}
			labels = append(labels, p.Label)
			srcs := p.Implementation.GeneratedSrcs(rel, cfg, lib)
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
