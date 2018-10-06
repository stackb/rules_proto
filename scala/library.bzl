load("//scala:compile.bzl", "scala_proto_compile", "scala_grpc_compile")
load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def scala_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    scala_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    scala_library(
        name = name,
        srcs = [name_pb],
        deps = ["//scala:proto_deps"],
        exports = [
            "//scala:proto_deps",
        ],
        visibility = visibility,
    )

def scala_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    scala_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    scala_library(
        name = name,
        srcs = [name_pb],
        deps = ["//scala:grpc_deps"],
        exports = [
            "//scala:grpc_deps",
        ],
        visibility = visibility,
    )
