load("@com_google_protobuf//:protobuf.bzl", "py_proto_library")

# Hack - providers indexing is by int, but I have not idea how to get the actual
# provider object here.
ProtoInfoProvider = 0

def _get_output_filename(src, suffix):
    basename = src.basename
    if basename.endswith(".proto"):
        basename = basename[:-5]
    return basename + "py"

def _py_proto_compile(ctx):
    protoc = ctx.executable.prototool
    #prototool = ctx.executable.prototool

    deps = [dep.proto for dep in ctx.attr.deps]
    sources = []
    outputs = []

    for dep in deps:
        print("dep: %r" % dep)
        for e in dep.transitive_imports:
            print("import: %r" % e)
        for e in dep.transitive_sources:
            print("source: %r" % e)
            sources.append(e)
            filename = _get_output_filename(e, "py")
            output = ctx.actions.declare_file("%s/%s" % (ctx.label.name, filename))
            outputs.append(output)
        for e in dep.transitive_proto_path:
            print("proto_path: %r" % e)
        for e in dep.transitive_descriptor_sets:
            print("descriptor_set: %r" % e)
    
    prototool_json = struct(
        protoc = struct(
            version = ctx.attr.proto_version,
        ),
        generate = struct(
            plugins = [
                struct(
                    name = "python",
                    type = "python",
                ),
            ],
        ),
    )

    ctx.actions.write(ctx.outputs.prototool_json, prototool_json.to_json())
    
    ctx.actions.run(
        executable = protoc,
        arguments = ["generate"],
        inputs = sources + [ctx.outputs.prototool_json],
        outputs = outputs,
        env = {
            "HOME": ctx.outputs.prototool_json.dirname,
        }
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
        "proto_version": attr.string(
            default = "3.6.1",
        ),        
        "prototool": attr.label(
            default = "@//:prototool",
            cfg = "host",
            executable = True,
        ),
        "protoc": attr.label(
            default = "@com_google_protobuf//:protoc",
            cfg = "host",
            executable = True,
        ),
    },
    outputs = {
        "prototool_json": "%{name}/prototool.json",
    },
)

# def py_proto_library(**kwargs):
#     compile_name = kwargs["name"] + ".pb"
#     py_proto_compile(
#         name = compile_name,
#         **kwargs
#     )