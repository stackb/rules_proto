"""proto_compile_gencopy.bzl provides the proto_compile_gencopy_run and proto_compile_gencopy_test rules.
"""

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

    # comprehend a mapping of relpath -> File
    srcfiles = {f.short_path[len(ctx.label.package):].lstrip("/"): f for f in ctx.files.srcs}

    for info in [dep[ProtoCompileInfo] for dep in ctx.attr.deps]:
        runfiles += info.outputs

        srcs = []  # list of string
        for f in info.outputs:
            if config.mode == "check":
                # if we are in 'check' mode, the src and dst cannot be the same file, so
                # make a copy of it...  but first, we need to find it in the srcs files!
                found = False
                for srcfilename, srcfile in srcfiles.items():
                    if srcfilename == f.basename:
                        replica = ctx.actions.declare_file(f.basename + ".actual", sibling = f)
                        _copy_file(ctx.actions, srcfile, replica)
                        runfiles.append(replica)
                        srcs.append(replica.short_path)
                        found = True
                        break
                if not found:
                    fail("could find matching source file for generated file %s in %r" % (f.short_path, srcfiles))

            else:
                srcs.append(f.short_path)

        config.packageConfigs.append(
            struct(
                targetLabel = str(info.label),
                targetPackage = info.label.package,
                generatedFiles = [f.short_path for f in info.outputs],
                sourceFiles = srcs,
            ),
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
            srcs = attr.label_list(
                doc = "The ProtoCompileInfo providers",
                allow_files = True,
            ),
        ),
        executable = True,
        test = is_test,
    )

proto_compile_gencopy_test = _proto_compile_gencopy_rule(True)
proto_compile_gencopy_run = _proto_compile_gencopy_rule(False)
