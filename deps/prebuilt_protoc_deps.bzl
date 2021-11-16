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
        sha256 = "6003de742ea3fcf703cfec1cd4a3380fd143081a2eb0e559065563496af27807",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-linux-x86_64.zip",
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
        sha256 = "0decc6ce5beed07f8c20361ddeb5ac7666f09cf34572cca530e16814093f9c0c",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-osx-x86_64.zip",
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
        sha256 = "0decc6ce5beed07f8c20361ddeb5ac7666f09cf34572cca530e16814093f9c0c",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.6.1/protoc-3.6.1-win32.zip",
        ],
        build_file_content = """
filegroup(
    name = "protoc",
    srcs = ["bin/protoc.exe"],
    visibility = ["//visibility:public"],
)
""",
    )
