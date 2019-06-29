package main

var grpcjsWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()`)

var grpcjsLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()`)

var grpcjsGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:grpcjs_grpc_compile.bzl", "grpcjs_grpc_compile")
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

    grpcjs_grpc_compile(
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
		BazelCIExcludePlatforms: []string{"windows"}, // grpc.js protoc plugin uses features not supported on Windows cpp compiler
		Rules: []*Rule{
			&Rule{
				Name:             "grpcjs_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/stackb/grpc.js:grpc.js"},
				WorkspaceExample: grpcjsWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates protobuf closure grpc *.js files",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "grpcjs_grpc_library",
				Kind:             "grpc",
				Implementation:   grpcjsGrpcLibraryRuleTemplate,
				WorkspaceExample: grpcjsLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates protobuf closure library *.js files",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
