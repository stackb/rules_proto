"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_js_deps():
    com_github_stackb_grpc_js()  # via <TOP>

def com_github_stackb_grpc_js():
    _maybe(
        http_archive,
        name = "com_github_stackb_grpc_js",
        sha256 = "b2a52d483f8a5f5cda64e7c714a3a4fdb8af4a96268a8d39ce686f0ae2f9bc06",
        strip_prefix = "grpc.js-0b49f2138b98eb676f6bb55620358dec4ee8d9a2",
        urls = [
            "https://github.com/stackb/grpc.js/archive/0b49f2138b98eb676f6bb55620358dec4ee8d9a2.tar.gz",
        ],
    )
