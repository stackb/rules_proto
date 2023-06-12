package protobuf

import (
	"log"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	// protoOverrideRulesKey is used to stash a list of proto_library rules in a
	// private attr for later deps resolution.
	protoLibrariesRuleKey = "_proto_library_rules"
	// overrideKindName is the name of the kind
	overrideKindName = "proto_library_override"
	// debugOverrides is a developer-flag.
	debugOverrides = false
)

var overrideKind = rule.KindInfo{
	ResolveAttrs: map[string]bool{"deps": true},
}

func makeProtoOverrideRule(libs []protoc.ProtoLibrary) *rule.Rule {
	// This rule is *only* used to trigger a Resolve() callback such that we can
	// process the proto_library rules we've captured here; the rule itself is
	// always deleted from the file.
	overrideRule := rule.NewRule(overrideKindName, protoLibrariesRuleKey)
	overrideRule.SetPrivateAttr(protoLibrariesRuleKey, libs)
	return overrideRule
}

func resolveOverrideRule(c *config.Config, rel string, overrideRule *rule.Rule, resolver protoc.ImportResolver) {

	libs := overrideRule.PrivateAttr(protoLibrariesRuleKey).([]protoc.ProtoLibrary)
	if len(libs) == 0 {
		return
	}

	for _, lib := range libs {
		r := lib.Rule()

		// re-resolve dependencies.
		deps := make([]label.Label, 0)

		imports := r.PrivateAttr(config.GazelleImportsKey)
		if imps, ok := imports.([]string); ok {
			for _, imp := range imps {
				result := resolver.Resolve("proto", "proto", imp)
				if len(result) > 0 {
					first := result[0]
					deps = append(deps, first.Label)
					if debugOverrides {
						log.Println("go_googleapis resolve imports HIT", imp, first.Label)
					}
				} else {
					if debugOverrides {
						log.Printf("go_googleapis resolve imports MISS %s: %+v", imp, resolver)
					}
				}
			}
		}

		if len(deps) > 0 {
			ss := make([]string, len(deps))
			for i, lbl := range deps {
				ss[i] = lbl.Rel("", rel).String()
			}
			r.SetAttr("deps", protoc.DeduplicateAndSort(ss))
		}
	}

	overrideRule.Delete()
}
