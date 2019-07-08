load("//d:d_proto_compile.bzl", "d_proto_compile")
load("@io_bazel_rules_d//d:d.bzl", "d_library")

def d_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    d_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create d library
    d_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        imports = ["external/com_github_dcarp_protobuf_d/src", name_pb],
        visibility = kwargs.get("visibility"),
    )

PROTO_DEPS = [
    "@com_github_dcarp_protobuf_d//:protosrc",
    "@com_github_dcarp_protobuf_d//:protobuf",
]
