"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_gateway_openapiv2_deps():
    """grpc_gateway_openapiv2 dependency macro
    """
    com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_darwin_x86_64()  # via <TOP>
    com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_linux_x86_64()  # via <TOP>
    com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_windows_x86_64()  # via <TOP>

def com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_darwin_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_darwin_x86_64",
        executable = True,
        sha256 = "ec2dde842ac14a7fc3d9704863b04f8c878255428b9639e0d2cbbb9bf47280cf",
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.10.0/protoc-gen-openapiv2-v2.10.0-darwin-x86_64",
        ],
    )

def com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_linux_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_linux_x86_64",
        executable = True,
        sha256 = "71f5d666a2af33817b8621d42cb7cf0f695fd4b4690a411951a547532f28f7c7",
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.10.0/protoc-gen-openapiv2-v2.10.0-linux-x86_64",
        ],
    )

def com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_windows_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_ecosystem_grpc_gateway_releases_download_v2_10_0_protoc_gen_openapiv2_2_10_0_windows_x86_64",
        executable = True,
        sha256 = "ae16a52026310bd9eb07b5ab33a3a72a9471ae7c62fcc0dcf6e36bd85281a8af",
        urls = [
            "https://github.com/grpc-ecosystem/grpc-gateway/releases/download/v2.10.0/protoc-gen-openapiv2-v2.10.0-windows-x86_64.exe",
        ],
    )
