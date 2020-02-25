load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")

load(
    "//:deps.bzl",
    "build_bazel_rules_android",
    "com_google_protobuf",
    "com_google_protobuf_lite",
    "io_grpc_grpc_java",
    "MAVEN_SERVER_URLS",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def com_google_guava_guava_android(**kwargs):
    if "com_google_guava_guava_android" not in native.existing_rules():
        jvm_maven_import_external(
            name = "com_google_guava_guava_android",
            artifact = "com.google.guava:guava:27.0.1-android",
            server_urls = MAVEN_SERVER_URLS,
            artifact_sha256 = "caf0955aed29a1e6d149f85cfb625a89161b5cf88e0e246552b7ffa358204e28",
        )

def android_proto_compile(**kwargs):
    protobuf(**kwargs)
    com_google_protobuf_lite(**kwargs)

def android_grpc_compile(**kwargs):
    android_proto_compile(**kwargs)
    io_grpc_grpc_java(**kwargs)

def android_proto_library(**kwargs):
    android_proto_compile(**kwargs)
    build_bazel_rules_android(**kwargs)
    com_google_guava_guava_android(**kwargs)

def android_grpc_library(**kwargs):
    android_grpc_compile(**kwargs)
    android_proto_library(**kwargs)
