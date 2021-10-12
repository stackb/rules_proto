package rules_java

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	grpcJavaLibraryRuleName   = "grpc_java_library"
	grpcJavaLibraryRuleSuffix = "_grpc_java_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:grpc_java_library", &grpcJavaLibrary{})
}

// grpcJavaLibrary implements LanguageRule
type grpcJavaLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *grpcJavaLibrary) Name() string {
	return grpcJavaLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *grpcJavaLibrary) KindInfo() rule.KindInfo {
	return javaLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *grpcJavaLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/java:grpc_java_library.bzl",
		Symbols: []string{grpcJavaLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *grpcJavaLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("grpc:grpc-java:protoc-gen-grpc-java")
	if len(outputs) == 0 {
		return nil
	}

	return &JavaLibrary{
		KindName:       grpcJavaLibraryRuleName,
		RuleNameSuffix: grpcJavaLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver: func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
			deps := append(r.AttrStrings("deps"), ":"+pc.Library.BaseName()+ProtoJavaLibraryRuleSuffix)

			if len(deps) > 0 {
				r.SetAttr("deps", deps)
				r.SetAttr("exports", deps)
			}
		},
	}
}
