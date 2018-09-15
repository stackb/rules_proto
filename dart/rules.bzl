load("@//:rules.bzl", "proto_compile")

def dart_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//dart:dart"],
        **kwargs
    )

def grpc_dart_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//dart:dart"],
        **kwargs
    )


def grpc_dart_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    grpc_dart_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )
    # native.dart_library(
    #     name = name,
    #     srcs = [name_pb],
    #     deps = ["@//dart:grpc_deps"],
    #     # This magically adds REPOSITORY_NAME/PACKAGE_NAME/{name_pb} to dartPATH
    #     imports = [name_pb],
    #     visibility = visibility,
    # )

