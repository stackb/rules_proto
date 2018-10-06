load("//cpp:compile.bzl", "cpp_proto_compile", "cpp_grpc_compile")

def cpp_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
    name_pb = name + "_pb"
    cpp_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        transitive = True,
    )
    native.cc_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
        ],
        includes = [name_pb],
        visibility = visibility,
    )


def cpp_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    cpp_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
    )
    native.cc_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
            "@com_github_grpc_grpc//:grpc++",
            "@com_github_grpc_grpc//:grpc++_reflection",
        ],
        # This seems magical to me.
        includes = [name_pb],
        visibility = visibility,
    )
