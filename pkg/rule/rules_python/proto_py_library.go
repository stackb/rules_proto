package rules_python

import (
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoPyLibraryRuleName   = "proto_py_library"
	ProtoPyLibraryRuleSuffix = "_py_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_py_library", &protoPyLibrary{})
}

// protoPyLibrary implements LanguageRule for the 'proto_py_library' rule from
// @rules_proto.
type protoPyLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoPyLibrary) Name() string {
	return ProtoPyLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoPyLibrary) KindInfo() rule.KindInfo {
	return pyLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoPyLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/py:proto_py_library.bzl",
		Symbols: []string{ProtoPyLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoPyLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("builtin:python")
	if len(outputs) == 0 {
		return nil
	}
	return &PyLibrary{
		KindName:       ProtoPyLibraryRuleName,
		RuleNameSuffix: ProtoPyLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver:       protoc.ResolveDepsAttr("deps", true),
	}
}
