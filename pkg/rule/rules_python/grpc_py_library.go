package rules_python

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	grpcPyLibraryRuleName   = "grpc_py_library"
	grpcPyLibraryRuleSuffix = "_grpc_py_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:grpc_py_library", &grpcPyLibrary{})
}

// grpcPyLibrary implements LanguageRule for the 'grpc_py_library' rule from
// @rules_proto.
type grpcPyLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *grpcPyLibrary) Name() string {
	return grpcPyLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *grpcPyLibrary) KindInfo() rule.KindInfo {
	return pyLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *grpcPyLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/py:grpc_py_library.bzl",
		Symbols: []string{grpcPyLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *grpcPyLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("grpc:grpc:protoc-gen-grpc-python")
	if len(outputs) == 0 {
		return nil
	}

	return &PyLibrary{
		KindName:       grpcPyLibraryRuleName,
		RuleNameSuffix: grpcPyLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver: func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
			deps := append(r.AttrStrings("deps"), ":"+pc.Library.BaseName()+ProtoPyLibraryRuleSuffix)

			if len(deps) > 0 {
				r.SetAttr("deps", deps)
			}
		},
	}
}
