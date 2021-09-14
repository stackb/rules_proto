"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

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
