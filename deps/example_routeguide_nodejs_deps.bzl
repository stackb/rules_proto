"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@build_bazel_rules_nodejs//:index.bzl", "npm_install")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def example_routeguide_nodejs_deps():
    """example_routeguide_nodejs dependency macro
    """
    build_bazel_rules_nodejs()  # via npm_example_routeguide_nodejs
    npm_example_routeguide_nodejs()  # via <TOP>

def build_bazel_rules_nodejs():
    _maybe(
        http_archive,
        name = "build_bazel_rules_nodejs",
        sha256 = "4501158976b9da216295ac65d872b1be51e3eeb805273e68c516d2eb36ae1fbb",
        urls = [
            "https://github.com/bazelbuild/rules_nodejs/releases/download/4.4.1/rules_nodejs-4.4.1.tar.gz",
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
