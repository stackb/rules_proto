"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def prebuilt_protoc_deps():
    prebuilt_protoc_linux()  # via <TOP>
    prebuilt_protoc_osx()  # via <TOP>
    prebuilt_protoc_windows()  # via <TOP>

def prebuilt_protoc_linux():
    _maybe(
        http_archive,
        name = "prebuilt_protoc_linux",
        sha256 = "3a0e900f9556fbcac4c3a913a00d07680f0fdf6b990a341462d822247b265562",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.20.1/protoc-3.20.1-linux-x86_64.zip",
        ],
        build_file_content = """
filegroup(
    name = "protoc",
    srcs = ["bin/protoc"],
    visibility = ["//visibility:public"],
)
""",
    )

def prebuilt_protoc_osx():
    _maybe(
        http_archive,
        name = "prebuilt_protoc_osx",
        sha256 = "b4f36b18202d54d343a66eebc9f8ae60809a2a96cc2d1b378137550bbe4cf33c",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.20.1/protoc-3.20.1-osx-x86_64.zip",
        ],
        build_file_content = """
filegroup(
    name = "protoc",
    srcs = ["bin/protoc"],
    visibility = ["//visibility:public"],
)
""",
    )

def prebuilt_protoc_windows():
    _maybe(
        http_archive,
        name = "prebuilt_protoc_windows",
        sha256 = "2291c634777242f3bf4891b082cebc6dd495ae621fbf751b27e800b83369a345",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.20.1/protoc-3.20.1-win32.zip",
        ],
        build_file_content = """
filegroup(
    name = "protoc",
    srcs = ["bin/protoc.exe"],
    visibility = ["//visibility:public"],
)
""",
    )
