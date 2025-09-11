load("@io_bazel_rules_go//go:def.bzl", "GoArchive")
load("@io_bazel_rules_go//go/private:common.bzl", "GO_TOOLCHAIN")

ProtoGoModulesInfo = provider(
    doc = "info provided from a go_modules rule",
    fields = {
        "direct": "List[GoArchive] deps of this rule",
        "label": "[Label]: the label of this rule",
    },
)

def _is_proto_dep(go_archive_data):
    for src in go_archive_data.srcs:
        if not src.basename.endswith(".pb.go"):
            return False
    return True

def _proto_go_modules_impl(ctx):
    # index the GoArchive objects by importpath.
    direct = {}

    for dep in ctx.attr.deps:
        go_archive = dep[GoArchive]
        direct[go_archive.data.importpath] = go_archive
    for module in ctx.attr.modules:
        for imp, go_archive in module[ProtoGoModulesInfo].direct.items():
            if direct.get(imp) == None:
                direct[imp] = go_archive

    # collect proto_modules info for debug purposes
    available_imports = [
        go_archive.data
        for go_archive in direct.values()
        if _is_proto_dep(go_archive.data)
    ]

    # collect the final list of deps we want to vendor sources for
    want = {}
    if len(ctx.attr.imports) == 0:
        want = {imp: d.data for imp, d in direct.items()}
    else:
        for imp in ctx.attr.imports:
            go_archive = direct.get(imp)
            if not go_archive:
                fail("no known GoArchive for %s.  Please ensure .deps or .modules provides the corresponding import." % imp)
            if want.get(imp) != None:
                want[imp] = go_archive.data
            for go_archive_data in go_archive.transitive.to_list():
                if want.get(go_archive_data.importpath) == None:
                    # try from the direct pool first
                    direct_go_archive = direct.get(go_archive_data.importpath)
                    if direct_go_archive and _is_proto_dep(direct_go_archive.data):
                        want[go_archive_data.importpath] = direct_go_archive.data

    # srcs will be a list of .go files that will be included in runfiles such
    # that they will be available for copy operations.
    srcs = []

    # get the go toolchain for (1) the go tool and (2) the version that we stamp
    # into generated go.mod files
    toolchain = ctx.toolchains[GO_TOOLCHAIN]

    # lines is a list of shell script commands to be written to the executable
    # output script.
    lines = [
        "set -euox pipefail",
        "cwd=$PWD",
        """gobin="$PWD/%s" """ % toolchain.sdk.go.short_path,
        """cd $BUILD_WORKING_DIRECTORY """,
    ]

    for go_archive_data in want.values():
        dstdir = "./%s/%s" % (ctx.attr.srcroot, go_archive_data.importpath)

        lines.append("")
        lines.append("# module=" + str(go_archive_data.importpath))
        lines.append("mkdir -p %s" % dstdir)
        lines.append("echo 'module %s' > %s/go.mod" % (go_archive_data.importpath, dstdir))
        lines.append("""echo "go %s" >> %s/go.mod""" % (toolchain.sdk.version, dstdir))

        for src in go_archive_data.srcs:
            srcs.append(src)
            dst = "%s/%s" % (dstdir, src.basename)
            lines.append("""cp -f "${cwd}/%s" %s""" % (src.short_path, dst))
            lines.append("""echo '# %s' >> %s/go.mod""" % (src.short_path, dstdir))  # to record where the src was copied from

        lines.append(""""${gobin}" mod edit -replace %s=%s""" % (go_archive_data.importpath, dstdir))

    lines.append("")
    lines.append("")

    lines.extend(sorted([
        "# %s provided by %s" % (d.importpath, d.label)
        for d in available_imports
    ]))

    ctx.actions.write(
        output = ctx.outputs.executable,
        content = "\n".join(lines),
    )

    return [
        DefaultInfo(
            executable = ctx.outputs.executable,
            files = depset([ctx.outputs.executable]),
            runfiles = ctx.runfiles(files = srcs + [toolchain.sdk.go]),
        ),
        ProtoGoModulesInfo(
            label = ctx.label,
            direct = direct,
        ),
    ]

proto_go_modules = rule(
    implementation = _proto_go_modules_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "list of libraries that provide GoArchive (proto_go_library, go_library, ...)",
            providers = [GoArchive],
        ),
        "modules": attr.label_list(
            doc = "list of labels that provide ProtoGoModulesInfo (go_modules)",
            providers = [ProtoGoModulesInfo],
        ),
        "imports": attr.string_list(
            doc = "list of go importpaths that represent top-level proto imports that are desired.  The transitive set of proto import dependencies will be computed from this set",
        ),
        "srcroot": attr.string(
            default = "local",
        ),
    },
    provides = [DefaultInfo, ProtoGoModulesInfo],
    toolchains = [GO_TOOLCHAIN],
    executable = True,
)
