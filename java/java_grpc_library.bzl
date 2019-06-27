load("//java:java_grpc_compile.bzl", "java_grpc_compile")

def java_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    java_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create java library
    native.java_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [Label("//java:grpc_deps")],
        exports = [
            Label("//java:grpc_deps"),
        ],
        visibility = kwargs.get("visibility"),
    )
