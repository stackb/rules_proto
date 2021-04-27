load("//cmd/gencopy:gencopy.bzl", "gencopy_action", "gencopy_attrs", "gencopy_config")
load(":providers.bzl", "ProtoCompileInfo")

def _proto_compile_gencopy_impl(ctx):
    outputs = []
    for info in [dep[ProtoCompileInfo] for dep in ctx.attr.deps]:
        outputs += info.outputs

    script, runfiles = gencopy_action(ctx, gencopy_config(ctx), outputs)

    return [DefaultInfo(
        files = depset(outputs),
        runfiles = runfiles,
        executable = script,
    )]


def _proto_compile_gencopy_rule(is_test):
    return rule(
        implementation = _proto_compile_gencopy_impl,
        attrs = dict(
            gencopy_attrs,
            deps = attr.label_list(
                doc = "The ProtoCompileInfo providers",
                providers = [ProtoCompileInfo],
            ),
        ),
        executable = True,
        test = is_test,
    )

proto_compile_gencopy_test = _proto_compile_gencopy_rule(True)
proto_compile_gencopy_run = _proto_compile_gencopy_rule(False)
