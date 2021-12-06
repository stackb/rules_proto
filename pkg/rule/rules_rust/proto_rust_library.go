package rules_rust

import (
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	protoRustLibraryRuleName   = "proto_rust_library"
	protoRustLibraryRuleSuffix = "_rust_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_rust_library", &protoRustLibrary{})
}

// protoRustLibrary implements LanguageRule for the 'proto_rust_library' rule.
type protoRustLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *protoRustLibrary) Name() string {
	return protoRustLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoRustLibrary) KindInfo() rule.KindInfo {
	return rustLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoRustLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/rust:proto_rust_library.bzl",
		Symbols: []string{protoRustLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoRustLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := pc.GetPluginOutputs("stepancheg:rust-protobuf:protoc-gen-rust")
	if len(outputs) == 0 {
		return nil
	}
	return &rustLibrary{
		KindName:       protoRustLibraryRuleName,
		RuleNameSuffix: protoRustLibraryRuleSuffix,
		Outputs:        outputs,
		RuleConfig:     cfg,
		Config:         pc,
		Resolver:       protoc.ResolveDepsAttr("deps", true),
	}
}
