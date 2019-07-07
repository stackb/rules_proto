def _rustc_resolve_impl(ctx):
    # Get resolved rust toolchain
    rust_toolchain = ctx.toolchains["@io_bazel_rules_rust//rust:toolchain"]

    # Copy rustc to output executable
    rustc = ctx.actions.declare_file("rustc/bin/rustc")
    ctx.actions.run_shell(
        inputs = [rust_toolchain.rustc],
        outputs = [rustc],
        command = "cp '{}' '{}'".format(rust_toolchain.rustc.path, rustc.path),
    )

    # Copy minimal lib files required to get rustc to work
    lib_files = []
    for f in rust_toolchain.rustc_lib.files.to_list():
        new_f = ctx.actions.declare_file("rustc/lib/{}".format(f.basename))
        ctx.actions.run_shell(
            inputs = [f],
            outputs = [new_f],
            command = "cp '{}' '{}'".format(f.path, new_f.path),
        )
        lib_files.append(new_f)

    # Return default info with an executable pointing to our local rustc
    return [
        DefaultInfo(
            runfiles = ctx.runfiles(
                files = lib_files,
                collect_data = True,
            ),
            executable = rustc,
        )
    ]

# Rule that consumes the rules_rust toolchain and exposes the resolved rustc binary
rustc_resolve = rule(
    _rustc_resolve_impl,
    toolchains = [
        "@io_bazel_rules_rust//rust:toolchain",
    ],
)