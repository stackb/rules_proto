package main

var csharpProtoWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@io_bazel_rules_dotnet//dotnet:deps.bzl", "dotnet_repositories")

dotnet_repositories()

load(
    "@io_bazel_rules_dotnet//dotnet:defs.bzl",
    "core_register_sdk",
    "dotnet_register_toolchains",
    "dotnet_repositories_nugets",
)

dotnet_register_toolchains()
dotnet_repositories_nugets()

core_register_sdk()

load("@rules_proto_grpc//csharp/nuget:nuget.bzl", "nuget_rules_proto_grpc_packages")

nuget_rules_proto_grpc_packages()`)

var csharpGrpcWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@io_bazel_rules_dotnet//dotnet:deps.bzl", "dotnet_repositories")

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

dotnet_repositories()

load(
    "@io_bazel_rules_dotnet//dotnet:defs.bzl",
    "core_register_sdk",
    "dotnet_register_toolchains",
    "dotnet_repositories_nugets",
)

dotnet_register_toolchains()
dotnet_repositories_nugets()

core_register_sdk()

load("@rules_proto_grpc//csharp/nuget:nuget.bzl", "nuget_rules_proto_grpc_packages")

nuget_rules_proto_grpc_packages()`)

var csharpLibraryRuleTemplateString = `load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    {{ .Lang.Name }}_{{ .Rule.Kind }}_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )
`

var csharpProtoLibraryRuleTemplate = mustTemplate(csharpLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    core_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

PROTO_DEPS = [
    "@google.protobuf//:core",
    "@io_bazel_rules_dotnet//dotnet/stdlib.core:netstandard.dll",
]`)

var csharpGrpcLibraryRuleTemplate = mustTemplate(csharpLibraryRuleTemplateString + `
    # Create {{ .Lang.Name }} library
    core_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    "@google.protobuf//:core",
    "@grpc.core//:core",
    "@io_bazel_rules_dotnet//dotnet/stdlib.core:netstandard.dll",
]`)

func makeCsharp() *Language {
	return &Language{
		Dir:   "csharp",
		Name:  "csharp",
		DisplayName: "C#",
		Flags: commonLangFlags,
		Notes: mustTemplate(`Rules for generating C# protobuf and gRPC ` + "`.cs`" + ` files and libraries using standard Protocol Buffers and gRPC. Libraries are created with ` + "`core_library`" + ` from [rules_dotnet](https://github.com/bazelbuild/rules_dotnet)`),
		Rules: []*Rule{
			&Rule{
				Name:             "csharp_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//csharp:csharp_plugin"},
				WorkspaceExample: csharpProtoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates C# protobuf `.cs` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "csharp_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//csharp:csharp_plugin", "//csharp:grpc_csharp_plugin"},
				WorkspaceExample: csharpGrpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates C# protobuf+gRPC `.cs` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "csharp_proto_library",
				Kind:             "proto",
				Implementation:   csharpProtoLibraryRuleTemplate,
				WorkspaceExample: csharpProtoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates a C# protobuf library using `core_library` from `rules_dotnet`. Note that the library name must end in `.dll`",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "csharp_grpc_library",
				Kind:             "grpc",
				Implementation:   csharpGrpcLibraryRuleTemplate,
				WorkspaceExample: csharpGrpcWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates a C# protobuf+gRPC library using `core_library` from `rules_dotnet`. Note that the library name must end in `.dll`",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
