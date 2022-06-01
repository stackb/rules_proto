"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_node_deps():
    build_bazel_rules_swift()  # via com_github_grpc_grpc
    bazel_skylib()  # via com_google_protobuf
    rules_pkg()  # via com_google_protobuf
    rules_python()  # via com_google_protobuf
    zlib()  # via com_google_protobuf
    com_google_protobuf()  # via com_github_grpc_grpc
    rules_jvm_external()  # via com_github_grpc_grpc
    com_github_grpc_grpc()  # via com_github_grpc_grpc_node_packages_grpc_tools_src
    com_github_grpc_grpc_node_packages_grpc_tools_src()  # via <TOP>

def build_bazel_rules_swift():
    _maybe(
        http_archive,
        name = "build_bazel_rules_swift",
        sha256 = "1f5499bb053736cda8905d89aac42e98011bbe9ca93b774a40c04759f045d7bf",
        strip_prefix = "rules_swift-dadd12190182530cf6f91ca7f9e70391644ce502",
        urls = [
            "https://github.com/bazelbuild/rules_swift/archive/dadd12190182530cf6f91ca7f9e70391644ce502.tar.gz",
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
        sha256 = "5361aad0b47621ee08ab16d5f997c4bb216eac999c0af9c7c7a7f6180d47e948",
        strip_prefix = "protobuf-7062d0a2d0075d5e7d5c294fd3984df67a976da3",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/7062d0a2d0075d5e7d5c294fd3984df67a976da3.tar.gz",
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
        sha256 = "e6c6b1ac9ba2257c93e49c98ef4fc96b2e2a1cdd90782a919f60e23fa8c2428b",
        strip_prefix = "grpc-5f759fcd1f602b38004b948b071f8b5726a9a4b1",
        urls = [
            "https://github.com/grpc/grpc/archive/5f759fcd1f602b38004b948b071f8b5726a9a4b1.tar.gz",
        ],
    )

def com_github_grpc_grpc_node_packages_grpc_tools_src():
    _maybe(
        http_archive,
        name = "com_github_grpc_grpc_node_packages_grpc_tools_src",
        sha256 = "7fbe9d04e45420c3c2e02456c0275fa8716fa894c48525b9a8f7db9ac0b4acb0",
        strip_prefix = "grpc-node-aeb42733d861883b82323e2dc6d1aba0e3a12aa0/packages/grpc-tools/src",
        urls = [
            "https://github.com/grpc/grpc-node/archive/aeb42733d861883b82323e2dc6d1aba0e3a12aa0.tar.gz",
        ],
        build_file_content = """
cc_library(
    name = "grpc_plugin_support",
    srcs = ["node_generator.cc"],
    hdrs = [
        "config.h",
        "config_protobuf.h",
        "generator_helpers.h",
        "node_generator.h",
        "node_generator_helpers.h",
    ],
    deps = ["//external:protobuf_clib"],
)

cc_binary(
    name = "grpc_node_plugin",
    srcs = ["node_plugin.cc"],
    visibility = ["//visibility:public"],
    deps = [":grpc_plugin_support"],
)
""",
    )
