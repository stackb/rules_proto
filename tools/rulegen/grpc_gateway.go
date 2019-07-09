package main

var grpcGatewayWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "gateway_deps")

gateway_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()`)

var grpcGatewayLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:gateway_grpc_compile.bzl", "gateway_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

def {{ .Rule.Name }}(**kwargs):
    # Apply default args
    if not kwargs.get("compilers"):
        kwargs["compilers"] = [
            "@io_bazel_rules_go//proto:go_grpc",
            "@grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
        ]

    # Create go library
    go_proto_library(
        proto = kwargs.get("deps")[0],
        deps = ["@go_googleapis//google/api:annotations_go_proto"],
        **{k: v for (k, v) in kwargs.items() if k != "deps"} # Forward args except deps
    )`)

var grpcGatewayCompileExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "api_gateway_grpc",
    deps = ["@build_stack_rules_proto//{{ .Lang.Dir }}/example/api:api_proto"],
)`)

var grpcGatewayLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "api_gateway_library",
    importpath = "github.com/stackb/rules_proto/github.com/grpc-ecosystem/grpc-gateway/examples/api",
    deps = ["@build_stack_rules_proto//{{ .Lang.Dir }}/example/api:api_proto"],
)`)

func makeGrpcGateway() *Language {
	return &Language{
		Dir:  "github.com/grpc-ecosystem/grpc-gateway",
		Name: "grpc-gateway",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "gateway_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/grpc-ecosystem/grpc-gateway:grpc-gateway"},
				WorkspaceExample: grpcGatewayWorkspaceTemplate,
				BuildExample:     grpcGatewayCompileExampleTemplate,
				Doc:              "Generates grpc-gateway *.go files",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "gateway_swagger_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/grpc-ecosystem/grpc-gateway:swagger"},
				WorkspaceExample: grpcGatewayWorkspaceTemplate,
				BuildExample:     grpcGatewayCompileExampleTemplate,
				Doc:              "Generates grpc-gateway swagger *.json files",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"}, // gRPC go lib rules fail on windows due to bad path
			},
			&Rule{
				Name:             "gateway_grpc_library",
				Kind:             "grpc",
				Implementation:   grpcGatewayLibraryRuleTemplate,
				WorkspaceExample: grpcGatewayWorkspaceTemplate,
				BuildExample:     grpcGatewayLibraryExampleTemplate,
				Doc:              "Generates grpc-gateway library files",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"}, // gRPC go lib rules fail on windows due to bad path
			},
		},
	}
}
