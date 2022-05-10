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
        sha256 = "4a3b26d1ebb9c1d23e933694a6669295f6a39ddc64c3db2adf671f0a6026f82e",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.13.0/protoc-3.13.0-linux-x86_64.zip",
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
        sha256 = "a201954cc7d1a309b5f4feacd23a0abcf3ffc20eb15e79c9a0856a5804f6c34c",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.13.0/protoc-3.13.0-osx-x86_64.zip",
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
        sha256 = "da7ac5d046810ee44c13bd92c6bc034763d483b506e697baf278e2769730716c",
        urls = [
            "https://github.com/google/protobuf/releases/download/v3.13.0/protoc-3.13.0-win32.zip",
        ],
        build_file_content = """
filegroup(
    name = "protoc",
    srcs = ["bin/protoc.exe"],
    visibility = ["//visibility:public"],
)
""",
    )
