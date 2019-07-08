package main

var grpcWebWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "grpc_web_deps")

grpc_web_deps()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)`)

var grpcWebGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:closure_grpc_compile.bzl", "closure_grpc_compile")
load("//closure:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    name_pb_grpc = kwargs.get("name") + "_pb_grpc"
    closure_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    closure_grpc_compile(
        name = name_pb_grpc,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create closure library
    closure_js_library(
        name = kwargs.get("name"),
        srcs = [name_pb, name_pb_grpc],
        deps = GRPC_DEPS,
        suppress = [
            "JSC_LATE_PROVIDE_ERROR",
            "JSC_UNDEFINED_VARIABLE",
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            "JSC_STRICT_INEXISTENT_PROPERTY",
            "JSC_POSSIBLE_INEXISTENT_PROPERTY",
            "JSC_UNRECOGNIZED_TYPE_ERROR",
            "JSC_UNUSED_PRIVATE_PROPERTY",
            "JSC_EXTRA_REQUIRE_WARNING",
            "JSC_INVALID_INTERFACE_MEMBER_DECLARATION",
        ],
        visibility = kwargs.get("visibility"),
    )

GRPC_DEPS = [
    "@com_github_grpc_grpc_web//javascript/net/grpc/web:abstractclientbase",
    "@com_github_grpc_grpc_web//javascript/net/grpc/web:clientreadablestream",
    "@com_github_grpc_grpc_web//javascript/net/grpc/web:grpcwebclientbase",
    "@com_github_grpc_grpc_web//javascript/net/grpc/web:error",
    "@io_bazel_rules_closure//closure/library",
    "@io_bazel_rules_closure//closure/protobuf:jspb",
]`)

func makeGithubComGrpcGrpcWeb() *Language {
	return &Language{
		Dir:  "github.com/grpc/grpc-web",
		Name: "grpc-web",
		Rules: []*Rule{
			&Rule{
				Name:             "closure_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/grpc/grpc-web:closure"},
				WorkspaceExample: grpcWebWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates a closure *.js protobuf+gRPC files",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "commonjs_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/grpc/grpc-web:commonjs"},
				WorkspaceExample: grpcWebWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates a commonjs *.js protobuf+gRPC files",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "commonjs_dts_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/grpc/grpc-web:commonjs_dts"},
				WorkspaceExample: grpcWebWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates a commonjs_dts *.js protobuf+gRPC files",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "ts_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/grpc/grpc-web:ts"},
				WorkspaceExample: grpcWebWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates a commonjs *.ts protobuf+gRPC files",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "closure_grpc_library",
				Kind:             "grpc",
				Implementation:   grpcWebGrpcLibraryRuleTemplate,
				WorkspaceExample: grpcWebWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates protobuf closure library *.js files",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
