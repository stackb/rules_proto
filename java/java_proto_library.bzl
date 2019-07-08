load("//java:java_proto_compile.bzl", "java_proto_compile")

def java_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    java_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create java library
    native.java_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        exports = PROTO_DEPS,
        visibility = kwargs.get("visibility"),
    )

PROTO_DEPS = [
    "@com_google_guava_guava//jar",
    "@com_google_protobuf//:protobuf_java",
    "@javax_annotation_javax_annotation_api//jar",
]
