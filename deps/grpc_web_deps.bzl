"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_web_deps():
    """grpc_web dependency macro
    """
    com_github_grpc_grpc_web()  # via <TOP>

def com_github_grpc_grpc_web():
    _maybe(
        http_archive,
        name = "com_github_grpc_grpc_web",
        sha256 = "d292df306b269ebf83fb53a349bbec61c07de4d628bd6a02d75ad3bd2f295574",
        strip_prefix = "grpc-web-1.3.1",
        urls = [
            "https://github.com/grpc/grpc-web/archive/1.3.1.tar.gz",
        ],
    )
