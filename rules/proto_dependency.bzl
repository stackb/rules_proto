"""proto_dependency.bzl 
"""

load("//rules:providers.bzl", "ProtoDependencyInfo")

def _proto_dependency_impl(ctx):
    return [
        ProtoDependencyInfo(
            buildFile = ctx.attr.build_file,
            buildFileContent = ctx.attr.build_file_content,
            buildFileProtoMode = ctx.attr.build_file_proto_mode,
            deps = [dep[ProtoDependencyInfo] for dep in ctx.attr.deps],
            importpath = ctx.attr.importpath,
            label = str(ctx.label),
            name = ctx.attr.name,
            path = ctx.attr.path,
            repositoryRule = ctx.attr.repository_rule,
            sha256 = ctx.attr.sha256,
            stripPrefix = ctx.attr.strip_prefix,
            sum = ctx.attr.sum,
            urls = ctx.attr.urls,
            version = ctx.attr.version,
            packageJson = ctx.attr.package_json,
            packageLockJson = ctx.attr.package_lock_json,
            yarnLock = ctx.attr.yarn_lock,
            workspaceSnippet = ctx.attr.workspace_snippet,
            frozenLockfile = ctx.attr.frozen_lockfile,
            symlinkNodeModules = ctx.attr.symlink_node_modules,
        ),
    ]

proto_dependency = rule(
    implementation = _proto_dependency_impl,
    attrs = {
        "build_file": attr.string(
            doc = "The build_file attribute for http_archive",
        ),
        "build_file_content": attr.string(
            doc = "The build_file_content attribute for http_archive",
        ),
        "build_file_proto_mode": attr.string(
            doc = "The build_file_proto_mode attribute for go_repository",
        ),
        "deps": attr.label_list(
            doc = "Additional transitive dependencies",
            providers = [ProtoDependencyInfo],
        ),
        "importpath": attr.string(
            doc = "The importpath attribute for go_repository",
        ),
        "path": attr.string(
            doc = "The path attribute for local_repository",
        ),
        "package_json": attr.string(
            doc = "The package_json attribute for npm_install",
        ),
        "package_lock_json": attr.string(
            doc = "The package_lock_json attribute for npm_install",
        ),
        "frozen_lockfile": attr.bool(
            doc = "The frozen_lockfile attribute for yarn_install",
        ),
        "repository_rule": attr.string(
            doc = "The repository rule that instantiates this dependency",
            values = ["http_archive", "http_file", "bind", "go_repository", "local_repository", "phony", "npm_install", "yarn_install"],
        ),
        "sha256": attr.string(
            doc = "The sha256 attribute for http_archive",
        ),
        "strip_prefix": attr.string(
            doc = "The strip_prefix attribute for http_archive",
        ),
        "symlink_node_modules": attr.bool(
            doc = "The symlink_node_modules attribute for npm_install",
        ),
        "sum": attr.string(
            doc = "The sum attribute for go_repository",
        ),
        "workspace_snippet": attr.string(
            doc = "The starlark code snippet for the WORKSPACE needed when using this dependency",
        ),
        "version": attr.string(
            doc = "The version attribute for go_repository",
        ),
        "urls": attr.string_list(
            doc = "The strip_prefix attribute for http_archive",
        ),
        "yarn_lock": attr.string(
            doc = "The yarn_lock attribute for yarn_install",
        ),
    },
)
