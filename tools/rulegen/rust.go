package main

var rustUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto/raze:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()`)

var rustProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir}}:rust_proto_compile.bzl", "rust_proto_compile")
load("//{{ .Lang.Dir }}:rust_proto_lib.bzl", "rust_proto_lib")
load("@io_bazel_rules_rust//rust:rust.bzl", "rust_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_lib = name + "_lib"

    rust_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    rust_proto_lib(
        name = name_lib,
        compilation = name_pb,
    )

    rust_library(
        name = name,
        srcs = [name_pb, name_lib],
        deps = [
            "@io_bazel_rules_rust//proto/raze:protobuf",
        ],
        visibility = visibility,
    )`)

var rustGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir}}:rust_grpc_compile.bzl", "rust_grpc_compile")
load("//{{ .Lang.Dir }}:rust_proto_lib.bzl", "rust_proto_lib")
load("@io_bazel_rules_rust//rust:rust.bzl", "rust_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_lib = name + "_lib"

    rust_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    rust_proto_lib(
        name = name_lib,
        compilation = name_pb,
    )

    rust_library(
        name = name,
        srcs = [name_pb, name_lib],
        deps = [
            "@io_bazel_rules_rust//proto/raze:protobuf",
            "@io_bazel_rules_rust//proto/raze:grpc",
            "@io_bazel_rules_rust//proto/raze:tls_api",
            "@io_bazel_rules_rust//proto/raze:tls_api_stub",
        ],
        visibility = visibility,
    )`)

func makeRust() *Language {
	return &Language{
		Dir:  "rust",
		Name: "rust",
		Flags: append(commonLangFlags, &Flag{
			Category:    "build",
			Name:        "incompatible_enable_cc_toolchain_resolution",
			Value:       "false",
			Description: "In order to use find_cpp_toolchain, you must include the '@bazel_tools//tools/cpp:toolchain_type' in the toolchains argument to your rule.",
		}, &Flag{
			Category:    "build",
			Name:        "incompatible_require_ctx_in_configure_features",
			Value:       "false",
			Description: `/external/io_bazel_rules_rust/rust/private/rustc.bzl", line 143, in _get_linker_and_args cc_common.configure_features(cc_toolchain = cc_toolchain, reque..., ...) Incompatible flag --incompatible_require_ctx_in_configure_features has been flipped, and the mandatory parameter 'ctx' of cc_common.configure_features is missing`,
		}),
		Rules: []*Rule{
			&Rule{
				Name:           "rust_proto_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//rust:rust"},
				Usage:          rustUsageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates rust protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "rust_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//rust:rust", "//rust:grpc_rust"},
				Usage:          rustUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates rust protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "rust_proto_library",
				Implementation: rustProtoLibraryRuleTemplate,
				Usage:          rustUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates rust protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "rust_grpc_library",
				Implementation: rustGrpcLibraryRuleTemplate,
				Usage:          rustUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates rust protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
