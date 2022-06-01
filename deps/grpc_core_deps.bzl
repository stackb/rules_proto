"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_core_deps():
    build_bazel_rules_swift()  # via com_github_grpc_grpc
    bazel_skylib()  # via com_google_protobuf
    rules_pkg()  # via com_google_protobuf
    rules_python()  # via com_google_protobuf
    zlib()  # via com_google_protobuf
    com_google_protobuf()  # via com_github_grpc_grpc
    rules_jvm_external()  # via com_github_grpc_grpc
    com_github_grpc_grpc()  # via <TOP>

def build_bazel_rules_swift():
    _maybe(
        http_archive,
        name = "build_bazel_rules_swift",
        sha256 = "7e361b49ff66486c6bd59a6a715a05b5038a2c48df866acd0044310be5363503",
        strip_prefix = "rules_swift-e52312c88e27f58f98f9c65526a44dcaba892863",
        urls = [
            "https://github.com/bazelbuild/rules_swift/archive/e52312c88e27f58f98f9c65526a44dcaba892863.tar.gz",
        ],
    )

def bazel_skylib():
    _maybe(
        http_archive,
        name = "bazel_skylib",
        sha256 = "ebdf850bfef28d923a2cc67ddca86355a449b5e4f38b0a70e584dc24e5984aa6",
        strip_prefix = "bazel-skylib-f80bc733d4b9f83d427ce3442be2e07427b2cc8d",
        urls = [
            "https://github.com/bazelbuild/bazel-skylib/archive/f80bc733d4b9f83d427ce3442be2e07427b2cc8d.tar.gz",
        ],
    )

def rules_pkg():
    _maybe(
        http_archive,
        name = "rules_pkg",
        sha256 = "de4cf980e4c5eba24f3897016a71daec6b8d3c36f9ecdfe4e6dbcabb5017ade0",
        strip_prefix = "rules_pkg-ea8c75a15c4ac9562da29f3d9a633decb384d4a3",
        urls = [
            "https://github.com/bazelbuild/rules_pkg/archive/ea8c75a15c4ac9562da29f3d9a633decb384d4a3.tar.gz",
        ],
    )

def rules_python():
    _maybe(
        http_archive,
        name = "rules_python",
        sha256 = "8cc0ad31c8fc699a49ad31628273529ef8929ded0a0859a3d841ce711a9a90d5",
        strip_prefix = "rules_python-c7e068d38e2fec1d899e1c150e372f205c220e27",
        urls = [
            "https://github.com/bazelbuild/rules_python/archive/c7e068d38e2fec1d899e1c150e372f205c220e27.tar.gz",
        ],
    )

def zlib():
    _maybe(
        http_archive,
        name = "zlib",
        sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        strip_prefix = "zlib-1.2.11",
        urls = [
            "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
            "https://zlib.net/zlib-1.2.11.tar.gz",
        ],
        build_file = "@build_stack_rules_proto//third_party:zlib.BUILD",
    )

def com_google_protobuf():
    _maybe(
        http_archive,
        name = "com_google_protobuf",
        sha256 = "8b28fdd45bab62d15db232ec404248901842e5340299a57765e48abe8a80d930",
        strip_prefix = "protobuf-3.20.1",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/v3.20.1.tar.gz",
        ],
    )

def rules_jvm_external():
    _maybe(
        http_archive,
        name = "rules_jvm_external",
        sha256 = "31701ad93dbfe544d597dbe62c9a1fdd76d81d8a9150c2bf1ecf928ecdf97169",
        strip_prefix = "rules_jvm_external-4.0",
        urls = [
            "https://github.com/bazelbuild/rules_jvm_external/archive/4.0.zip",
        ],
    )

def com_github_grpc_grpc():
    _maybe(
        http_archive,
        name = "com_github_grpc_grpc",
        sha256 = "b2f2620c762427bfeeef96a68c1924319f384e877bc0e084487601e4cc6e434c",
        strip_prefix = "grpc-1.42.0",
        urls = [
            "https://github.com/grpc/grpc/archive/v1.42.0.tar.gz",
        ],
    )
