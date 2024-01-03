"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_node_deps():
    """grpc_node dependency macro
    """
    bazel_skylib()  # via com_google_protobuf
    rules_pkg()  # via com_google_protobuf
    rules_python()  # via com_google_protobuf
    zlib()  # via com_google_protobuf
    com_google_protobuf()  # via com_github_grpc_grpc
    com_github_grpc_grpc()  # via com_github_grpc_grpc_node_packages_grpc_tools_src
    com_github_grpc_grpc_node_packages_grpc_tools_src()  # via <TOP>

def bazel_skylib():
    _maybe(
        http_archive,
        name = "bazel_skylib",
        sha256 = "118e313990135890ee4cc8504e32929844f9578804a1b2f571d69b1dd080cfb8",
        strip_prefix = "bazel-skylib-1.5.0",
        urls = [
            "https://github.com/bazelbuild/bazel-skylib/archive/1.5.0.tar.gz",
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
        sha256 = "e85ae30de33625a63eca7fc40a94fea845e641888e52f32b6beea91e8b1b2793",
        strip_prefix = "rules_python-0.27.1",
        urls = [
            "https://github.com/bazelbuild/rules_python/archive/0.27.1.tar.gz",
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
        sha256 = "d594b561fb41bf243233d8f411c7f2b7d913e5c9c1be4ca439baf7e48384c893",
        strip_prefix = "protobuf-f0dc78d7e6e331b8c6bb2d5283e06aa26883ca7c",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/f0dc78d7e6e331b8c6bb2d5283e06aa26883ca7c.tar.gz",
        ],
    )

def com_github_grpc_grpc():
    _maybe(
        http_archive,
        name = "com_github_grpc_grpc",
        sha256 = "17e4e1b100657b88027721220cbfb694d86c4b807e9257eaf2fb2d273b41b1b1",
        strip_prefix = "grpc-1.54.3",
        urls = [
            "https://github.com/grpc/grpc/archive/v1.54.3.tar.gz",
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
