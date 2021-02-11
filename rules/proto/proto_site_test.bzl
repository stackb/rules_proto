load("@build_stack_rules_proto//rules/proto:proto_rule.bzl", "ProtoRuleInfo")
load("@build_stack_rules_proto//rules/proto:proto_language.bzl", "ProtoLanguageInfo")
load(
    "@build_stack_rules_proto//tools/gencopy:gencopy.bzl",
    "gencopy_action",
    "gencopy_attrs",
    "gencopy_config",
)

def _proto_site_test_impl(ctx):
    outputs = []

    config = gencopy_config(ctx)

    for info in [dep[ProtoRuleInfo] for dep in ctx.attr.rules]:
        outputs.append(info.markdown_file)
    for info in [dep[ProtoLanguageInfo] for dep in ctx.attr.languages]:
        outputs.append(info.markdown_file)

    script, runfiles = gencopy_action(ctx, config, outputs)

    return [
        DefaultInfo(
            files = depset(outputs),
            runfiles = runfiles,
            executable = script,
        ),
    ]

def _proto_site_rule(is_test):
    return rule(
        implementation = _proto_site_test_impl,
        attrs = dict(
            gencopy_attrs,
            languages = attr.label_list(
                doc = "The ProtoLanguageInfo provider rules",
                providers = [ProtoLanguageInfo],
            ),
            rules = attr.label_list(
                doc = "The ProtoRuleInfo provider rules",
                providers = [ProtoRuleInfo],
            ),
        ),
        executable = True,
        test = is_test,
    )

_proto_site_test = _proto_site_rule(True)
_proto_site_run = _proto_site_rule(False)

def proto_site_test(**kwargs):
    rules = kwargs.pop("rules", [])
    languages = kwargs.pop("languages", [])
    srcs = kwargs.pop("srcs", [])
    name = kwargs.pop("name")

    update_target_label_name = "golden"
    update_name = "%s.%s" % (name, update_target_label_name)

    _proto_site_test(
        name = name,
        rules = rules,
        languages = languages,
        srcs = srcs,
        mode = "check",
        update_target_label_name = update_target_label_name,
    )

    _proto_site_run(
        name = update_name,
        rules = rules,
        languages = languages,
        mode = "update",
        update_target_label_name = update_target_label_name,
    )
