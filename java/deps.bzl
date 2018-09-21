load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

PLUGIN_VERSION = "1.9.0"

def java_grpc_library_deps():
    existing = native.existing_rules()

    if "io_grpc_grpc_java" not in existing:
        http_archive(
            name = "io_grpc_grpc_java",
            urls = ["https://github.com/grpc/grpc-java/archive/v1.15.0.tar.gz"],
            strip_prefix = "grpc-java-1.15.0",
            sha256 = "8a131e773b1c9c0442e606b7fc85d7fc6739659281589d01bd917ceda218a1c7",
        )
