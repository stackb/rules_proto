load("//:compile.bzl", "invoke_transitive")
load("//csharp:compile.bzl", "csharp_proto_compile", "csharp_grpc_compile")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def csharp_proto_library(**kwargs):
    kwargs["srcs"] = [invoke_transitive(csharp_proto_compile, "_pb", kwargs)]   
    kwargs["deps"] = [
        "@google.protobuf//:core",
        "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
    ]
    kwargs["verbose"] = None

    core_library(**kwargs)

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
