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
        transitive = True,
    )

    native.cc_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "@com_google_protobuf//:protobuf_lite",
            "@com_github_grpc_grpc//:grpc++_codegen_proto",
        ],
        # This seems magical to me.
        includes = [name_pb],
        visibility = visibility,
    )
