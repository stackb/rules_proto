load(
    "@build_stack_rules_proto//rules:java_grpc_compile.bzl",
    "java_grpc_compile",
)
load("@rules_java//java:defs.bzl", "java_library")

GRPC_DEPS = [  # From https://github.com/grpc/grpc-java/blob/3ce5df3f78e8fd4a619a791914087dd4c0562835/compiler/BUILD.bazel#L21-L27
    "@io_grpc_grpc_java//api",
    "@io_grpc_grpc_java//protobuf",
    "@io_grpc_grpc_java//stub",
    "@io_grpc_grpc_java//stub:javax_annotation",
    "@com_google_code_findbugs_jsr305//jar",
    "@com_google_guava_guava//jar",
    "@com_google_protobuf//:protobuf_java",
    "@com_google_protobuf//:protobuf_java_util",
]

def java_grpc_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"

    java_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    java_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        runtime_deps = ["@io_grpc_grpc_java//netty"],
        exports = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags", []),
    )
