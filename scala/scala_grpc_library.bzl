load("//scala:scala_grpc_compile.bzl", "scala_grpc_compile")
load("@io_bazel_rules_scala//scala:scala.bzl", "scala_library")

def scala_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    scala_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create scala library
    scala_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        exports = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
    )

GRPC_DEPS = [
    "//external:io_bazel_rules_scala/dependency/com_google_protobuf/protobuf_java",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_fastparse",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_lenses",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_runtime",
    "//external:io_bazel_rules_scala/dependency/proto/google_instrumentation",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_context",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_core",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_netty",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_protobuf",
    "//external:io_bazel_rules_scala/dependency/proto/grpc_stub",
    "//external:io_bazel_rules_scala/dependency/proto/guava",
    "//external:io_bazel_rules_scala/dependency/proto/netty_buffer",
    "//external:io_bazel_rules_scala/dependency/proto/netty_codec",
    "//external:io_bazel_rules_scala/dependency/proto/netty_codec_http",
    "//external:io_bazel_rules_scala/dependency/proto/netty_codec_http2",
    "//external:io_bazel_rules_scala/dependency/proto/netty_codec_socks",
    "//external:io_bazel_rules_scala/dependency/proto/netty_common",
    "//external:io_bazel_rules_scala/dependency/proto/netty_handler",
    "//external:io_bazel_rules_scala/dependency/proto/netty_handler_proxy",
    "//external:io_bazel_rules_scala/dependency/proto/netty_resolver",
    "//external:io_bazel_rules_scala/dependency/proto/netty_transport",
    "//external:io_bazel_rules_scala/dependency/proto/opencensus_api",
    "//external:io_bazel_rules_scala/dependency/proto/opencensus_contrib_grpc_metrics",
    "//external:io_bazel_rules_scala/dependency/proto/scalapb_runtime_grpc",
]
