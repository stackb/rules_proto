package main

var grpcWebUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Name }}:deps.bzl", "{{ .Rule.Name }}")

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
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates a closure *.js protobuf+gRPC files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
