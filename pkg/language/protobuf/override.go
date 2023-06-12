package protobuf

import (
	"fmt"
	"log"
	"os"

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
	debugOverrides = true
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

		// filter out go_googleapis dependencies and re-resolve them anew.
		keep := make([]label.Label, 0)

		// log.Printf("override resolve //%s:%s", rel, r.Name())

		debug := rel == "google/api/servicecontrol/v1" && r.Name() == "log_entry_proto"

		// if debug {
		// 	printRules(r)
		// }
		for _, dep := range r.AttrStrings("deps") {
			lbl, _ := label.Parse(dep)
			if debug {
				log.Printf("override resolve dep rel=%s, r.Name()=%s, lbl=%s, c.RepoName=%s", rel, r.Name(), lbl, c.RepoName)
			}
			if lbl.Relative {
				// relative labels will be repopulated via resolution (below)
				continue
			}
			// if lbl.Repo == "go_googleapis" || (lbl.Repo == "" && c.RepoName == "googleapis") {
			if lbl.Repo == "go_googleapis" {
				if debug {
					log.Printf("override resolve //%s:%s dep %v", rel, r.Name(), lbl)
				}
				continue
			}
			if lbl.Repo == "com_google_protobuf" {
				if debug {
					log.Printf("override resolve //%s:%s dep %v", rel, r.Name(), lbl)
				}
				continue
			}
			// keep = append(keep, lbl)
		}

		var hasLogSeverityProto bool
		imports := r.PrivateAttr(config.GazelleImportsKey)
		if imps, ok := imports.([]string); ok {
			for _, imp := range imps {
				if imp == "XXX google/logging/type/log_severity.proto" {
					hasLogSeverityProto = true
				}
				result := resolver.Resolve("proto", "proto", imp)
				if len(result) > 0 {
					first := result[0]
					keep = append(keep, first.Label)
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

		if hasLogSeverityProto {
			log.Println("deps before:")
			printRules(r)
			log.Println("deps keep:", keep)
		}

		if len(keep) > 0 {
			ss := make([]string, len(keep))
			for i, lbl := range keep {
				ss[i] = lbl.Rel("", rel).String()
			}
			r.SetAttr("deps", protoc.DeduplicateAndSort(ss))
		}
		if hasLogSeverityProto {
			log.Println("deps after:")
			printRules(r)
		}
	}

	overrideRule.Delete()
}

func printRules(rules ...*rule.Rule) {
	file := rule.EmptyFile("", "")
	for _, r := range rules {
		r.Insert(file)
	}
	fmt.Fprintln(os.Stderr, string(file.Format()))
}
