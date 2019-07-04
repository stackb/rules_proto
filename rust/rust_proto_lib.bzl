load("//:common.bzl", "ProtoCompileInfo")

RustProtoLibInfo = provider(fields = {
    "name": "rule name",
    "lib": "lib.rs file",
})

def _strip_extension(f):
    return f.basename[:-len(f.extension) - 1]

def _rust_proto_lib_impl(ctx):
    """Generate a lib.rs file for the crates."""
    compilation = ctx.attr.compilation[ProtoCompileInfo]
    lib_rs = ctx.actions.declare_file("%s/lib.rs" % compilation.label.name)

    # Add externs
    content = ["extern crate protobuf;"]
    if ctx.attr.grpc:
        content.append("extern crate grpcio;")
        content.append("extern crate futures;")
    content.append("") # Newline

    # List each output
    srcs = [f for files in compilation.output_files.values() for f in files]
    for f in srcs:
        content.append("pub mod %s;" % _strip_extension(f))
        content.append("pub use %s::*;" % _strip_extension(f))

    # Write file
    ctx.actions.write(
        lib_rs,
        "\n".join(content) + "\n",
        False,
    )

    return [RustProtoLibInfo(
        name = ctx.label.name,
        lib = lib_rs,
    ), DefaultInfo(
        files = depset([lib_rs]),
    )]

rust_proto_lib = rule(
    implementation = _rust_proto_lib_impl,
    attrs = {
        "compilation": attr.label(
            providers = [ProtoCompileInfo],
            mandatory = True,
        ),
        "grpc": attr.bool(
            mandatory = True,
        ),
    },
    output_to_genfiles = True,
)
