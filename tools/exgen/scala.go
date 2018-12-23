package main

var scalaUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

# rules_go used here to compile a wrapper around the protoc-gen-scala plugin
load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()`)

var scalaProtoLibraryRuleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:scala_proto_compile.bzl", "scala_proto_compile")
load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    scala_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    scala_library(
        name = name,
        srcs = [name_pb],
        deps = [str(Label("//scala:proto_deps"))],
        exports = [
            str(Label("//scala:proto_deps")),
        ],
        visibility = visibility,
    )`)

var scalaGrpcLibraryRuleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:scala_grpc_compile.bzl", "scala_grpc_compile")
load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    scala_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    scala_library(
        name = name,
        srcs = [name_pb],
        deps = [str(Label("//scala:grpc_deps"))],
        exports = [
            str(Label("//scala:grpc_deps")),
        ],
        visibility = visibility,
    )`)

func makeScala() *Language {
	return &Language{
		Dir:   "scala",
		Name:  "scala",
		Notes: mustTemplate("Rules for `scala_grpc_{compile|library}` don't produce code that compiles!  Use `@//io_bazel_rules_scala//scala_proto:scala_proto.bzl` instead"),
		Rules: []*Rule{
			&Rule{
				Name:           "scala_proto_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//scala:scala"},
				Usage:          scalaUsageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates *.scala protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "scala_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//scala:grpc_scala"},
				Usage:          scalaUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates *.scala protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Experimental:   true,
			},
			&Rule{
				Name:           "scala_proto_library",
				Implementation: scalaProtoLibraryRuleTemplate,
				Usage:          scalaUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates *.py protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "scala_grpc_library",
				Implementation: scalaGrpcLibraryRuleTemplate,
				Usage:          scalaUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates *.py protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Experimental:   true,
			},
		},
	}
}
