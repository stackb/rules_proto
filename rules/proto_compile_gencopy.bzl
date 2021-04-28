load("//cmd/gencopy:gencopy.bzl", "gencopy_action", "gencopy_attrs", "gencopy_config")
load(":providers.bzl", "ProtoCompileInfo")


def _copy_file(actions, src, dst):
    """Copy a file to a new path destination
    Args:
      actions: the <ctx.actions> object
      src: the source file <File>
      dst: the destination path of the file
    Returns:
      <Generated File> for the copied file
    """
    actions.run_shell(
        mnemonic = "CopyFile",
        inputs = [src],
        outputs = [dst],
        command = "cp '{}' '{}'".format(src.path, dst.path),
        progress_message = "copying {} to {}".format(src.path, dst.path),
    )

def _proto_compile_gencopy_impl(ctx):

    config = gencopy_config(ctx)

    runfiles = []

    for info in [dep[ProtoCompileInfo] for dep in ctx.attr.deps]:
        srcs = []
        # here's the tricky bit: if we have a source file and generated file
        # that have the same relative path, the source file will get shadowed by
        # the generated one.  In the "update" case, that's not a problem.  In
        # the "check" case, it means we have to disambiguate the source files
        # with a different name.
        if config.mode == "check":
            for src in info.srcs:
                replica = ctx.actions.declare_file(src.basename+".actual", sibling=src)
                _copy_file(ctx.actions, src, replica)
                srcs.append(replica)
        else:
            srcs += info.srcs

        runfiles += info.outputs + srcs

        config.packageConfigs.append(
            struct(
                targetLabel = str(info.label),
                targetPackage = info.label.package,
                generatedFiles = [f.short_path for f in info.outputs],
                sourceFiles = [f.short_path for f in srcs],
            )
        )

    config_json, script, runfiles = gencopy_action(ctx, config, runfiles)

    return [DefaultInfo(
        files = depset([config_json]),
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
