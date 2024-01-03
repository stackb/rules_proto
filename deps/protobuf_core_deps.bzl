"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def protobuf_core_deps():
    """protobuf_core dependency macro
    """
    bazel_skylib()  # via com_google_protobuf
    rules_pkg()  # via com_google_protobuf
    rules_python()  # via com_google_protobuf
    zlib()  # via com_google_protobuf
    com_google_protobuf()  # via <TOP>

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

def rules_pkg():
    _maybe(
        http_archive,
        name = "rules_pkg",
        sha256 = "de4cf980e4c5eba24f3897016a71daec6b8d3c36f9ecdfe4e6dbcabb5017ade0",
        strip_prefix = "rules_pkg-ea8c75a15c4ac9562da29f3d9a633decb384d4a3",
        urls = [
            "https://github.com/bazelbuild/rules_pkg/archive/ea8c75a15c4ac9562da29f3d9a633decb384d4a3.tar.gz",
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
        sha256 = "d594b561fb41bf243233d8f411c7f2b7d913e5c9c1be4ca439baf7e48384c893",
        strip_prefix = "protobuf-f0dc78d7e6e331b8c6bb2d5283e06aa26883ca7c",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/f0dc78d7e6e331b8c6bb2d5283e06aa26883ca7c.tar.gz",
        ],
    )
