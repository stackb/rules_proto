"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def closure_deps():
    """closure dependency macro
    """
    com_google_protobuf_javascript()  # via io_bazel_rules_closure
    io_bazel_rules_closure()  # via <TOP>

def com_google_protobuf_javascript():
    _maybe(
        http_archive,
        name = "com_google_protobuf_javascript",
        sha256 = "06fc35c7d35c48bdc99a6ab72211086532d1de2bc4ec28011cde607a4025ea95",
        strip_prefix = "protobuf-javascript-e1a52f9a897653985b0649cca17615cb1b0eb3b7",
        urls = [
            "https://github.com/protocolbuffers/protobuf-javascript/archive/e1a52f9a897653985b0649cca17615cb1b0eb3b7.tar.gz",
        ],
    )

def io_bazel_rules_closure():
    _maybe(
        http_archive,
        name = "io_bazel_rules_closure",
        sha256 = "00d492551233d7548ca2a983f4e19d6aabb0bc716957ade62d691baf1dcef374",
        strip_prefix = "rules_closure-42195b5ca136f78d28819ef486e3a7b02ad45146",
        urls = [
            "https://github.com/bazelbuild/rules_closure/archive/42195b5ca136f78d28819ef486e3a7b02ad45146.tar.gz",
        ],
    )
