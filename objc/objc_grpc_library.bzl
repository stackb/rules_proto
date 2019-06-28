load("//objc:objc_grpc_compile.bzl", "objc_grpc_compile")
def objc_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    objc_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create objc library
    native.objc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_google_protobuf//:protobuf_objc",
            "@com_github_grpc_grpc//:grpc++",
        ]
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )

