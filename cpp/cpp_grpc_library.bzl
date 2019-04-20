load("//cpp:cpp_grpc_compile.bzl", "cpp_grpc_compile")

def cpp_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    cpp_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
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
