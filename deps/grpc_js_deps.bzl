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
        sha256 = "ac90cd56d28ea3841e97ad18fe50ac5dce04120e3b91116d5c3ffc21bc33c098",
        strip_prefix = "grpc.js-349d42b41fc52cb4c0a91662961b9ecff9201396",
        urls = [
            "https://github.com/stackb/grpc.js/archive/349d42b41fc52cb4c0a91662961b9ecff9201396.tar.gz",
        ],
    )
