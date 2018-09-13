load("@com_google_protobuf//:protobuf.bzl", "py_proto_library")

# Hack - providers indexing is by int, but I have not idea how to get the actual
# provider object here.
ProtoInfoProvider = 0

def _get_output_filename(src, suffix):
    basename = src.basename
    if basename.endswith(".proto"):
        basename = basename[:-6]
    return basename + suffix

def _py_proto_compile(ctx):
    protoc = ctx.executable.protoc
    descriptor = ctx.outputs.descriptor

    deps = [dep.proto for dep in ctx.attr.deps]
    protos = []
    outputs = []
    args = []

    for dep in deps:
        print("dep: %r" % dep)
        # for e in dep.transitive_imports:
        #     print("import: %r" % e)
        #     args += ["--proto_path", path]
        for src in dep.transitive_sources:
            print("source: %r" % src)
            proto = ctx.actions.declare_file(src.short_path, sibling = descriptor)
            protos.append(proto)
            ctx.actions.run_shell(
                mnemonic = "CopyProto",
                inputs = [src],
                outputs = [proto],
                command = "find . && cp %s %s" % (src.path, proto.path),
            )
            gen = ctx.actions.declare_file(_get_output_filename(src, "_pb2.py"), sibling = proto)
            outputs.append(gen)
        for e in dep.transitive_proto_path:
            print("proto_path: %r" % e)
        for e in dep.transitive_descriptor_sets:
            print("descriptor_set: %r" % e)

    args += ["--descriptor_set_out=%s" % descriptor.path]
    args += ["--proto_path=%s" % descriptor.dirname]        
    args += ["--python_out=%s" % descriptor.dirname]        
    args += [proto.path for proto in protos]

    ctx.actions.run_shell(
        command = "find . && " +  " ".join([protoc.path] + args) + " && find .",
        inputs = [protoc] + protos,
        outputs = outputs + [descriptor],
    )

    return [DefaultInfo(
        files = depset(outputs),
    )]

py_proto_compile = rule(
    implementation = _py_proto_compile,
    attrs = {
        "deps": attr.label_list(
            mandatory = True,
            providers = ["proto"],
        ),
        "protoc": attr.label(
            default = "@com_google_protobuf//:protoc",
            cfg = "host",
            executable = True,
        ),
    },
    outputs = {
        "descriptor": "%{name}/descriptor.bin",
    }
)

# def py_proto_library(**kwargs):
#     compile_name = kwargs["name"] + ".pb"
#     py_proto_compile(
#         name = compile_name,
#         **kwargs
#     )