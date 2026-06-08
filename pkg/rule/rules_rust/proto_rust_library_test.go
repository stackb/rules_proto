package rules_rust

import (
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
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
