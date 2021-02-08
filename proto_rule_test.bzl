load("@build_stack_rules_proto//:proto_rule.bzl", "ProtoRuleInfo")
load(
    "@build_stack_rules_proto//tools/gencopy:gencopy.bzl",
    "gencopy_action",
    "gencopy_attrs",
    "gencopy_config",
)

def _proto_rule_test_impl(ctx):
    outputs = []

    config = gencopy_config(ctx)

    for info in [dep[ProtoRuleInfo] for dep in ctx.attr.deps]:
        outputs.append(info.bzl_file)
        outputs.append(info.deps_file)

    script, runfiles = gencopy_action(ctx, config, outputs)

    return [
        DefaultInfo(
            files = depset(outputs),
            runfiles = runfiles,
            executable = script,
        ),
    ]

def _proto_rule_rule(is_test):
    return rule(
        implementation = _proto_rule_test_impl,
        attrs = dict(
            gencopy_attrs,
            deps = attr.label_list(
                doc = "The ProtoRuleInfo provider rules",
                providers = [ProtoRuleInfo],
            ),
        ),
        executable = True,
        test = is_test,
    )

_proto_rule_test = _proto_rule_rule(True)
_proto_rule_run = _proto_rule_rule(False)

def proto_rule_test(**kwargs):
    deps = kwargs.pop("deps", [])
    srcs = kwargs.pop("srcs", [])
    name = kwargs.pop("name")

    update_target_label_name = "golden"
    update_name = "%s.%s" % (name, update_target_label_name)

    _proto_rule_test(
        name = name,
        deps = deps,
        srcs = srcs,
        mode = "check",
        target_package = kwargs.get("target_package", None),
        update_target_label_name = update_target_label_name,
    )

    _proto_rule_run(
        name = update_name,
        deps = deps,
        mode = "update",
        target_package = kwargs.get("target_package", None),
        update_target_label_name = update_target_label_name,
    )
