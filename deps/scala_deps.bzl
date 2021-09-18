"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def scala_deps():
    github_com_scalapb_scalapb_releases_download_v0_11_5_protoc_gen_scala_0_11_5_linux_x86_64_zip()  # via <TOP>
    github_com_scalapb_scalapb_releases_download_v0_11_5_protoc_gen_scala_0_11_5_osx_x86_64_zip()  # via <TOP>
    bazel_skylib()  # via io_bazel_rules_scala
    io_bazel_rules_scala()  # via <TOP>


def github_com_scalapb_scalapb_releases_download_v0_11_5_protoc_gen_scala_0_11_5_linux_x86_64_zip():
    _maybe(
        http_archive,
        name = "github_com_scalapb_scalapb_releases_download_v0_11_5_protoc_gen_scala_0_11_5_linux_x86_64_zip",
        sha256 = "1a58ae65d06e7894d7153812fc1c9ceeecab7891dcabeb33ee0b3e91e7502889",
        urls = [
            "https://github.com/scalapb/ScalaPB/releases/download/v0.11.5/protoc-gen-scala-0.11.5-linux-x86_64.zip",
        ],
        build_file_content = """
filegroup(
    name = "exe",
    srcs = ["protoc-gen-scala"],
    visibility = ["//visibility:public"],
)
""",
    )

def github_com_scalapb_scalapb_releases_download_v0_11_5_protoc_gen_scala_0_11_5_osx_x86_64_zip():
    _maybe(
        http_archive,
        name = "github_com_scalapb_scalapb_releases_download_v0_11_5_protoc_gen_scala_0_11_5_osx_x86_64_zip",
        sha256 = "312d802a999df87e7c046f2fef9213d0b3b8868b7771fa68ec8057ebf68397c0",
        urls = [
            "https://github.com/scalapb/ScalaPB/releases/download/v0.11.5/protoc-gen-scala-0.11.5-osx-x86_64.zip",
        ],
        build_file_content = """
filegroup(
    name = "exe",
    srcs = ["protoc-gen-scala"],
    visibility = ["//visibility:public"],
)
    """,
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
