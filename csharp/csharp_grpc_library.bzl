load("//csharp:csharp_grpc_compile.bzl", "csharp_grpc_compile")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def csharp_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    csharp_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create csharp library
    core_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@google.protobuf//:netstandard1.0_core",
            "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
            "@grpc.core//:netstandard1.5_core",
            "@system.interactive.async//:netstandard2.0_core",
        ],
        visibility = kwargs.get("visibility"),
    )
