"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_gateway_ts_deps():
    """grpc_gateway_ts dependency macro
    """
    grpc_gateway_ts_darvin()  # via <TOP>
    grpc_gateway_ts_linux()  # via <TOP>
    grpc_gateway_ts_windows()  # via <TOP>

def grpc_gateway_ts_darvin():
    _maybe(
        http_archive,
        name = "grpc_gateway_ts_darvin",
        sha256 = "847349db0bcf0dc48ea8b9887de5bff2326561be3ad89b20c8b21f7163d1903d",
        urls = [
            "https://github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/releases/download/v1.1.1/protoc-gen-grpc-gateway-ts_1.1.1_Darwin_amd64.tar.gz",
        ],
        build_file_content = """
filegroup(
    name = "exe",
    srcs = ["protoc-gen-grpc-gateway-ts"],
    visibility = ["//visibility:public"],
)
""",
    )

def grpc_gateway_ts_linux():
    _maybe(
        http_archive,
        name = "grpc_gateway_ts_linux",
        sha256 = "48f69195f2e07ad04058e07561c8bc7253c8e3f758165301985b0109ae74057b",
        urls = [
            "https://github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/releases/download/v1.1.1/protoc-gen-grpc-gateway-ts_1.1.1_Linux_amd64.tar.gz",
        ],
        build_file_content = """
filegroup(
    name = "exe",
    srcs = ["protoc-gen-grpc-gateway-ts"],
    visibility = ["//visibility:public"],
)
""",
    )

def grpc_gateway_ts_windows():
    _maybe(
        http_archive,
        name = "grpc_gateway_ts_windows",
        sha256 = "ba769d45f99fded91329f68053a8755573e4fac9842ca11cd3864dd7f8425817",
        urls = [
            "https://github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts/releases/download/v1.1.1/protoc-gen-grpc-gateway-ts_1.1.1_Windows_amd64.tar.gz",
        ],
        build_file_content = """
filegroup(
    name = "exe",
    srcs = ["protoc-gen-grpc-gateway-ts.exe"],
    visibility = ["//visibility:public"],
)
""",
    )
