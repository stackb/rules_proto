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
        sha256 = "d488d8ecca98a4042442a4ae5f1ab0b614f896c0ebf6e3eafff363bcc51c6e62",
        strip_prefix = "bazel-lib-1.33.0",
        urls = [
            "https://github.com/aspect-build/bazel-lib/releases/download/v1.33.0/bazel-lib-v1.33.0.tar.gz",
        ],
    )

def aspect_rules_js():
    _maybe(
        http_archive,
        name = "aspect_rules_js",
        sha256 = "e3e6c3d42491e2938f4239a3d04259a58adc83e21e352346ad4ef62f87e76125",
        strip_prefix = "rules_js-1.30.0",
        urls = [
            "https://github.com/aspect-build/rules_js/releases/download/v1.30.0/rules_js-v1.30.0.tar.gz",
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
        sha256 = "4c3f34fff9f96ffc9c26635d8235a32a23a6797324486c7d23c1dfa477e8b451",
        strip_prefix = "rules_ts-1.4.5",
        urls = [
            "https://github.com/aspect-build/rules_ts/releases/download/v1.4.5/rules_ts-v1.4.5.tar.gz",
        ],
    )
