load("//rules/internal:common.bzl", "ProtoCompileInfo")
load(
    "@build_bazel_rules_nodejs//:providers.bzl",
    "JSModuleInfo",
    "LinkablePackageInfo",
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
    direct_outputs = [package_file, index_file]
    transitive_depsets = []

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

    has_output_dirs = False

    for dep in ctx.attr.deps:
        if not ProtoCompileInfo in dep:
            print("unknown provider: %r" % dep)
            pass

        info = dep[ProtoCompileInfo]
        for [genrootdir, genfiles] in info.output_files.items():
            for src in genfiles:
                relname = get_source_relname(genrootdir, src)
                base, ext = paths.split_extension(src.basename)
                index_lines.append("  '%s': '%s'" % (base, relname))
                dst = ctx.actions.declare_file(relname, sibling = package_file)
                direct_outputs.append(dst)
                copy_file(ctx, src, dst)

        for dir in info.output_dirs.to_list():
            has_output_dirs = True
            direct_outputs.append(dir)
            ctx.actions.run(
                inputs = [dir],
                outputs = [index_file],
                executable = ctx.executable._tool,
                arguments = [dir.path, index_file.path],
                mnemonic = "CreateIndexJs",
                progress_message = "Creating index.js file",
            )

    for dep in ctx.attr.js_deps:
        if not JSModuleInfo in dep:
            print("skipping unknown provider %r" % dep)
            pass
        transitive_depsets.append(dep[JSModuleInfo].sources)

    index_lines.append("};")

    ctx.actions.write(package_file, package.to_json())

    if not has_output_dirs:
        ctx.actions.write(index_file, "\n".join(index_lines))

    return [
        LinkablePackageInfo(
            package_name = package.name,
            path = ctx.label.package,
            files = depset(direct = direct_outputs, transitive = transitive_depsets),
        ),
        DefaultInfo(
            files = depset(direct = direct_outputs),
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
        "js_deps": attr.label_list(
            doc = "List of rules that provide rules_nodejs providers",
        ),
        "_tool": attr.label(
            doc = "The js_module generating tool",
            default = ":proto_compile_js_module",
            executable = True,
            cfg = "exec",
        ),
    },
    outputs = {
        "package": "%{name}.package.json",
    },
)
