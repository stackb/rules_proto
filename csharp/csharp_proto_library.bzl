load("//csharp:csharp_proto_compile.bzl", "csharp_proto_compile")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def csharp_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")
    transitive = kwargs.get("transitive", True)

    name_pb = name + "_pb"
    csharp_proto_compile(
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
            "@google.protobuf//:core",
            "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
        ],
        visibility = visibility,
    )
