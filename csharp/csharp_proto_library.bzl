load("//csharp:csharp_proto_compile.bzl", "csharp_proto_compile")
load("//:compile.bzl", "invoke_transitive")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def csharp_proto_library(**kwargs):
    kwargs["srcs"] = [invoke_transitive(csharp_proto_compile, "_pb", kwargs)]   
    kwargs["deps"] = [
        "@google.protobuf//:core",
        "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
    ]
    kwargs["verbose"] = None

    core_library(**kwargs)
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

