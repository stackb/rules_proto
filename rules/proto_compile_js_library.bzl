load("//rules/internal:common.bzl", "ProtoCompileInfo")
load(
    "@build_bazel_rules_nodejs//:providers.bzl",
    "LinkablePackageInfo",
    # "direct_sources": "Depset of direct JavaScript files and sourcemaps",
    # "sources": "Depset of direct and transitive JavaScript files and sourcemaps",
)
load("@bazel_skylib//lib:paths.bzl", "paths")

def copy_file(ctx, src, dst):
    ctx.actions.run_shell(
        tools = [src],
        outputs = [dst],
        command = "cp -f \"$1\" \"$2\"",
        arguments = [src.path, dst.path],
        mnemonic = "CopyFile",
        progress_message = "Copying files",
        use_default_shell_env = True,
    )

def get_source_relname(genrootdir, file):
    before, sep, after = file.path.partition(genrootdir)
    return after[1:]

def package_name_from_label(label):
    return "@" + label.name.lower()

def _proto_compile_js_library_impl(ctx):
    package_file = ctx.outputs.package
    index_file = ctx.actions.declare_file(ctx.label.name + ".js", sibling = package_file)
    outputs = [package_file, index_file]

    index_lines = [
        "// Generated file - DO NOT EDIT",
        "",
        "module.exports = {",
    ]

    package = struct(
        name = package_name_from_label(ctx.label),
        main = index_file.basename,
        files = [],
    )

    for dep in ctx.attr.deps:
        if ProtoCompileInfo in dep:
            info = dep[ProtoCompileInfo]
            for [genrootdir, genfiles] in info.output_files.items():
                for src in genfiles:
                    relname = get_source_relname(genrootdir, src)
                    base, ext = paths.split_extension(src.basename)
                    index_lines.append("  '%s': '%s'" % (base, relname))
                    dst = ctx.actions.declare_file(relname, sibling = package_file)
                    outputs.append(dst)
                    copy_file(ctx, src, dst)
        else:
            print("unknown provider: %r" % dep)

    index_lines.append("};")

    ctx.actions.write(package_file, package.to_json())
    ctx.actions.write(index_file, "\n".join(index_lines))

    sources = depset(direct = outputs)

    return [
        LinkablePackageInfo(
            package_name = package.name,
            path = ctx.label.package,
            files = depset(direct = outputs),
        ),
        # JSModuleInfo(
        #     sources = sources,
        # ),
        DefaultInfo(
            files = sources,
        ),
    ]

proto_compile_js_library = rule(
    implementation = _proto_compile_js_library_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "List of proto_compile rules that are providing *.js proto outputs",
            mandatory = True,
            providers = [ProtoCompileInfo],
        ),
    },
    outputs = {
        "package": "%{name}package.json",
    },
)
