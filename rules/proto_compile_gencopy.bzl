"""proto_compile_gencopy.bzl provides the proto_compile_gencopy_run and proto_compile_gencopy_test rules.
"""

load("//cmd/gencopy:gencopy.bzl", "gencopy_action", "gencopy_attrs", "gencopy_config")
load(":providers.bzl", "ProtoCompileInfo")

def _proto_compile_gencopy_run_impl(ctx):
    config = gencopy_config(ctx)

    runfiles = []
    for info in [dep[ProtoCompileInfo] for dep in ctx.attr.deps]:
        # List[String]: names of files that represent the source files.  In an
        # update, these are the target filenames of a file copy operation.
        source_files = []

        # List[String]: names of files that represent the generated files.  In
        # an update, these are the source filenames of a file copy operation
        # (although the file itself was generated by proto_compile).
        generated_files = []
        for rel, generated_file in info.output_files_by_rel_path.items():
            runfiles.append(generated_file)
            source_files.append(rel)
            generated_files.append(generated_file.short_path)

        config.packageConfigs.append(
            struct(
                targetLabel = str(info.label),
                targetPackage = info.label.package,
                targetWorkspaceRoot = info.label.workspace_root,
                generatedFiles = generated_files,
                sourceFiles = source_files,
            ),
        )

    config_json, script, runfiles = gencopy_action(ctx, config, runfiles)

    return [DefaultInfo(
        files = depset([config_json]),
        runfiles = runfiles,
        executable = script,
    )]

proto_compile_gencopy_run = rule(
    implementation = _proto_compile_gencopy_run_impl,
    attrs = dict(
        gencopy_attrs,
        deps = attr.label_list(
            doc = "The ProtoCompileInfo providers",
            providers = [ProtoCompileInfo],
        ),
    ),
    executable = True,
    test = False,
)

def _proto_compile_gencopy_test_impl(ctx):
    config = gencopy_config(ctx)

    runfiles = []

    source_file_map = {f.short_path: f for f in ctx.files.srcs}

    for info in [dep[ProtoCompileInfo] for dep in ctx.attr.deps]:
        # List[String]: names of files that represent the source files.  In a
        # test, these are the file paths of actual source files that are in the
        # workspace (and checked into source control).
        source_files = []

        # List[String]: names of files that represent the generated files.  In a
        # test, these are the outputs files from the proto_compile rule.
        generated_files = []
        for rel, generated_file in info.output_files_by_rel_path.items():
            source_file = source_file_map.get(rel)
            if not source_file:
                fail("could not find matching source file for generated file %s in %r" % (rel, source_file_map.keys()))

            if source_file.short_path == generated_file.short_path:
                fail("source file path must be distinct from generated file path (%s)" % source_file.short_path)

            runfiles.append(source_file)
            runfiles.append(generated_file)
            source_files.append(source_file.short_path)
            generated_files.append(generated_file.short_path)

        config.packageConfigs.append(
            struct(
                targetLabel = str(info.label),
                targetPackage = info.label.package,
                targetWorkspaceRoot = info.label.workspace_root,
                generatedFiles = generated_files,
                sourceFiles = source_files,
            ),
        )

    config_json, script, runfiles = gencopy_action(ctx, config, runfiles)

    return [DefaultInfo(
        files = depset([config_json]),
        runfiles = runfiles,
        executable = script,
    )]

proto_compile_gencopy_test = rule(
    implementation = _proto_compile_gencopy_test_impl,
    attrs = dict(
        gencopy_attrs,
        deps = attr.label_list(
            doc = "The ProtoCompileInfo providers",
            providers = [ProtoCompileInfo],
        ),
        srcs = attr.label_list(
            doc = "The source files",
            allow_files = True,
        ),
    ),
    executable = True,
    test = True,
)
