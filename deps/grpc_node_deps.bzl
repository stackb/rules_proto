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
    com_google_protobuf()  # via com_github_grpc_grpc
    com_github_grpc_grpc()  # via com_github_grpc_grpc_node_packages_grpc_tools_src
    com_github_grpc_grpc_node_packages_grpc_tools_src()  # via <TOP>

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
        sha256 = "437068b8b777d3b339da94d3498f1dc20642ac9bfa76db43abdd522186b1542b",
        strip_prefix = "grpc-1.60.0",
        urls = [
            "https://github.com/grpc/grpc/archive/v1.60.0.tar.gz",
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
