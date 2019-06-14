package main

var grpcjsLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()`)

var grpcjsUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()`)

var grpcjsGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:closure_grpc_compile.bzl", "closure_grpc_compile")
load("//closure:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    name_pb_lib = kwargs.get("name") + "_pb_lib"
    name_pb_grpc = kwargs.get("name") + "_pb_grpc"

    closure_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k not in ("name", "closure_deps")} # Forward args except name and closure_deps
    )

    closure_grpc_compile(
        name = name_pb_grpc,
        **{k: v for (k, v) in kwargs.items() if k not in ("name", "closure_deps")} # Forward args except name and closure_deps
    )

    # Create libraries
    closure_js_library(
        name = name_pb_lib,
        srcs = [name_pb],
        deps = [
            "@io_bazel_rules_closure//closure/protobuf:jspb",
        ] + kwargs.get("closure_deps", []),
        internal_descriptors = [
            name_pb + "/descriptor.source.bin",
        ],
        suppress = [
            "JSC_LATE_PROVIDE_ERROR",
            "JSC_UNDEFINED_VARIABLE",
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            "JSC_STRICT_INEXISTENT_PROPERTY",
            "JSC_POSSIBLE_INEXISTENT_PROPERTY",
            "JSC_UNRECOGNIZED_TYPE_ERROR",
        ],
        visibility = kwargs.get("visibility"),
    )

    closure_js_library(
        name = kwargs.get("name"),
        srcs = [name_pb_grpc],
        deps = [
            name_pb_lib,
            "@io_bazel_rules_closure//closure/library/promise",
            "@com_github_stackb_grpc_js//js/grpc/stream:observer",
            "@com_github_stackb_grpc_js//js/grpc/stream/observer:call",
            "@com_github_stackb_grpc_js//js/grpc",
            "@com_github_stackb_grpc_js//js/grpc:api",
            "@com_github_stackb_grpc_js//js/grpc:options",
        ] + kwargs.get("closure_deps", []),
        internal_descriptors = [
            name_pb + "/descriptor.source.bin",
            name_pb_grpc + "/descriptor.source.bin",
        ],
        exports = [
            name_pb_lib,
        ],
        suppress = [
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
        ],
        visibility = kwargs.get("visibility"),
    )`)

func makeGrpcJs() *Language {
	return &Language{
		Dir:  "github.com/stackb/grpc.js",
		Name: "grpc.js",
		Rules: []*Rule{
			&Rule{
				Name:           "closure_grpc_compile",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/stackb/grpc.js:grpc.js"},
				Usage:          grpcjsUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates protobuf closure grpc *.js files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "closure_grpc_library",
				Kind:           "grpc",
				Implementation: grpcjsGrpcLibraryRuleTemplate,
				Usage:          grpcjsLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates protobuf closure library *.js files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
				Flags: []*Flag{
					{
						Category: "build",
						Name:     "incompatible_use_toolchain_resolution_for_java_rules",
						Value:    "false",
					},
					{
						Category: "build",
						Name:     "incompatible_disallow_struct_provider_syntax",
						Value:    "false",
					},
				},
			},
		},
	}
}
