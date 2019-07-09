package main

var scalaWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

# rules_go used here to compile a wrapper around the protoc-gen-scala plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()`)

var scalaLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var scalaProtoLibraryRuleTemplate = mustTemplate(scalaLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    scala_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        exports = PROTO_DEPS,
        visibility = kwargs.get("visibility"),
    )

PROTO_DEPS = [
    "//external:io_bazel_rules_scala/dependency/com_google_protobuf/protobuf_java",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_fastparse",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_lenses",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_runtime",
]`)

var scalaGrpcLibraryRuleTemplate = mustTemplate(scalaLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    scala_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        exports = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
    )

GRPC_DEPS = [
    "//external:io_bazel_rules_scala/dependency/com_google_protobuf/protobuf_java",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_fastparse",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_lenses",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_runtime",
    "//external:io_bazel_rules_scala/dependency/proto/google_instrumentation",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_context",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_core",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_netty",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_protobuf",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_stub",
    "//external:io_bazel_rules_scala/dependency/proto/guava",
    "//external:io_bazel_rules_scala/dependency/proto/netty_buffer",
    "//external:io_bazel_rules_scala/dependency/proto/netty_codec",
    "//external:io_bazel_rules_scala/dependency/proto/netty_codec_http",
    "//external:io_bazel_rules_scala/dependency/proto/netty_codec_http2",
    "//external:io_bazel_rules_scala/dependency/proto/netty_codec_socks",
    "//external:io_bazel_rules_scala/dependency/proto/netty_common",
    "//external:io_bazel_rules_scala/dependency/proto/netty_handler",
    "//external:io_bazel_rules_scala/dependency/proto/netty_handler_proxy",
    "//external:io_bazel_rules_scala/dependency/proto/netty_resolver",
    "//external:io_bazel_rules_scala/dependency/proto/netty_transport",
    "//external:io_bazel_rules_scala/dependency/proto/opencensus_api",
    "//external:io_bazel_rules_scala/dependency/proto/opencensus_contrib_grpc_metrics",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_runtime_grpc",
]`)

func makeScala() *Language {
	return &Language{
		Dir:   "scala",
		Name:  "scala",
		Flags: commonLangFlags,
		SkipDirectoriesMerge: true,
		SkipTestPlatforms: []string{"windows"},
		Rules: []*Rule{
			&Rule{
				Name:             "scala_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//scala:scala"},
				WorkspaceExample: scalaWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates *.scala protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "scala_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//scala:grpc_scala"},
				WorkspaceExample: scalaWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates *.scala protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
				Experimental:     true,
			},
			&Rule{
				Name:             "scala_proto_library",
				Kind:             "proto",
				Implementation:   scalaProtoLibraryRuleTemplate,
				WorkspaceExample: scalaWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates *.scala protobuf library",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "scala_grpc_library",
				Kind:             "grpc",
				Implementation:   scalaGrpcLibraryRuleTemplate,
				WorkspaceExample: scalaWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates *.scala protobuf+gRPC library",
				Attrs:            aspectProtoCompileAttrs,
				Experimental:     true,
			},
		},
	}
}
