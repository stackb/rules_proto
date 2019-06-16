load("//java:java_proto_compile.bzl", "java_proto_compile")

def java_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    java_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create java library
    native.java_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [Label("//java:proto_deps")],
        exports = [
            Label("//java:proto_deps"),
        ],
        visibility = kwargs.get("visibility"),
    )
