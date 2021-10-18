"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def closure_deps():
    io_bazel_rules_closure()  # via <TOP>

def io_bazel_rules_closure():
    _maybe(
        http_archive,
        name = "io_bazel_rules_closure",
        sha256 = "50096b6be0052055ba4f0577d8aa3d82adf077377ffa86e2b7a67a335442f01b",
        strip_prefix = "rules_closure-2c59208867759800a37d0f008c3a4398af4c0cb2",
        urls = [
            "https://github.com/bazelbuild/rules_closure/archive/2c59208867759800a37d0f008c3a4398af4c0cb2.tar.gz",
        ],
    )
