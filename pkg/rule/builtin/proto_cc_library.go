package builtin

import (
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoCcLibraryRuleName   = "proto_cc_library"
	ProtoCcLibraryRuleSuffix = "_cc_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_cc_library", &protoCcLibrary{})
}

// protoCcLibrary implements LanguageRule for the 'proto_cc_library' rule from
// @rules_proto.
type protoCcLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoCcLibrary) Name() string {
	return ProtoCcLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoCcLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs":       true,
			"hdrs":       true,
			"deps":       true,
			"visibility": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoCcLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/cc:proto_cc_library.bzl",
		Symbols: []string{ProtoCcLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoCcLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("builtin:cpp")
	if len(outputs) == 0 {
		return nil
	}
	return &CcLibraryRule{
		KindName:       ProtoCcLibraryRuleName,
		RuleNameSuffix: ProtoCcLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver: func(impl *CcLibraryRule, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
			deps := impl.Deps()

			for _, d := range impl.Config.Library.Deps() {
				if strings.HasPrefix(d, "@com_google_protobuf//") {
					continue
				}
				d = strings.TrimSuffix(d, "_proto")
				deps = append(deps, d+ProtoCcLibraryRuleSuffix)
			}

			if len(deps) > 0 {
				r.SetAttr("deps", deps)
			}
		},
	}
}
