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
        sha256 = "33042aa893625ec5bf6d59bf38b3954e5558b7e549b1cb2eeee66cd2ccf8ab29",
        strip_prefix = "grpc.js-c938ee76ee462abf4f83d758f63d52f03fa24c7c",
        urls = [
            "https://github.com/stackb/grpc.js/archive/c938ee76ee462abf4f83d758f63d52f03fa24c7c.tar.gz",
        ],
    )
