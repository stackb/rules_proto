"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""


load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def js_core_deps():
    """js_core dependency macro
    """
    com_google_protobuf_javascript()  # via <TOP>


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
