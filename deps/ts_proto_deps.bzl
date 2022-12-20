"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@build_bazel_rules_nodejs//:index.bzl", "npm_install", "yarn_install")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def ts_proto_deps():
    """ts_proto dependency macro
    """
    rules_nodejs()  # via build_bazel_rules_nodejs
    build_bazel_rules_nodejs()  # via npm_ts_proto
    npm_ts_proto()  # via <TOP>
    npm_tsc()  # via <TOP>

def rules_nodejs():
    _maybe(
        http_archive,
        name = "rules_nodejs",
        sha256 = "08337d4fffc78f7fe648a93be12ea2fc4e8eb9795a4e6aa48595b66b34555626",
        urls = [
            "https://github.com/bazelbuild/rules_nodejs/releases/download/5.8.0/rules_nodejs-core-5.8.0.tar.gz",
        ],
    )

def build_bazel_rules_nodejs():
    _maybe(
        http_archive,
        name = "build_bazel_rules_nodejs",
        sha256 = "dcc55f810142b6cf46a44d0180a5a7fb923c04a5061e2e8d8eb05ccccc60864b",
        urls = [
            "https://github.com/bazelbuild/rules_nodejs/releases/download/5.8.0/rules_nodejs-5.8.0.tar.gz",
        ],
    )

def npm_ts_proto():
    _maybe(
        npm_install,
        name = "npm_ts_proto",
        package_json = "@build_stack_rules_proto//plugin/stephenh/ts-proto:package.json",
        package_lock_json = "@build_stack_rules_proto//plugin/stephenh/ts-proto:package-lock.json",
        symlink_node_modules = False,
    )

def npm_tsc():
    _maybe(
        yarn_install,
        name = "npm_tsc",
        package_json = "@build_stack_rules_proto//rules/ts:package.json",
        yarn_lock = "@build_stack_rules_proto//rules/ts:yarn.lock",
        frozen_lockfile = True,
    )
