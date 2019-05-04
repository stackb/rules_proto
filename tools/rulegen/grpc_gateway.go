package main

var grpcGatewayUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

bazel_gazelle()

load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()`)

var grpcGatewayLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:gateway_grpc_compile.bzl", "gateway_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

def {{ .Rule.Name }}(**kwargs):
    name = kwargs.get("name")
    importpath = kwargs.get("importpath")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    compilers = kwargs.get("compilers")
    if not compilers:
        compilers = [
            "@io_bazel_rules_go//proto:go_grpc",
            "@grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
        ]

    go_proto_library(
        name = name,
        compilers = compilers,
        importpath = importpath,
        proto = deps[0],
        deps = ["@go_googleapis//google/api:annotations_go_proto"],
        visibility = visibility,
    )`)

var grpcGatewayCompileExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "api_gateway_grpc",
    deps = ["@build_stack_rules_proto//{{ .Lang.Dir }}/example/api:api_proto"],
)`)

var grpcGatewayLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "api_gateway_library",
    importpath = "github.com/stackb/rules_proto/github.com/grpc-ecosystem/grpc-gateway/examples/api",
    deps = ["@build_stack_rules_proto//{{ .Lang.Dir }}/example/api:api_proto"],
)`)

func makeGrpcGateway() *Language {
	return &Language{
		Dir:  "github.com/grpc-ecosystem/grpc-gateway",
		Name: "grpc-gateway",
		Rules: []*Rule{
			&Rule{
				Name:           "gateway_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc-ecosystem/grpc-gateway:grpc-gateway"},
				Usage:          grpcGatewayUsageTemplate,
				Example:        grpcGatewayCompileExampleTemplate,
				Doc:            "Generates grpc-gateway *.go files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "gateway_swagger_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//github.com/grpc-ecosystem/grpc-gateway:swagger"},
				Usage:          grpcGatewayUsageTemplate,
				Example:        grpcGatewayCompileExampleTemplate,
				Doc:            "Generates grpc-gateway swagger *.json files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "gateway_grpc_library",
				Implementation: grpcGatewayLibraryRuleTemplate,
				Usage:          grpcGatewayUsageTemplate,
				Example:        grpcGatewayLibraryExampleTemplate,
				Doc:            "Generates grpc-gateway library files",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
