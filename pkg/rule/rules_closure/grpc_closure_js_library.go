package rules_closure

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	GrpcClosureJsLibraryRuleName   = "grpc_closure_js_library"
	GrpcClosureJsLibraryRuleSuffix = "_grpc_closure_js_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:grpc_closure_js_library", &grpcClosureJsLibrary{})
}

// grpcClosureJsLibrary implements LanguageRule for the
// 'grpc_closure_js_library' rule.
type grpcClosureJsLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *grpcClosureJsLibrary) Name() string {
	return GrpcClosureJsLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *grpcClosureJsLibrary) KindInfo() rule.KindInfo {
	return closureJsLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *grpcClosureJsLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/closure:grpc_closure_js_library.bzl",
		Symbols: []string{GrpcClosureJsLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *grpcClosureJsLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("stackb:grpc.js:protoc-gen-grpc-js")
	if len(outputs) == 0 {
		return nil
	}
	return &ClosureJsLibrary{
		KindName:       GrpcClosureJsLibraryRuleName,
		RuleNameSuffix: GrpcClosureJsLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver: func(impl protoc.DepsProvider, pc *protoc.ProtocConfiguration, c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
			protoDep := ":" + pc.Library.BaseName() + ProtoClosureJsLibraryRuleSuffix

			deps := impl.Deps()
			deps = append(deps, protoDep)

			if len(deps) > 0 {
				r.SetAttr("deps", deps)
				r.SetAttr("exports", deps)
			}
		},
	}
}
