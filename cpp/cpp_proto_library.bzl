load("//cpp:cpp_proto_compile.bzl", "cpp_proto_compile")

def cpp_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    cpp_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create cpp library
    native.cc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "//external:protobuf_clib",
        ],
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
    )
