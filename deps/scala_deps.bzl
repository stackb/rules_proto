"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""


load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def scala_deps():
    """scala dependency macro
    """
    rules_jvm_external()  # via <TOP>
    bazel_skylib()  # via io_bazel_rules_scala
    io_bazel_rules_scala()  # via <TOP>


def rules_jvm_external():
    _maybe(
        http_archive,
        name = "rules_jvm_external",
        sha256 = "31701ad93dbfe544d597dbe62c9a1fdd76d81d8a9150c2bf1ecf928ecdf97169",
        strip_prefix = "rules_jvm_external-4.0",
        urls = [
            "https://github.com/bazelbuild/rules_jvm_external/archive/4.0.zip",
        ],
    )

def bazel_skylib():
    _maybe(
        http_archive,
        name = "bazel_skylib",
        sha256 = "ebdf850bfef28d923a2cc67ddca86355a449b5e4f38b0a70e584dc24e5984aa6",
        strip_prefix = "bazel-skylib-f80bc733d4b9f83d427ce3442be2e07427b2cc8d",
        urls = [
            "https://github.com/bazelbuild/bazel-skylib/archive/f80bc733d4b9f83d427ce3442be2e07427b2cc8d.tar.gz",
        ],
    )

def io_bazel_rules_scala():
    _maybe(
        http_archive,
        name = "io_bazel_rules_scala",
        sha256 = "0701ee4e1cfd59702d780acde907ac657752fbb5c7d08a0ec6f58ebea8cd0efb",
        strip_prefix = "rules_scala-2437e40131072cadc1628726775ff00fa3941a4a",
        urls = [
            "https://github.com/bazelbuild/rules_scala/archive/2437e40131072cadc1628726775ff00fa3941a4a.tar.gz",
        ],
    )
