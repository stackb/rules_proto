"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

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
        sha256 = "9a37844f00eab4236d8ddbe9844d52c87fda58aa2631d05d9961a07976edb9c0",
        strip_prefix = "rules_go-0.32.0",
        urls = [
            "https://github.com/bazelbuild/rules_go/archive/v0.32.0.tar.gz",
        ],
    )

def bazel_gazelle():
    _maybe(
        http_archive,
        name = "bazel_gazelle",
        sha256 = "f09b77a3fa22ea467b98cbdd0387573705076dc7463da5481a5d5fd37c9deae6",
        strip_prefix = "bazel-gazelle-1dbcd58297322ddeeafbfd006b288cede1892352",
        urls = [
            "https://github.com/bazelbuild/bazel-gazelle/archive/1dbcd58297322ddeeafbfd006b288cede1892352.tar.gz",
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
