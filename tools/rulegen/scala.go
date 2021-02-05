package main

var scalaProtoWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()`)

var scalaGrpcWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories()`)

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
        tags = kwargs.get("tags"),
    )

PROTO_DEPS = [
    "@io_bazel_rules_scala//scala_proto:default_scalapb_compile_dependencies",
]`)

var scalaGrpcLibraryRuleTemplate = mustTemplate(scalaLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    scala_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        exports = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    "@io_bazel_rules_scala//scala_proto:default_scalapb_compile_dependencies",
    "@io_bazel_rules_scala//scala_proto:default_scalapb_grpc_dependencies",
]`)

func makeScala() *Language {
	return &Language{
		Dir:   "scala",
		Name:  "scala",
		DisplayName: "Scala",
		Notes: mustTemplate("Rules for generating Scala protobuf and gRPC `.jar` files and libraries using [ScalaPB](https://github.com/scalapb/ScalaPB). Libraries are created with `scala_library` from [rules_scala](https://github.com/bazelbuild/rules_scala)"),
		Flags: commonLangFlags,
		SkipDirectoriesMerge: true,
		SkipTestPlatforms: []string{},
		Rules: []*Rule{
			&Rule{
				Name:             "scala_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//scala:scala_plugin"},
				WorkspaceExample: scalaProtoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates a Scala protobuf `.jar` artifact",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"},
			},
			&Rule{
				Name:             "scala_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//scala:grpc_scala_plugin"},
				WorkspaceExample: scalaGrpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Scala protobuf+gRPC `.jar` artifacts",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"},
				Experimental:     true,
			},
			&Rule{
				Name:             "scala_proto_library",
				Kind:             "proto",
				Implementation:   scalaProtoLibraryRuleTemplate,
				WorkspaceExample: scalaProtoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates a Scala protobuf library using `scala_library` from `rules_scala`",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"},
			},
			&Rule{
				Name:             "scala_grpc_library",
				Kind:             "grpc",
				Implementation:   scalaGrpcLibraryRuleTemplate,
				WorkspaceExample: scalaGrpcWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates a Scala protobuf+gRPC library using `scala_library` from `rules_scala`",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"},
				Experimental:     true,
			},
		},
	}
}
