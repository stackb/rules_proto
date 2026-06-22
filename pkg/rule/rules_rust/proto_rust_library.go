package rules_rust

import (
	"path"
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

	// Compute output_mappings whenever the directory protoc-gen-prost writes
	// to (derived from the proto package name, with Rust keyword segments
	// r#-escaped) differs from the bazel package the rule lives in. Two
	// distinct causes:
	//
	//   1. Rust keyword escapes — proto package "google.type" → prost writes
	//      to "google/r#type/" while the bazel pkg is "google/type".
	//
	//   2. Proto package path simply differs from bazel package path —
	//      e.g. proto package "trumid.common.auth" living at bazel pkg
	//      "trumid/common/auth/proto", or "grpc.health.v1" living at
	//      "thirdparty/protobuf/grpc/src/main/protobuf/grpc/health/v1".
	//      Without a mapping, proto_compile.bzl's rename step looks for the
	//      output at <bazel-bin>/<bazel_pkg>/<file>.rs and fails with
	//      `mv: ... No such file or directory`.
	if files := pc.Library.Files(); len(files) > 0 {
		pkg := files[0].Package().Name
		protocDir := protoc.RustProtocOutputDir(pkg)
		if protocDir != "" && protocDir != pc.Rel {
			if pc.Mappings == nil {
				pc.Mappings = make(map[string]string)
			}
			for _, output := range outputs {
				base := path.Base(output)
				pc.Mappings[base] = path.Join(protocDir, base)
			}
		}
	}

	rl := &RustLibrary{
		KindName:             ProtoRustLibraryRuleName,
		RuleNameSuffix:       ProtoRustLibraryRuleSuffix,
		Outputs:              outputs,
		RuleConfig:           cfg,
		Config:               pc,
		Resolver:             protoc.ResolveDepsAttr("deps", false),
		PerFile:              pc.IsProtoFileMode(),
		protoLibrariesByRule: s.protoLibrariesByRule,
	}
	rl.id = label.New("", pc.Rel, rl.Name())
	return rl
}
