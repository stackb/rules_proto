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
        sha256 = "a3944abe45f51c04ab3f567d26d6b7c15f9726349a8bb801b26d6b8495637bb9",
        strip_prefix = "grpc.js-84d000d03910220625e6c7c0ad434d598e0b6cd9",
        urls = [
            "https://github.com/stackb/grpc.js/archive/84d000d03910220625e6c7c0ad434d598e0b6cd9.tar.gz",
        ],
    )
