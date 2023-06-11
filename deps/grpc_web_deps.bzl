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
    github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_darwin_aarch64()  # via <TOP>
    github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_darwin_x86_64()  # via <TOP>
    github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_linux_aarch64()  # via <TOP>
    github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_linux_x86_64()  # via <TOP>
    github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_windows_x86_64_exe()  # via <TOP>

def github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_darwin_aarch64():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_darwin_aarch64",
        executable = True,
        sha256 = "87263950cd36ec875c86b1e50625215727264c5495d6625ddf9e4f79aeef727e",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-darwin-aarch64",
        ],
    )

def github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_darwin_x86_64():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_darwin_x86_64",
        executable = True,
        sha256 = "6b73e8e9ef2deb114d39c9eea177ff8450d92e7154b5e47dea668a43499a2383",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-darwin-x86_64",
        ],
    )

def github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_linux_aarch64():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_linux_aarch64",
        executable = True,
        sha256 = "db3b1e18284a157b1361ecb79502e5f4623a212daf7185f17e35bd8e239055a4",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-linux-aarch64",
        ],
    )

def github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_linux_x86_64():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_linux_x86_64",
        executable = True,
        sha256 = "5e82c3f1f435e176c94b94de9669911ab3bfb891608b7e80adff358f777ff857",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-linux-x86_64",
        ],
    )

def github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_windows_x86_64_exe():
    _maybe(
        http_file,
        name = "github_com_grpc_grpc_web_releases_download_1_4_2_protoc_gen_grpc_web_1_4_2_windows_x86_64_exe",
        executable = True,
        sha256 = "3a0fc44662cb89a5c101b632e3e8841d04d9bea3103512deae82591e2acdff93",
        urls = [
            "https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-windows-x86_64.exe",
        ],
    )
