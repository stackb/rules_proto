package main

var dartUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

dart_repositories()

load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

dart_protoc_plugin_deps()`)

var dartGrpcLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

dart_repositories()

load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

dart_protoc_plugin_deps()

load("@dart_pub_deps_grpc//:deps.bzl", dart_grpc_deps = "pub_deps")

dart_grpc_deps()`)

var dartProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir}}:dart_proto_compile.bzl", "dart_proto_compile")
load("@io_bazel_rules_dart//dart/build_rules:core.bzl", "dart_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    dart_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    dart_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("@vendor_protobuf//:protobuf")),
        ],
        pub_pkg_name = name,
        visibility = visibility,
    )`)

var dartGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir}}:dart_grpc_compile.bzl", "dart_grpc_compile")
load("@io_bazel_rules_dart//dart/build_rules:core.bzl", "dart_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    dart_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    dart_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("@vendor_protobuf//:protobuf")),
            str(Label("@vendor_grpc//:grpc")),
        ],
        pub_pkg_name = name,
        visibility = visibility,
    )`)

var dartFlags = []*Flag{
	// {
	// 	Category:    "build",
	// 	Name:        "incompatible_disallow_data_transition",
	// 	Value:       "false",
	// 	Description: "vm.bzl is still using cfg=data",
	// },
	{
		Category: "build",
		Name:     "incompatible_no_transitive_loads",
		Value:    "false",
	},
	{
		Category: "build",
		Name:     "incompatible_disable_deprecated_attr_params",
		Value:    "false",
	},
	{
		Category: "build",
		Name:     "incompatible_enable_cc_toolchain_resolution",
		Value:    "false",
	},
	{
		Category: "build",
		Name:     "incompatible_require_ctx_in_configure_features",
		Value:    "false",
	},
	{
		Category: "build",
		Name:     "incompatible_depset_is_not_iterable",
		Value:    "false",
	},
	{
		Category: "build",
		Name:     "incompatible_depset_union",
		Value:    "false",
	},
	{
		Category: "build",
		Name:     "incompatible_disallow_struct_provider_syntax",
		Value:    "false",
	},
}

func makeDart() *Language {
	return &Language{
		Dir:   "dart",
		Name:  "dart",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:           "dart_proto_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//dart:dart"},
				Usage:          dartUsageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates dart protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags:          dartFlags,
			},
			&Rule{
				Name:           "dart_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//dart:grpc_dart"},
				Usage:          dartUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates dart protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags:          dartFlags,
			},
			&Rule{
				Name:           "dart_proto_library",
				Implementation: dartProtoLibraryRuleTemplate,
				Usage:          dartUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Doc:            "Generates dart protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags:          dartFlags,
			},
			&Rule{
				Name:           "dart_grpc_library",
				Implementation: dartGrpcLibraryRuleTemplate,
				Usage:          dartGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates dart protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags:          dartFlags,
			},
		},
	}
}
