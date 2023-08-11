package rules_nodejs

import (
	"flag"
	"log"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoTsLibraryRuleName   = "proto_ts_library"
	ProtoTsLibraryRuleSuffix = "_ts_proto"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_ts_library", &protoTsLibrary{})
}

// protoTsLibrary implements LanguageRule for the 'proto_ts_library' rule from
// @rules_proto.
type protoTsLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoTsLibrary) Name() string {
	return ProtoTsLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoTsLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs":     true,
			"tsc":      true,
			"args":     true,
			"data":     true,
			"tsconfig": true,
			"out_dir":  true,
		},
		ResolveAttrs: map[string]bool{
			"deps": true,
		},
	}

}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoTsLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/ts:proto_ts_library.bzl",
		Symbols: []string{ProtoTsLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoTsLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	flags := parseProtoTsLibraryFlags(ProtoTsLibraryRuleName, cfg.GetOptions())

	outputs := make([]string, 0)
	for _, out := range pc.Outputs {
		if strings.HasSuffix(out, ".ts") {
			outputs = append(outputs, out)
		}
	}
	if len(outputs) == 0 {
		return nil
	}
	return &tsLibrary{
		flags:          flags,
		KindName:       ProtoTsLibraryRuleName,
		RuleNameSuffix: ProtoTsLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
	}
}

// tsLibrary implements RuleProvider for 'ts_library'-like rules.
type tsLibrary struct {
	flags          *protoTsLibraryFlags
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *tsLibrary) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *tsLibrary) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *tsLibrary) Srcs() []string {
	return s.Outputs
}

// Deps computes the deps list for the rule.
func (s *tsLibrary) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility provides visibility labels.
func (s *tsLibrary) Visibility() []string {
	return s.RuleConfig.GetVisibility()
}

// Rule implements part of the ruleProvider interface.
func (s *tsLibrary) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())

	deps := s.Deps()
	if len(deps) > 0 {
		newRule.SetAttr("deps", deps)
	}

	args := s.RuleConfig.GetAttr("args")
	if len(args) > 0 {
		newRule.SetAttr("args", args)
	}

	tsconfig := s.RuleConfig.GetAttr("tsconfig")
	if len(tsconfig) > 0 {
		if len(tsconfig) > 1 {
			log.Printf("warning (%s) found multiple entries for 'tsconfig', choosing last one: %v", s.Kind(), tsconfig)
		}
		newRule.SetAttr("tsconfig", tsconfig[len(tsconfig)-1])
	}

	outdir := s.RuleConfig.GetAttr("out_dir")
	if len(outdir) > 0 {
		if len(outdir) > 1 {
			log.Printf("warning (%s) found multiple entries for 'out_dir', choosing last one: %v", s.Kind(), outdir)
		}
		newRule.SetAttr("out_dir", outdir[len(outdir)-1])
	}

	if s.flags.includeProtoLibraryData {
		newRule.SetAttr("data", []string{s.Config.Library.Name()})
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *tsLibrary) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	if lib, ok := r.PrivateAttr(protoc.ProtoLibraryKey).(protoc.ProtoLibrary); ok {
		return protoc.ProtoLibraryImportSpecsForKind(r.Kind(), lib)
	}
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *tsLibrary) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	protoc.ResolveDepsAttr("deps", false /* resolve wkts */)(c, ix, r, imports, from)
}

// protoTsLibraryFlags represents the parsed flag configuration for the
// proto_ts_library rule.
type protoTsLibraryFlags struct {
	// includeProtoLibraryData is true if the rule should populate the data
	// attribute with the proto_library rule.
	includeProtoLibraryData bool
}

func parseProtoTsLibraryFlags(kindName string, args []string) *protoTsLibraryFlags {
	flags := flag.NewFlagSet(kindName, flag.ExitOnError)

	var includeProtoLibraryData bool
	flags.BoolVar(&includeProtoLibraryData, "include_proto_library_data", false, "--include_proto_library_data=true populates the data attribute with the proto_library rule")

	if err := flags.Parse(args); err != nil {
		log.Fatalf("failed to parse flags for %q: %v", kindName, err)
	}

	return &protoTsLibraryFlags{includeProtoLibraryData}
}
