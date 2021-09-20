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
        sha256 = "825da2c522405cb5fa6279b051c9f5a1b052c03938a3cf738df393aaada335aa",
        strip_prefix = "rules_closure-56cb92640e02fe9d354d46087211abfbd6300b06",
        urls = [
            "https://github.com/bazelbuild/rules_closure/archive/56cb92640e02fe9d354d46087211abfbd6300b06.tar.gz",
        ],
    )
