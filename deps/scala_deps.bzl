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
        sha256 = "1ce86ffee65725300dc1f0017b7df89715c832de550137432dc1985d60a13155",
        strip_prefix = "rules_jvm_external-e6c1ff21e002bf97a7b1c07d63edd508a8dc9659",
        urls = [
            "https://github.com/bazelbuild/rules_jvm_external/archive/e6c1ff21e002bf97a7b1c07d63edd508a8dc9659.tar.gz",
        ],
    )

def bazel_skylib():
    _maybe(
        http_archive,
        name = "bazel_skylib",
        sha256 = "118e313990135890ee4cc8504e32929844f9578804a1b2f571d69b1dd080cfb8",
        strip_prefix = "bazel-skylib-1.5.0",
        urls = [
            "https://github.com/bazelbuild/bazel-skylib/archive/1.5.0.tar.gz",
        ],
    )

def io_bazel_rules_scala():
    _maybe(
        http_archive,
        name = "io_bazel_rules_scala",
        sha256 = "9a23058a36183a556a9ba7229b4f204d3e68c8c6eb7b28260521016b38ef4e00",
        strip_prefix = "rules_scala-6.4.0",
        urls = [
            "https://github.com/bazelbuild/rules_scala/archive/v6.4.0.tar.gz",
        ],
    )
