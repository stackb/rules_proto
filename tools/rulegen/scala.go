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
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    scala_library(
        name = name,
        srcs = [name_pb],
        deps = [str(Label("//scala:{{ .Rule.Kind }}_deps"))],
        exports = [
            str(Label("//scala:{{ .Rule.Kind }}_deps")),
        ],
        visibility = visibility,
    )`)

func makeScala() *Language {
	return &Language{
		Dir:   "scala",
		Name:  "scala",
		Notes: mustTemplate("Rules for `scala_grpc_{compile|library}` don't produce code that compiles!  Use `@//io_bazel_rules_scala//scala_proto:scala_proto.bzl` instead"),
		Flags: append(commonLangFlags, &Flag{
			Category:    "build",
			Name:        "incompatible_enable_cc_toolchain_resolution",
			Value:       "false",
			Description: "In order to use find_cpp_toolchain, you must include the '@bazel_tools//tools/cpp:toolchain_type' in the toolchains argument to your rule.",
		}, &Flag{
			Category: "build",
			Name:     "incompatible_require_ctx_in_configure_features",
			Value:    "false",
		}, &Flag{
			Category:    "build",
			Name:        "incompatible_disallow_legacy_javainfo",
			Value:       "false",
			Description: `/external/io_bazel_rules_scala/scala/scala_import.bzl", line 73, in _create_provider java_common.create_provider(use_ijar = False, compile_time_jar...]), <3 more arguments>) java_common.create_provider is deprecated and cannot be used when --incompatible_disallow_legacy_javainfo is set. Please migrate to the JavaInfo constructor.`,
		}, &Flag{
			Category:    "build",
			Name:        "incompatible_disallow_struct_provider_syntax",
			Value:       "false",
			Description: `external/io_bazel_rules_scala/scala/scala_import.bzl`,
		}, &Flag{
			Category: "build",
			Name:     "incompatible_use_toolchain_resolution_for_java_rules",
			Value:    "false",
		}),
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
