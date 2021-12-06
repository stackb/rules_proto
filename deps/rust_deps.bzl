"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def rust_deps():
    rules_rust()  # via <TOP>

def rules_rust():
    _maybe(
        http_archive,
        name = "rules_rust",
        sha256 = "608ac74a2892af88cd2ddbf51a8c5d9586641479deda3397ba7bf35743b6d0b7",
        strip_prefix = "rules_rust-20f4ff5ef691de251da7253dc3fafc2ab9add550",
        urls = [
            "https://github.com/bazelbuild/rules_rust/archive/20f4ff5ef691de251da7253dc3fafc2ab9add550.tar.gz",
        ],
    )
