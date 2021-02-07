load("@build_stack_languages_proto//:proto_language.bzl", "ProtoLanguageInfo")
load(
    "@build_stack_languages_proto//tools/gencopy:gencopy.bzl",
    "gencopy_action",
    "gencopy_attrs",
    "gencopy_config",
)

def _proto_language_test_impl(ctx):
    outputs = []

    config = gencopy_config(ctx)

    for info in [dep[ProtoLanguageInfo] for dep in ctx.attr.deps]:
        outputs.append(info.markdown_file)

    script, runfiles = gencopy_action(ctx, config, outputs)

    return [
        DefaultInfo(
            files = depset(outputs),
            runfiles = runfiles,
            executable = script,
        ),
    ]

def _proto_language_rule(is_test):
    return rule(
        implementation = _proto_language_test_impl,
        attrs = dict(
            gencopy_attrs,
            deps = attr.label_list(
                doc = "The ProtoLanguageInfo provider languages",
                providers = [ProtoLanguageInfo],
            ),
        ),
        executable = True,
        test = is_test,
    )

_proto_language_test = _proto_language_rule(True)
_proto_language_run = _proto_language_rule(False)

def proto_language_test(**kwargs):
    deps = kwargs.pop("deps", [])
    srcs = kwargs.pop("srcs", [])
    name = kwargs.pop("name")

    update_target_label_name = "golden"
    update_name = "%s.%s" % (name, update_target_label_name)

    _proto_language_test(
        name = name,
        deps = deps,
        srcs = srcs,
        mode = "check",
        update_target_label_name = update_target_label_name,
    )

    _proto_language_run(
        name = update_name,
        deps = deps,
        mode = "update",
        update_target_label_name = update_target_label_name,
    )
