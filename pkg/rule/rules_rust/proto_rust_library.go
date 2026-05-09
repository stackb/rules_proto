package rules_rust

import (
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

const (
	ProtoRustLibraryRuleName   = "proto_rust_library"
	ProtoRustLibraryRuleSuffix = "_rust_library"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:proto_rust_library", &protoRustLibrary{
		protoLibrariesByRule: make(map[label.Label][]protoc.ProtoLibrary),
	})
}

// protoRustLibrary implements LanguageRule for the 'proto_rust_library' rule.
type protoRustLibrary struct {
	protoLibrariesByRule map[label.Label][]protoc.ProtoLibrary
}

// Name implements part of the LanguageRule interface.
func (s *protoRustLibrary) Name() string {
	return ProtoRustLibraryRuleName
}

// KindInfo implements part of the LanguageRule interface.
func (s *protoRustLibrary) KindInfo() rule.KindInfo {
	return rustLibraryKindInfo
}

// LoadInfo implements part of the LanguageRule interface.
func (s *protoRustLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    "@build_stack_rules_proto//rules/rust:proto_rust_library.bzl",
		Symbols: []string{ProtoRustLibraryRuleName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *protoRustLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	outputs := make([]string, 0)
	for _, plugin := range pc.Plugins {
		for _, out := range plugin.Outputs {
			if strings.HasSuffix(out, ".rs") {
				outputs = append(outputs, out)
			}
		}
	}
	if len(outputs) == 0 {
		return nil
	}

	// Compute Rust keyword escape mappings for proto packages containing
	// Rust reserved keywords (e.g., "google.type" → prost writes to
	// "google/r#type/" instead of "google/type/").
	if files := pc.Library.Files(); len(files) > 0 {
		pkg := files[0].Package().Name
		for output, escapedPath := range protoc.RustKeywordEscapeMappings(pkg, outputs) {
			if pc.Mappings == nil {
				pc.Mappings = make(map[string]string)
			}
			pc.Mappings[output] = escapedPath
		}
	}

	rl := &RustLibrary{
		KindName:             ProtoRustLibraryRuleName,
		RuleNameSuffix:       ProtoRustLibraryRuleSuffix,
		Outputs:              outputs,
		RuleConfig:           cfg,
		Config:               pc,
		Resolver:             protoc.ResolveDepsAttr("deps", false),
		protoLibrariesByRule: s.protoLibrariesByRule,
	}
	rl.id = label.New("", pc.Rel, rl.Name())
	return rl
}
