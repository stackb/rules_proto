"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def ts_proto_deps():
    """ts_proto dependency macro
    """
    aspect_bazel_lib()  # via aspect_rules_ts
    aspect_rules_js()  # via aspect_rules_ts
    rules_nodejs()  # via aspect_rules_ts
    aspect_rules_ts()  # via <TOP>

def aspect_bazel_lib():
    _maybe(
        http_archive,
        name = "aspect_bazel_lib",
        sha256 = "a7bfc7aed7b86a4caaba382116e0214ebbaa623f393a9e716d87a3e1bab29d78",
        strip_prefix = "bazel-lib-1.19.0",
        urls = [
            "https://github.com/aspect-build/bazel-lib/archive/refs/tags/v1.19.0.tar.gz",
        ],
    )

def aspect_rules_js():
    _maybe(
        http_archive,
        name = "aspect_rules_js",
        sha256 = "66ecc9f56300dd63fb86f11cfa1e8affcaa42d5300e2746dba08541916e913fd",
        strip_prefix = "rules_js-1.13.0",
        urls = [
            "https://github.com/aspect-build/rules_js/archive/refs/tags/v1.13.0.tar.gz",
        ],
    )

def rules_nodejs():
    _maybe(
        http_archive,
        name = "rules_nodejs",
        sha256 = "08337d4fffc78f7fe648a93be12ea2fc4e8eb9795a4e6aa48595b66b34555626",
        urls = [
            "https://github.com/bazelbuild/rules_nodejs/releases/download/5.8.0/rules_nodejs-core-5.8.0.tar.gz",
        ],
    )

def aspect_rules_ts():
    _maybe(
        http_archive,
        name = "aspect_rules_ts",
        sha256 = "e81f37c4fe014fc83229e619360d51bfd6cb8ac405a7e8018b4a362efa79d000",
        strip_prefix = "rules_ts-1.0.4",
        urls = [
            "https://github.com/aspect-build/rules_ts/archive/refs/tags/v1.0.4.tar.gz",
        ],
    )
