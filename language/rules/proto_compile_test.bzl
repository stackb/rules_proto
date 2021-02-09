load("@build_stack_rules_proto//internal:common.bzl", "ProtoCompileInfo")
load(
    "@build_stack_rules_proto//tools/gencopy:gencopy.bzl",
    "gencopy_action",
    "gencopy_attrs",
    "gencopy_config",
)

def _proto_compile_test_impl(ctx):
    outputs = []

    config = gencopy_config(ctx)

    for info in [dep[ProtoCompileInfo] for dep in ctx.attr.deps]:
        for [genrootdir, genfiles] in info.output_files.items():
            for f in genfiles:
                outputs.append(f)

    script, runfiles = gencopy_action(ctx, config, outputs)

    return [
        DefaultInfo(
            files = depset(outputs),
            runfiles = runfiles,
            executable = script,
        ),
    ]

def _proto_compile_rule(is_test):
    return rule(
        implementation = _proto_compile_test_impl,
        attrs = dict(
            gencopy_attrs,
            deps = attr.label_list(
                doc = "The ProtoCompileInfo provider rules",
                providers = [ProtoCompileInfo],
            ),
        ),
        executable = True,
        test = is_test,
    )

_proto_compile_test = _proto_compile_rule(True)
_proto_compile_run = _proto_compile_rule(False)

def proto_compile_test(**kwargs):
    proto_compile_rule = kwargs.pop("rule")
    srcs = kwargs.pop("srcs", [])
    name = kwargs.pop("name")

    out_name = name + "_out"
    update_target_label_name = "golden"
    update_name = "%s.%s" % (name, update_target_label_name)

    proto_compile_rule(name = out_name, **kwargs)

    _proto_compile_test(
        name = name,
        deps = [out_name],
        srcs = srcs,
        mode = "check",
        update_target_label_name = update_target_label_name,
    )

    _proto_compile_run(
        name = update_name,
        deps = [out_name],
        mode = "update",
        update_target_label_name = update_target_label_name,
    )
