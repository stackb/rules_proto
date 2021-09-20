"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def nodejs_deps():
    build_bazel_rules_nodejs()  # via <TOP>

def build_bazel_rules_nodejs():
    _maybe(
        http_archive,
        name = "build_bazel_rules_nodejs",
        sha256 = "482741b49b730b4055e5bb3936b4fe97e27365e917d1e4d442d5b71a6180aaf2",
        strip_prefix = "rules_nodejs-4.2.0",
        urls = [
            "https://github.com/bazelbuild/rules_nodejs/archive/4.2.0.tar.gz",
        ],
    )
