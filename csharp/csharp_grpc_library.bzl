load("//csharp:csharp_grpc_compile.bzl", "csharp_grpc_compile")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def csharp_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")
    transitive = kwargs.get("transitive")

    name_pb = name + "_pb"
    csharp_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        transitive = transitive,
        verbose = verbose,
    )

    core_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "@google.protobuf//:netstandard1.0_core",
            "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
            "@grpc.core//:netstandard1.5_core",
            "@system.interactive.async//:netstandard2.0_core",
        ],
        visibility = visibility,
    )
