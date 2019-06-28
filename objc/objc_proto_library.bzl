load("//objc:objc_proto_compile.bzl", "objc_proto_compile")
def objc_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    objc_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create objc library
    native.objc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )

