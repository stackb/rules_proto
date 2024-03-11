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
    github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_darwin_aarch64()  # via <TOP>
    github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_darwin_x86_64()  # via <TOP>
    github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_linux_aarch64()  # via <TOP>
    github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_linux_x86_64()  # via <TOP>
    github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_windows_x86_64_exe()  # via <TOP>

def github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_darwin_aarch64():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_darwin_aarch64",
        executable = True,
        sha256 = "a12b759629b943ebac5528f58fac5039d9aa2fb7abb9e9684d1b481b35afbfc6",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-darwin-aarch64",
        ],
    )

def github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_darwin_x86_64():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_darwin_x86_64",
        executable = True,
        sha256 = "1fa3ef92194d06c03448a5cba82759e9773e43d8b188866a1f1d4fc23bb1ecb7",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-darwin-x86_64",
        ],
    )

def github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_linux_aarch64():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_linux_aarch64",
        executable = True,
        sha256 = "522e958568cdeabdd20ef3c97394fc067ff8e374a728c08b7410bf5de8f57783",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-linux-aarch64",
        ],
    )

def github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_linux_x86_64():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_linux_x86_64",
        executable = True,
        sha256 = "2e6e074497b221045a14d5a54e9fc910945bfdd1198b12b9fc23686a95671d64",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-linux-x86_64",
        ],
    )

def github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_windows_x86_64_exe():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_5_0_protoc_gen_grpc_web_1_5_0_windows_x86_64_exe",
        executable = True,
        sha256 = "c8f6191072d09344555fb12d45e82cff9f8b8f29200b0d6497680e696feaf8a1",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-windows-x86_64.exe",
        ],
    )
