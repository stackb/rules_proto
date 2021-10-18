package protobuf

import (
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

func resolveOverrideRule(rel string, overrideRule *rule.Rule, resolver protoc.ImportResolver) {

	libs := overrideRule.PrivateAttr(protoLibrariesRuleKey).([]protoc.ProtoLibrary)
	if len(libs) == 0 {
		// printf("skipping resolve (no private libs)")
		return
	}

	for _, lib := range libs {
		r := lib.Rule()

		// filter out go_googleapis dependencies and re-resolve them anew.
		keep := make([]label.Label, 0)

		for _, dep := range r.AttrStrings("deps") {
			lbl, _ := label.Parse(dep)
			// printf("override resolve //%s:%s dep %v", rel, r.Name(), lbl)
			if lbl.Repo == "go_googleapis" {
				continue
			}
			if lbl.Relative {
				// relative labels will be repopulated via resolution (below)
				continue
			}
			keep = append(keep, lbl)
		}

		imports := r.PrivateAttr(config.GazelleImportsKey)
		if imps, ok := imports.([]string); ok {
			for _, imp := range imps {
				result := resolver.Resolve("proto", "proto", imp)
				// printf("go_googleapis resolve imports result: %+v", result)
				if len(result) > 0 {
					first := result[0]
					keep = append(keep, first.Label)
					// printf("go_googleapis resolve imports HIT", first.Label)
				} else {
					// printf("go_googleapis resolve imports MISS %s: %+v", imp, resolver)
				}
			}
		}

		if len(keep) > 0 {
			ss := make([]string, len(keep))
			for i, lbl := range keep {
				ss[i] = lbl.Rel("", rel).String()
			}
			r.SetAttr("deps", protoc.DeduplicateAndSort(ss))
		}
	}

	overrideRule.Delete()
}
