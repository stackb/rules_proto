package main

var grpcGatewayWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//:repositories.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_gateway_repos="gateway_repos")

rules_proto_grpc_gateway_repos()

load("@grpc_ecosystem_grpc_gateway//:repositories.bzl", "go_repositories")

go_repositories()`)

var grpcGatewayLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:gateway_grpc_compile.bzl", "gateway_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    gateway_{{ .Rule.Kind }}_compile(
        name = name_pb,
        prefix_path = kwargs.get("importpath", ""),
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create go library
    go_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = kwargs.get("go_deps", []) + GRPC_DEPS,
        importpath = kwargs.get("importpath"),
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    "@com_github_golang_protobuf//descriptor:go_default_library",
    "@com_github_golang_protobuf//proto:go_default_library",
    "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
    "@org_golang_google_protobuf//runtime/protoimpl:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
    "@org_golang_google_grpc//codes:go_default_library",
    "@org_golang_google_grpc//grpclog:go_default_library",
    "@org_golang_google_grpc//metadata:go_default_library",
    "@org_golang_google_grpc//status:go_default_library",
    "@org_golang_x_net//context:go_default_library",
    "@grpc_ecosystem_grpc_gateway//runtime:go_default_library",
    "@grpc_ecosystem_grpc_gateway//utilities:go_default_library",
    "@go_googleapis//google/api:annotations_go_proto",
]`)

var grpcGatewayCompileExampleTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "api_gateway_grpc",
    deps = ["@rules_proto_grpc//{{ .Lang.Dir }}/example/api:api_proto"],
)`)

var grpcGatewayLibraryExampleTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:defs.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "api_gateway_library",
    importpath = "github.com/rules-proto-grpc/rules_proto_grpc/github.com/grpc-ecosystem/grpc-gateway/examples/api",
    deps = ["@rules_proto_grpc//{{ .Lang.Dir }}/example/api:api_proto"],
)`)

func makeGrpcGateway() *Language {
	return &Language{
		Dir:  "github.com/grpc-ecosystem/grpc-gateway",
		Name: "grpc-gateway",
		DisplayName: "grpc-gateway",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "gateway_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/grpc-ecosystem/grpc-gateway:grpc_gateway_plugin", "//go:grpc_go_plugin"},
				WorkspaceExample: grpcGatewayWorkspaceTemplate,
				BuildExample:     grpcGatewayCompileExampleTemplate,
				Doc:              "Generates grpc-gateway `.go` files",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "gateway_swagger_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//github.com/grpc-ecosystem/grpc-gateway:swagger_plugin"},
				WorkspaceExample: grpcGatewayWorkspaceTemplate,
				BuildExample:     grpcGatewayCompileExampleTemplate,
				Doc:              "Generates grpc-gateway swagger `.json` files",
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
				Attrs:            append(aspectProtoCompileAttrs, goProtoAttrs...),
				SkipTestPlatforms: []string{"windows"}, // gRPC go lib rules fail on windows due to bad path
			},
		},
	}
}
