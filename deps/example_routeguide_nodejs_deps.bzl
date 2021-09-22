"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@build_bazel_rules_nodejs//:index.bzl", "npm_install")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def example_routeguide_nodejs_deps():
    build_bazel_rules_nodejs()  # via npm_example_routeguide_nodejs
    npm_example_routeguide_nodejs()  # via <TOP>

def build_bazel_rules_nodejs():
    _maybe(
        http_archive,
        name = "build_bazel_rules_nodejs",
        sha256 = "5c40083120eadec50a3497084f99bc75a85400ea727e82e0b2f422720573130f",
        urls = [
            "https://github.com/bazelbuild/rules_nodejs/releases/download/4.0.0-beta.0/rules_nodejs-4.0.0-beta.0.tar.gz",
        ],
    )

def npm_example_routeguide_nodejs():
    _maybe(
        npm_install,
        name = "npm_example_routeguide_nodejs",
        package_json = "//example/routeguide/nodejs:package.json",
        package_lock_json = "//example/routeguide/nodejs:package-lock.json",
        symlink_node_modules = False,
    )
