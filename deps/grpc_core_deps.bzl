"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_core_deps():
    """grpc_core dependency macro
    """
    com_google_protobuf()  # via com_github_grpc_grpc
    com_github_grpc_grpc()  # via <TOP>

def com_google_protobuf():
    _maybe(
        http_archive,
        name = "com_google_protobuf",
        sha256 = "7ed5fc41fe1614e551025f8e14b79b026a015b3ed337d38920c586f3ea35d818",
        strip_prefix = "protobuf-6b5d8db01fe47478e8d400f550e797e6230d464e",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/6b5d8db01fe47478e8d400f550e797e6230d464e.tar.gz",
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
