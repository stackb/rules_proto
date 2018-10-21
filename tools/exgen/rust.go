package main

var rustUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@build_stack_rules_proto//rust/cargo:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()
`)

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
        transitive = True,
        visibility = visibility,
    )

    rust_proto_lib(
        name = name_lib,
        compilation = name_pb,
    )

    rust_library(
        name = name,
        srcs = [name_pb, name_lib],
        deps = [
            str(Label("//rust/cargo:protobuf")),
        ],
        visibility = visibility,
    )
`)

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
        transitive = True,
        visibility = visibility,
    )

    rust_proto_lib(
        name = name_lib,
        compilation = name_pb,
    )

    rust_library(
        name = name,
        srcs = [name_pb, name_lib],
        deps = [
            str(Label("//rust/cargo:protobuf")),
            str(Label("//rust/cargo:grpc")),
            str(Label("//rust/cargo:tls_api")),
            str(Label("//rust/cargo:tls_api_stub")),
        ],
        visibility = visibility,
    )
`)

func makeRust() *Language {
	return &Language{
		Dir:  "rust",
		Name: "rust",
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
