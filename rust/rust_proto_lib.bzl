load("//:compile.bzl", "ProtoCompileInfo")

RustProtoLibInfo = provider(fields = {
    "name": "rule name",
    "lib": "lib.rs file",
})

def _basename(f):
    return f.basename[:-len(f.extension) - 1]

def _rust_proto_lib_impl(ctx):
  """Generate a lib.rs file for the crates."""
  compilation = ctx.attr.compilation[ProtoCompileInfo]
  deps = ctx.attr.deps
  grpc = True
  srcs = compilation.files
  lib_rs = ctx.actions.declare_file("%s/lib.rs" % compilation.label.name)

  content = ["extern crate protobuf;"]
  if grpc:
    content.append("extern crate grpc;")
    content.append("extern crate tls_api;")
  # for dep in deps:
  #   content.append("extern crate %s;" % dep.label.name)
  #   content.append("pub use %s::*;" % dep.label.name)
  for f in srcs:
    content.append("pub mod %s;" % _basename(f))
    content.append("pub use %s::*;" % _basename(f))

  ctx.actions.write(
      lib_rs,
      "\n".join(content),
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
    "deps": attr.label_list(
      # providers = [""],
    ),
  },
  output_to_genfiles = True,
)
