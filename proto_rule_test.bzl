load(
    "@bazel_skylib//lib:shell.bzl",
    "shell",
)
load("//:proto_rule.bzl", "ProtoRuleInfo")

def _proto_rule_test_impl(ctx):
    go = go_context(ctx)

    return [
        ProtoRuleInfo(
            name = ctx.attr.name,
            rule = rule,
            bzl_file = rule_bzl,
            build_file = rule_build,
            workspace_file = rule_workspace,
            test_file = rule_test,
        ),
        DefaultInfo(
            files = depset(outputs),
        ),
    ]

proto_rule_test = rule(
    implementation = _proto_rule_test_impl,
    attrs = {
        "data": attr.label_list(allow_files = True),
        "deps": attr.label_list(
            providers = [ProtoRuleInfo],
            allow_files = True,
        ),
        "_go_context_data": attr.label(
            default = "@io_bazel_rules_go//:go_context_data",
        ),
    },
    toolchains = ["@io_bazel_rules_go//go:toolchain"],
    test_only = True,
)
