package main

var grpcjsUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")
{{ .Rule.Name }}()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")
closure_repositories(omit_com_google_protobuf = True)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()`)

var grpcjsGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:closure_grpc_compile.bzl", "closure_grpc_compile")
load("//closure:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_pb_grpc = name + "_pb_grpc"

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )
    
    closure_grpc_compile(
        name = name_pb_grpc,
        deps = deps,
        transitive = True,
        visibility = visibility,
        verbose = verbose,
    )

    closure_js_library(
        name = name,
        srcs = [name_pb, name_pb_grpc],
        deps = [
            "@io_bazel_rules_closure//closure/library",
            "@io_bazel_rules_closure//closure/protobuf:jspb",
            "@com_github_stackb_grpc_js//js/grpc/stream:observer",
            "@com_github_stackb_grpc_js//js/grpc/stream/observer:call",
            "@com_github_stackb_grpc_js//js/grpc",
            "@com_github_stackb_grpc_js//js/grpc:api",
            "@com_github_stackb_grpc_js//js/grpc:options",
        ],
        internal_descriptors = [
            name_pb + "/descriptor.source.bin",
            name_pb_grpc + "/descriptor.source.bin",
        ],
        lenient = True,
        suppress = [
            "JSC_WRONG_ARGUMENT_COUNT",
        ],
        visibility = visibility,
    )
`)

func makeGrpcJs() *Language {
	return &Language{
		Dir:  "github.com/stackb/grpc.js",
		Name: "grpc.js",
		Rules: []*Rule{
			&Rule{
				Name:           "closure_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/stackb/grpc.js:grpc.js"},
				Usage:          grpcjsUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates protobuf closure grpc *.js files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "closure_grpc_library",
				Implementation: grpcjsGrpcLibraryRuleTemplate,
				Usage:          grpcjsUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates protobuf closure library *.js files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
