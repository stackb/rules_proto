package main

var csharpProtoLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")
{{ .Rule.Name }}()

load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "dotnet_register_toolchains", "dotnet_repositories")
dotnet_register_toolchains("host")
#dotnet_register_toolchains(dotnet_version="4.2.3")
dotnet_repositories()

load("@build_stack_rules_proto//csharp/nuget:packages.bzl", nuget_packages = "packages")
nuget_packages()

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_protobuf_packages")
nuget_protobuf_packages()
`)

var csharpGrpcLibraryUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")
{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "dotnet_register_toolchains", "dotnet_repositories")

dotnet_register_toolchains("host")
#dotnet_register_toolchains(dotnet_version="4.2.3")

dotnet_repositories()

load("@io_bazel_rules_dotnet//csharp/nuget:packages.bzl", nuget_packages = "packages")
nuget_packages()

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_protobuf_packages")
load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_grpc_packages")

nuget_protobuf_packages()
nuget_grpc_packages()
`)

var csharpProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("//:compile.bzl", "invoke_transitive")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def {{ .Rule.Name }}(**kwargs):
    kwargs["srcs"] = [invoke_transitive({{ .Rule.Base }}_{{ .Rule.Kind }}_compile, "_pb", kwargs)]   
    kwargs["deps"] = [
        "@google.protobuf//:core",
        "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
    ]
    kwargs["verbose"] = None

    core_library(**kwargs)
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
`)

var csharpGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Rule.Base}}_{{ .Rule.Kind }}_compile.bzl", "{{ .Rule.Base }}_{{ .Rule.Kind }}_compile")
load("//:compile.bzl", "invoke_transitive")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def {{ .Rule.Name }}(**kwargs):
    kwargs["srcs"] = [invoke_transitive({{ .Rule.Base }}_{{ .Rule.Kind }}_compile, "_pb", kwargs)]   
    kwargs["deps"] = [
        "@google.protobuf//:core",
        "@grpc.core//:core",
        "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
        "@system.interactive.async//:core",
    ]
    kwargs["verbose"] = None

    core_library(**kwargs)
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
`)

func makeCsharp() *Language {
	return &Language{
		Dir:  "csharp",
		Name: "csharp",
		Rules: []*Rule{
			&Rule{
				Name:           "csharp_proto_compile",
				Base:           "csharp",
				Kind:           "proto",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//csharp:csharp"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates csharp protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "csharp_grpc_compile",
				Base:           "csharp",
				Kind:           "grpc",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//csharp:csharp", "//csharp:grpc_csharp"},
				Usage:          grpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates csharp protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "csharp_proto_library",
				Base:           "csharp",
				Kind:           "proto",
				Usage:          csharpProtoLibraryUsageTemplate,
				Example:        protoLibraryExampleTemplate,
				Implementation: csharpProtoLibraryRuleTemplate,
				Doc:            "Generates csharp protobuf library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "csharp_grpc_library",
				Base:           "csharp",
				Kind:           "grpc",
				Implementation: csharpGrpcLibraryRuleTemplate,
				Usage:          csharpGrpcLibraryUsageTemplate,
				Example:        grpcLibraryExampleTemplate,
				Doc:            "Generates csharp protobuf+gRPC library",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
