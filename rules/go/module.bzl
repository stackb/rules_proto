load("@io_bazel_rules_go//go:def.bzl", "GoArchive")

GoModuleInfo = provider(
    doc = "info about a go_module",
    fields = {
        "output_file": "[File] the go.mod output file",
        "srcs": "List[File] the source files of the module (go_library)",
        "label": "[Label]: the label of the originating rule",
        "importpath": "[String] the module importpath",
    },
)

GoModulesInfo = provider(
    doc = "info provided from a go_modules rule",
    fields = {
        "deps": "List[GoArchive] deps of this rule",
        "modules": "List[GoModuleInfo] module deps of this rule",
        "label": "[Label]: the label of this rule",
    },
)

def _go_module_impl(ctx):
    output_file = ctx.actions.declare_file("go.mod")

    inputs = ctx.files.srcs

    ctx.actions.run_shell(
        mnemonic = "WriteGoModule",
        inputs = inputs,
        outputs = [output_file],
        command = """echo 'module {importpath}\ngo {version}' > {output_file}""".format(
            importpath = ctx.attr.importpath,
            version = ctx.attr.go_version,
            output_file = output_file.path,
        ),
    )

    return [
        DefaultInfo(files = depset([output_file])),
        GoModuleInfo(
            output_file = output_file,
            label = ctx.label,
            importpath = ctx.attr.importpath,
            srcs = ctx.files.srcs,
        ),
    ]

go_module = rule(
    implementation = _go_module_impl,
    attrs = {
        "importpath": attr.string(
            doc = "go importpath",
            mandatory = True,
        ),
        "go_version": attr.string(
            doc = "go version",
            default = "1.23.0",
        ),
        "srcs": attr.label_list(
            doc = "generated sources included in the module",
            allow_files = True,
        ),
    },
    provides = [
        DefaultInfo,
        GoModuleInfo,
    ],
)

def _is_proto_dep(go_archive_data):
    for src in go_archive_data.srcs:
        if src.basename.endswith(".pb.go"):
            return True
    return False

def _go_modules_impl(ctx):
    # comprehend a list of GoArchive providers
    deps = [
        dep[GoArchive]
        for dep in ctx.attr.deps
    ]
    modules = [
        dep[GoModulesInfo]
        for dep in ctx.attr.modules
    ]
    for module in modules:
        deps.extend(module.deps)

    # index the GoArchiveData objects by importpath.  Collect the direct and
    # transitive deps into different containers
    direct = {}
    transitive = {}
    for dep in deps:
        direct[dep.data.importpath] = dep
        for d in dep.direct:
            transitive[d.data.importpath] = d.data
        for d in dep.transitive.to_list():
            transitive[d.importpath] = d

    # collect the final list of deps we want to vendor sources for
    want = []
    if len(ctx.attr.imports) == 0:
        want = direct.values()
    else:
        for imp in ctx.attr.imports:
            dep = direct.get(imp)
            if not dep:
                fail("no known GoArchive for %s.  Please ensure .deps or .modules provides the corresponding import." % imp)
            want.append(dep)
            for d in dep.direct:
                if _is_proto_dep(d.data):
                    want.append(d)

    # srcs will be a list of .go files that will be included in runfiles such
    # that they will be available for copy operations.
    srcs = []

    # lines is a list of shell script commands to be written to the executable
    # output script.
    lines = [
        "set -euox pipefail",
        """cd $BUILD_WORKING_DIRECTORY""",
    ]

    for dep in want:
        dstdir = "./%s/%s" % (ctx.attr.srcroot, dep.data.importpath)

        lines.append("")
        lines.append("# module=" + str(dep.data.importpath))
        lines.append("mkdir -p %s" % dstdir)
        lines.append("echo 'module %s' > %s/go.mod" % (dep.data.importpath, dstdir))

        if ctx.attr.go_version == "go.mod":
            lines.append("""echo "$(grep '^go' < %s)" >> %s/go.mod""" % (ctx.attr.go_version, dstdir))
        else:
            lines.append("""echo "go %s" >> %s/go.mod""" % (ctx.attr.go_version, dstdir))

        for src in dep.data.srcs:
            srcs.append(src)
            dst = "%s/%s" % (dstdir, src.basename)
            lines.append("cp -f %s %s" % (src.path, dst))

        lines.append("go mod edit -replace %s=%s" % (dep.data.importpath, dstdir))

    lines.append("")
    lines.append("")

    ctx.actions.write(
        output = ctx.outputs.executable,
        content = "\n".join(lines),
    )

    return [
        DefaultInfo(
            executable = ctx.outputs.executable,
            files = depset([ctx.outputs.executable]),
            runfiles = ctx.runfiles(files = srcs),
        ),
        GoModulesInfo(
            deps = deps,
            modules = modules,
        ),
    ]

go_modules = rule(
    implementation = _go_modules_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "list of libraries that provide GoArchive (proto_go_library, go_library, ...)",
            providers = [GoArchive],
        ),
        "modules": attr.label_list(
            doc = "list of labels that provide GoModulesInfo (go_modules)",
            providers = [GoModulesInfo],
        ),
        "imports": attr.string_list(
            doc = "list of go importpaths that represent top-level proto imports that are desired.  The transitive set of proto import dependencies will be computed from this set",
        ),
        "go_version": attr.string(
            default = "1.23.0",
        ),
        "srcroot": attr.string(
            default = "local",
        ),
    },
    provides = [DefaultInfo, GoModulesInfo],
    executable = True,
)

# def _go_vendor_impl(ctx):
#     deps = [
#         dep[GoArchive]
#         for dep in ctx.attr.deps
#     ]

#     src_map = {}
#     dir_map = {}
#     src_files = []
#     for dep in deps:
#         for src in dep.data.srcs:
#             dst_dir = "generated/%s" % dep.data.importpath
#             dir_map[dst_dir] = dep.data.importpath
#             src_files.append(src)
#             src_map[src.short_path] = "%s/%s" % (dst_dir, src.basename)

#     # # modules = [
#     # #     dep[GoModuleInfo]
#     # #     for dep in ctx.attr.modules
#     # # ]
#     # # module_files = []
#     # # module_srcs = []

#     # for module in modules:
#     #     module_files.append(module.output_file)
#     #     module_srcs += module.srcs

#     ctx.actions.write(
#         output = ctx.outputs.executable,
#         content = "\n".join([
#             'mkdir -p "$BUILD_WORKING_DIRECTORY/{dir}"'.format(dir = dir)
#             for dir in dir_map.keys()
#         ]) + "\n\n" + "\n".join([
#             'echo "module {importpath}" > "$BUILD_WORKING_DIRECTORY/{dir}/go.mod"'.format(
#                 dir = dir,
#                 importpath = importpath,
#             )
#             for dir, importpath in dir_map.items()
#         ]) + "\n\n" + "\n".join([
#             'cp -f %s "$BUILD_WORKING_DIRECTORY/%s"' % (src, dst)
#             for src, dst in src_map.items()
#         ]) + """
# {protogomodulereplace} \\
# --file "$BUILD_WORKING_DIRECTORY/go.mod" \\
# """.format(
#             protogomodulereplace = ctx.executable._protogomodulereplace.short_path,
#         ) + "\n".join([
#             "--replace='%s=>./%s' \\\\" % (dep.data.importpath, dep.data.file.dirname)
#             for dep in deps
#         ]),
#         is_executable = True,
#     )

#     return [
#         DefaultInfo(
#             executable = ctx.outputs.executable,
#             files = depset([ctx.outputs.executable]),
#             runfiles = ctx.runfiles(files = [
#                 ctx.outputs.executable,
#                 ctx.executable._protogomodulereplace,
#             ] + src_files),
#         ),
#     ]

# go_vendor = rule(
#     implementation = _go_vendor_impl,
#     attrs = {
#         "modules": attr.label_list(
#             doc = "go_module rule dependencies",
#             providers = [GoModuleInfo],
#         ),
#         "deps": attr.label_list(
#             doc = "go_module rule dependencies",
#             providers = [GoArchive],
#         ),
#         "_protogomodulereplace": attr.label(
#             doc = "the protogomodulereplace tool",
#             default = "//cmd/protogomodulereplace",
#             executable = True,
#             cfg = "exec",
#         ),
#     },
#     provides = [DefaultInfo],
#     executable = True,
# )
