load("//rules:providers.bzl", "ProtoCompileInfo")

def _strip_extension(f):
    return f.basename[:-len(f.extension) - 1]

def _rust_proto_lib_impl(ctx):
    """Generate a lib.rs file for the crates."""
    compilations = [info for info in ctx.attr.compilation[ProtoCompileInfo]]

    content = []
    for extern in ctx.attr.externs:
        content.append("extern crate {};".format(extern))

    for compilation in compilations:
        for f in compilation.outputs:
            content.append("pub mod %s;" % _strip_extension(f))
            content.append("pub use %s::*;" % _strip_extension(f))

    # Write file
    lib_rs = ctx.actions.declare_file("lib.rs")
    ctx.actions.write(
        lib_rs,
        "\n".join(content) + "\n",
        False,
    )

    return [DefaultInfo(
        files = depset([lib_rs]),
    )]

rust_proto_lib = rule(
    implementation = _rust_proto_lib_impl,
    attrs = {
        "compilation": attr.label(
            providers = [ProtoCompileInfo],
            mandatory = True,
        ),
        "externs": attr.string_list(
            mandatory = True,
        ),
    },
)
