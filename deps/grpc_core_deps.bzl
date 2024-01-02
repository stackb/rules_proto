"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_core_deps():
    """grpc_core dependency macro
    """
    rules_pkg()  # via com_google_protobuf
    com_google_protobuf()  # via com_github_grpc_grpc
    com_github_grpc_grpc()  # via <TOP>

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

def com_github_grpc_grpc():
    _maybe(
        http_archive,
        name = "com_github_grpc_grpc",
        sha256 = "17e4e1b100657b88027721220cbfb694d86c4b807e9257eaf2fb2d273b41b1b1",
        strip_prefix = "grpc-1.54.3",
        urls = [
            "https://github.com/grpc/grpc/archive/v1.54.3.tar.gz",
        ],
    )
