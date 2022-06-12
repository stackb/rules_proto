"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@build_bazel_rules_nodejs//:index.bzl", "npm_install", "yarn_install")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def ts_proto_deps():
    """ts_proto dependency macro
    """
    npm_ts_proto()  # via <TOP>
    npm_tsc()  # via <TOP>

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
