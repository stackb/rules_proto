load("//swift:compile.bzl", "swift_proto_compile", "swift_grpc_compile")
load("@build_bazel_rules_swift//swift:swift.bzl", "swift_library")

def swift_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    swift_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    swift_library(
        name = name,
        srcs = [name_pb],
        #deps = ["//swift:proto_deps"],
        visibility = visibility,
    )

def swift_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    swift_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    swift_library(
        name = name,
        srcs = [name_pb],
        #deps = ["//swift:grpc_deps"],
        visibility = visibility,
    )
