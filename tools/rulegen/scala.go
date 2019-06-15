package main

var scalaUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

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

var scalaLibraryRuleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create {{ .Lang.Name }} library
    scala_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [str(Label("//scala:{{ .Rule.Kind }}_deps"))],
        exports = [
            str(Label("//scala:{{ .Rule.Kind }}_deps")),
        ],
        visibility = kwargs.get("visibility"),
    )`)

func makeScala() *Language {
	return &Language{
		Dir:   "scala",
		Name:  "scala",
		Notes: mustTemplate("Rules for `scala_grpc_{compile|library}` don't produce code that compiles!  Use `@//io_bazel_rules_scala//scala_proto:scala_proto.bzl` instead"),
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:           "scala_proto_compile",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//scala:scala"},
				Usage:          scalaUsageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates *.scala protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "scala_grpc_compile",
				Kind:           "grpc",
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
				Kind:           "proto",
				Implementation: scalaLibraryRuleTemplate,
				Usage:          scalaUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates *.scala protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "scala_grpc_library",
				Kind:           "grpc",
				Implementation: scalaLibraryRuleTemplate,
				Usage:          scalaUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates *.scala protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Experimental:   true,
			},
		},
	}
}
