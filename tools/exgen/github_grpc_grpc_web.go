package main

var grpcWebUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)`)

func makeGithubComGrpcGrpcWeb() *Language {
	return &Language{
		Dir:  "github.com/grpc/grpc-web",
		Name: "grpc_web",
		Rules: []*Rule{
			&Rule{
				Name:           "closure_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc/grpc-web:closure"},
				Usage:          grpcWebUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a closure *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "commonjs_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc/grpc-web:commonjs"},
				Usage:          grpcWebUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a commonjs *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "commonjs_dts_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc/grpc-web:commonjs_dts"},
				Usage:          grpcWebUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a commonjs_dts *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "ts_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc/grpc-web:ts"},
				Usage:          grpcWebUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates a commonjs *.ts protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
