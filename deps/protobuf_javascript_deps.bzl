"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def protobuf_javascript_deps():
    bazel_skylib()  # via com_google_protobuf
    rules_pkg()  # via com_google_protobuf
    rules_python()  # via com_google_protobuf
    zlib()  # via com_google_protobuf
    com_google_protobuf()  # via <TOP>
    com_google_protobuf_javascript()  # via <TOP>

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

def rules_pkg():
    _maybe(
        http_archive,
        name = "rules_pkg",
        sha256 = "67e4dc634cb7237bc501fd03101b6b935c9f991c6f46d31a36b35ed5461f51b6",
        strip_prefix = "rules_pkg-4f8f6ed027c07b92e4ee5a8b4b51d8673fcc60ee",
        urls = [
            "https://github.com/bazelbuild/rules_pkg/archive/4f8f6ed027c07b92e4ee5a8b4b51d8673fcc60ee.tar.gz",
        ],
    )

def rules_python():
    _maybe(
        http_archive,
        name = "rules_python",
        sha256 = "8cc0ad31c8fc699a49ad31628273529ef8929ded0a0859a3d841ce711a9a90d5",
        strip_prefix = "rules_python-c7e068d38e2fec1d899e1c150e372f205c220e27",
        urls = [
            "https://github.com/bazelbuild/rules_python/archive/c7e068d38e2fec1d899e1c150e372f205c220e27.tar.gz",
        ],
    )

def zlib():
    _maybe(
        http_archive,
        name = "zlib",
        sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        strip_prefix = "zlib-1.2.11",
        urls = [
            "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
            "https://zlib.net/zlib-1.2.11.tar.gz",
        ],
        build_file = "@build_stack_rules_proto//third_party:zlib.BUILD",
    )

def com_google_protobuf():
    _maybe(
        http_archive,
        name = "com_google_protobuf",
        sha256 = "8b28fdd45bab62d15db232ec404248901842e5340299a57765e48abe8a80d930",
        strip_prefix = "protobuf-3.20.1",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/v3.20.1.tar.gz",
        ],
    )

def com_google_protobuf_javascript():
    _maybe(
        http_archive,
        name = "com_google_protobuf_javascript",
        sha256 = "392cef95222eb8ad7726c489ca9a02e46e93c404716ee80d7f9d3778975b4349",
        strip_prefix = "protobuf-javascript-3561b05cbf706aa14b0f8886c1167f402cf87b77",
        patches = [
            "@build_stack_rules_proto//third_party:protobuf-javascript-5.patch",
        ],
        patch_args = [
            "-p1",
        ],
        urls = [
            "https://github.com/protocolbuffers/protobuf-javascript/archive/3561b05cbf706aa14b0f8886c1167f402cf87b77.tar.gz",
        ],
    )
