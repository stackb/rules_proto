"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_gateway_deps():
    com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_darwin_x86_64()  # via <TOP>
    com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_linux_x86_64()  # via <TOP>
    com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_windows_x86_64()  # via <TOP>

def com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_darwin_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_darwin_x86_64",
        executable = True,
        sha256 = "015932596f3ade410de21526a4823d22d3772b2796269e39e0836ca2c41578f8",
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.10.0/protoc-gen-grpc-gateway-v2.10.0-darwin-x86_64",
        ],
    )

def com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_linux_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_linux_x86_64",
        executable = True,
        sha256 = "5137dbf4643ccda82a75903adbaf5a5bedb6c9fb476bf449879345d1e93c6517",
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.10.0/protoc-gen-grpc-gateway-v2.10.0-linux-x86_64",
        ],
    )

def com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_windows_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_grpc_gateway_2_10_0_windows_x86_64",
        executable = True,
        sha256 = "835b61b8a1ca82ba794862ce32caa37f1bb67bd7684e0a47d08aec3a372a5c69",
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.10.0/protoc-gen-grpc-gateway-v2.10.0-windows-x86_64.exe",
        ],
    )
