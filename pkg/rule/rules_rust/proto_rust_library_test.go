package rules_rust

import (
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	gproto "github.com/bazelbuild/bazel-gazelle/language/proto"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/google/go-cmp/cmp"
	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

func makeTestProtoLibrary(files ...*protoc.File) protoc.ProtoLibrary {
	r := rule.NewRule("proto_library", "test_proto")
	return protoc.NewOtherProtoLibrary(nil, r, files...)
}

func makeFile(dir, basename, protoContent string) *protoc.File {
	f := protoc.NewFile(dir, basename)
	if err := f.ParseReader(strings.NewReader(protoContent)); err != nil {
		panic("bad proto: " + err.Error())
	}
	return f
}

func languageRuleConfigWithResolve(spec string) protoc.LanguageRuleConfig {
	cfg := protoc.NewLanguageRuleConfig(config.New(), "rust")
	rewrite, err := protoc.ParseRewrite(spec)
	if err != nil {
		panic("bad resolve: " + err.Error())
	}
	cfg.Resolves = append(cfg.Resolves, *rewrite)
	return *cfg
}

func TestProtoRustLibraryRule(t *testing.T) {
	for name, tc := range map[string]struct {
		cfg  protoc.LanguageRuleConfig
		pc   protoc.ProtocConfiguration
		want string
	}{
		"degenerate": {
			cfg: *protoc.NewLanguageRuleConfig(config.New(), "rust"),
			pc: protoc.ProtocConfiguration{
				Library: makeTestProtoLibrary(),
			},
		},
		"simple": {
			cfg: *protoc.NewLanguageRuleConfig(config.New(), "rust"),
			pc: protoc.ProtocConfiguration{
				Library: makeTestProtoLibrary(
					makeFile("google/api", "http.proto", `syntax = "proto3"; package google.api; message HttpRule {}`),
				),
				Plugins: []*protoc.PluginConfiguration{
					{
						Config:  &protoc.LanguagePluginConfig{},
						Outputs: []string{"google.api.rs"},
					},
				},
			},
			want: `
proto_rust_library(
    name = "google_api",
    srcs = ["google.api.rs"],
    deps = [
        "@crates//:pbjson",
        "@crates//:prost",
        "@crates//:serde",
    ],
)
`,
		},
		"multiple srcs": {
			cfg: *protoc.NewLanguageRuleConfig(config.New(), "rust"),
			pc: protoc.ProtocConfiguration{
				Library: makeTestProtoLibrary(
					makeFile("trumid/common/proto", "date_range.proto", `syntax = "proto3"; package trumid.common.proto; message DateRange {}`),
				),
				Plugins: []*protoc.PluginConfiguration{
					{
						Config:  &protoc.LanguagePluginConfig{},
						Outputs: []string{"trumid.common.proto.rs", "trumid.common.proto.serde.rs"},
					},
				},
			},
			want: `
proto_rust_library(
    name = "trumid_common_proto",
    srcs = [
        "trumid.common.proto.rs",
        "trumid.common.proto.serde.rs",
    ],
    deps = [
        "@crates//:pbjson",
        "@crates//:prost",
        "@crates//:serde",
    ],
)
`,
		},
		"with well-known types": {
			cfg: *protoc.NewLanguageRuleConfig(config.New(), "rust"),
			pc: protoc.ProtocConfiguration{
				Library: makeTestProtoLibrary(
					makeFile("example/wkt", "thing.proto",
						`syntax = "proto3"; package example.wkt; import "google/protobuf/duration.proto"; message Thing { google.protobuf.Duration d = 1; }`),
				),
				Plugins: []*protoc.PluginConfiguration{
					{
						Config:  &protoc.LanguagePluginConfig{},
						Outputs: []string{"example.wkt.rs"},
					},
				},
			},
			want: `
proto_rust_library(
    name = "example_wkt",
    srcs = ["example.wkt.rs"],
    deps = [
        "@crates//:pbjson",
        "@crates//:prost",
        "@crates//:prost-types",
        "@crates//:serde",
    ],
)
`,
		},
		"with services": {
			cfg: *protoc.NewLanguageRuleConfig(config.New(), "rust"),
			pc: protoc.ProtocConfiguration{
				Library: makeTestProtoLibrary(
					makeFile("example/grpc", "greeter.proto", `syntax = "proto3"; package example.grpc; message HelloRequest {} service Greeter { rpc SayHello (HelloRequest) returns (HelloRequest); }`),
				),
				Plugins: []*protoc.PluginConfiguration{
					{
						Config:  &protoc.LanguagePluginConfig{},
						Outputs: []string{"example.grpc.rs", "example.grpc.tonic.rs"},
					},
				},
			},
			want: `
proto_rust_library(
    name = "example_grpc",
    srcs = [
        "example.grpc.rs",
        "example.grpc.tonic.rs",
    ],
    deps = [
        "@crates//:pbjson",
        "@crates//:prost",
        "@crates//:serde",
        "@crates//:tonic",
    ],
)
`,
		},
		"with rewritten dependency": {
			cfg: languageRuleConfigWithResolve(
				`^trumid/common/proto/guid[.]proto$ //trumid/common/proto/rust-types:trumid_common_proto_rs`,
			),
			pc: protoc.ProtocConfiguration{
				Library: makeTestProtoLibrary(
					makeFile("example/api", "thing.proto",
						`syntax = "proto3"; package example.api; import "trumid/common/proto/guid.proto"; message Thing {}`),
				),
				Plugins: []*protoc.PluginConfiguration{
					{
						Config:  &protoc.LanguagePluginConfig{},
						Outputs: []string{"example.api.rs"},
					},
				},
			},
			want: `
proto_rust_library(
    name = "example_api",
    srcs = ["example.api.rs"],
    deps = [
        "//trumid/common/proto/rust-types:trumid_common_proto_rs",
        "@crates//:pbjson",
        "@crates//:prost",
        "@crates//:serde",
    ],
)
`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			lib := protoRustLibrary{
				protoLibrariesByRule: make(map[label.Label][]protoc.ProtoLibrary),
			}
			impl := lib.ProvideRule(&tc.cfg, &tc.pc)
			var got string
			if impl != nil {
				r := impl.Rule()
				got = formatRules(r)
			}
			if diff := cmp.Diff(strings.TrimSpace(tc.want), strings.TrimSpace(got)); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

// TestProtoRustLibraryRuleMerge verifies that when Rule() is called with an
// existing rule of the same kind/name (otherGen), the new srcs/deps/imports are
// merged into it instead of creating a duplicate rule.
func TestProtoRustLibraryRuleMerge(t *testing.T) {
	cfg := protoc.NewLanguageRuleConfig(config.New(), "rust")
	pc1 := &protoc.ProtocConfiguration{
		Library: makeTestProtoLibrary(
			makeFile("merge/pkg", "first.proto",
				`syntax = "proto3"; package merge.pkg; message First {}`),
		),
		Plugins: []*protoc.PluginConfiguration{
			{
				Config:  &protoc.LanguagePluginConfig{},
				Outputs: []string{"merge.pkg.rs"},
			},
		},
	}
	pc2 := &protoc.ProtocConfiguration{
		Library: makeTestProtoLibrary(
			makeFile("merge/pkg", "second.proto",
				`syntax = "proto3"; package merge.pkg; message Second {}`),
		),
		Plugins: []*protoc.PluginConfiguration{
			{
				Config:  &protoc.LanguagePluginConfig{},
				Outputs: []string{"merge.pkg.serde.rs"},
			},
		},
	}

	lib := protoRustLibrary{
		protoLibrariesByRule: make(map[label.Label][]protoc.ProtoLibrary),
	}

	// First library generates a fresh rule.
	first := lib.ProvideRule(cfg, pc1).Rule()
	if first == nil {
		t.Fatal("first ProvideRule returned nil")
	}

	// Second library should merge into the first.
	merged := lib.ProvideRule(cfg, pc2).Rule(first)
	if merged != first {
		t.Errorf("expected second Rule() to return the same *Rule as the first (merge), got a different pointer")
	}

	gotSrcs := merged.AttrStrings("srcs")
	wantSrcs := []string{"merge.pkg.rs", "merge.pkg.serde.rs"}
	if diff := cmp.Diff(wantSrcs, gotSrcs); diff != "" {
		t.Errorf("merged srcs mismatch (-want +got):\n%s", diff)
	}
}

func formatRules(rules ...*rule.Rule) string {
	file := rule.EmptyFile("", "")
	for _, r := range rules {
		r.Insert(file)
	}
	return string(file.Format())
}

// TestProtoRustLibraryRule_PerFileMode_OneRulePerFile is the core regression
// test for the per-file scheme: in `gazelle:proto file` mode each per-file
// proto_library must produce a uniquely-named `proto_rust_library` rule, so
// the merge branch in `Rule()` is a no-op and consumer dep resolution can
// route every cross-file/cross-package reference to the specific per-file
// label (rather than collapsing through a per-package facade and re-closing
// the proto-package-level cycle).
//
// Name convention asserted: `<sanitized_pkg>__<stem>` where `<stem>` is the
// basename (minus `.proto`) of the first content-bearing file in the
// proto_library — `package.proto` co-includes are skipped so they don't
// steal the suffix.
func TestProtoRustLibraryRule_PerFileMode_OneRulePerFile(t *testing.T) {
	c := &config.Config{Exts: map[string]interface{}{}}
	c.Exts["proto"] = &gproto.ProtoConfig{Mode: gproto.FileMode}
	pkgCfg := protoc.NewPackageConfig(c)
	ruleCfg := protoc.NewLanguageRuleConfig(c, "rust")

	libA := protoc.NewOtherProtoLibrary(nil, rule.NewRule("proto_library", "a_proto"),
		makeFile("p/k", "a.proto", `syntax = "proto3"; package p.k; message A {}`))
	libB := protoc.NewOtherProtoLibrary(nil, rule.NewRule("proto_library", "b_proto"),
		makeFile("p/k", "b.proto", `syntax = "proto3"; package p.k; message B {}`))

	pcA := &protoc.ProtocConfiguration{
		PackageConfig: pkgCfg,
		Rel:           "p/k",
		Library:       libA,
		Plugins: []*protoc.PluginConfiguration{
			{Config: &protoc.LanguagePluginConfig{}, Outputs: []string{"a.p.k.rs"}},
		},
	}
	pcB := &protoc.ProtocConfiguration{
		PackageConfig: pkgCfg,
		Rel:           "p/k",
		Library:       libB,
		Plugins: []*protoc.PluginConfiguration{
			{Config: &protoc.LanguagePluginConfig{}, Outputs: []string{"b.p.k.rs"}},
		},
	}

	impl := protoRustLibrary{
		protoLibrariesByRule: make(map[label.Label][]protoc.ProtoLibrary),
	}
	rA := impl.ProvideRule(ruleCfg, pcA).Rule()
	rB := impl.ProvideRule(ruleCfg, pcB).Rule(rA)

	if rA.Name() != "p_k__a" {
		t.Errorf("first rule name: got %q, want %q", rA.Name(), "p_k__a")
	}
	if rB.Name() != "p_k__b" {
		t.Errorf("second rule name: got %q, want %q", rB.Name(), "p_k__b")
	}
	if rA == rB {
		t.Errorf("per-file rules must NOT merge — got the same *Rule pointer back from Rule(rA)")
	}
}

// TestProtoRustLibraryRule_PerFileMode_SkipsPackageProto verifies that when
// a per-file proto_library carries a shared `package.proto` co-include, the
// rule's `__<stem>` suffix is derived from the content-bearing file (not
// `package`).
func TestProtoRustLibraryRule_PerFileMode_SkipsPackageProto(t *testing.T) {
	c := &config.Config{Exts: map[string]interface{}{}}
	c.Exts["proto"] = &gproto.ProtoConfig{Mode: gproto.FileMode}
	pkgCfg := protoc.NewPackageConfig(c)
	ruleCfg := protoc.NewLanguageRuleConfig(c, "rust")

	// Note `package.proto` comes FIRST in the file list — the suffix must
	// still resolve to `order`, the content-bearing file.
	lib := protoc.NewOtherProtoLibrary(nil, rule.NewRule("proto_library", "order_proto"),
		makeFile("p/k", "package.proto", `syntax = "proto3"; package p.k;`),
		makeFile("p/k", "order.proto", `syntax = "proto3"; package p.k; message Order {}`),
	)

	pc := &protoc.ProtocConfiguration{
		PackageConfig: pkgCfg,
		Rel:           "p/k",
		Library:       lib,
		Plugins: []*protoc.PluginConfiguration{
			{Config: &protoc.LanguagePluginConfig{}, Outputs: []string{"order.p.k.rs"}},
		},
	}

	impl := protoRustLibrary{
		protoLibrariesByRule: make(map[label.Label][]protoc.ProtoLibrary),
	}
	r := impl.ProvideRule(ruleCfg, pc).Rule()
	if r.Name() != "p_k__order" {
		t.Errorf("rule name: got %q, want %q", r.Name(), "p_k__order")
	}
}
