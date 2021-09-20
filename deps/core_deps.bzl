"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def core_deps():
    io_bazel_rules_go()  # via bazel_gazelle
    bazel_gazelle()  # via <TOP>
    rules_proto()  # via <TOP>

def io_bazel_rules_go():
    _maybe(
        http_archive,
        name = "io_bazel_rules_go",
        sha256 = "38171ce619b2695fa095427815d52c2a115c716b15f4cd0525a88c376113f584",
        strip_prefix = "rules_go-0.28.0",
        urls = [
            "https://github.com/bazelbuild/rules_go/archive/v0.28.0.tar.gz",
        ],
    )

def bazel_gazelle():
    _maybe(
        http_archive,
        name = "bazel_gazelle",
        sha256 = "7359a9a7071dff5343a52626fce1aae2a78936d3004ef89d038daaefd3fd6608",
        strip_prefix = "bazel-gazelle-e4496b956eb46bdf8c9bf95b8d1d85e3a086c4be",
        urls = [
            "https://github.com/bazelbuild/bazel-gazelle/archive/e4496b956eb46bdf8c9bf95b8d1d85e3a086c4be.tar.gz",
        ],
    )

def rules_proto():
    _maybe(
        http_archive,
        name = "rules_proto",
        sha256 = "9fc210a34f0f9e7cc31598d109b5d069ef44911a82f507d5a88716db171615a8",
        strip_prefix = "rules_proto-f7a30f6f80006b591fa7c437fe5a951eb10bcbcf",
        urls = [
            "https://github.com/bazelbuild/rules_proto/archive/f7a30f6f80006b591fa7c437fe5a951eb10bcbcf.tar.gz",
        ],
    )
