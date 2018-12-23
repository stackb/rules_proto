package main

var grpcWebUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)`)

var grpcWebGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:closure_grpc_compile.bzl", "closure_grpc_compile")
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
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:abstractclientbase",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:clientreadablestream",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:grpcwebclientbase",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:error",
            "@io_bazel_rules_closure//closure/library",
            "@io_bazel_rules_closure//closure/protobuf:jspb",
        ],
        lenient = True,
        visibility = visibility,
    )`)

func makeGithubComGrpcGrpcWeb() *Language {
	return &Language{
		Dir:  "github.com/grpc/grpc-web",
		Name: "grpc-web",
		Rules: []*Rule{
			&Rule{
				Name:           "closure_grpc_compile",
				Usage:          grpcWebUsageTemplate,
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc/grpc-web:closure"},
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a closure *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "commonjs_grpc_compile",
				Usage:          grpcWebUsageTemplate,
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc/grpc-web:commonjs"},
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a commonjs *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "commonjs_dts_grpc_compile",
				Usage:          grpcWebUsageTemplate,
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc/grpc-web:commonjs_dts"},
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a commonjs_dts *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "ts_grpc_compile",
				Usage:          grpcWebUsageTemplate,
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc/grpc-web:ts"},
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a commonjs *.ts protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "closure_grpc_library",
				Usage:          grpcWebUsageTemplate,
				Implementation: grpcWebGrpcLibraryRuleTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates protobuf closure library *.js files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
