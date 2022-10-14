"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_web_deps():
    """grpc_web dependency macro
    """
    com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_darwin_x86_64()  # via <TOP>
    com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_linux_x86_64()  # via <TOP>
    com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_windows_x86_64()  # via <TOP>

def com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_darwin_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_darwin_x86_64",
        executable = True,
        sha256 = "466ffe6d2096a2e09823ad02170a90a3e9f79d24094ec8ddcaf6c6d4e673aa2c",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.3.1/protoc-gen-grpc-web-1.3.1-darwin-x86_64",
        ],
    )

def com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_linux_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_linux_x86_64",
        executable = True,
        sha256 = "12d3cfefb22e077251e6d1fae2aeaafd7a66518898397c667ba69cdd1260e72a",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.3.1/protoc-gen-grpc-web-1.3.1-linux-x86_64",
        ],
    )

def com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_windows_x86_64():
    _maybe(
        http_file,
        name = "com_github_grpc_grpc_web_releases_download_v1_3_1_protoc_gen_grpc_web_1_3_1_windows_x86_64",
        executable = True,
        sha256 = "f7f3d3b8ddcc7f0f8e432e744768682c070491fc1dcacb922343ec8f39c0d115",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.3.1/protoc-gen-grpc-web-1.3.1-windows-x86_64.exe",
        ],
    )
