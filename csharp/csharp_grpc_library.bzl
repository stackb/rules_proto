load("//csharp:csharp_grpc_compile.bzl", "csharp_grpc_compile")
load("//:compile.bzl", "invoke_transitive")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def csharp_grpc_library(**kwargs):
    kwargs["srcs"] = [invoke_transitive(csharp_grpc_compile, "_pb", kwargs)]   
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

