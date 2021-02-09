ProtoDependencyInfo = provider(fields = {
    "build_file": "The build_file of this dependency",
    "deps": "The deps of this dependency",
    "label": "The proto dependency label",
    "name": "The proto dependency name (should correspond to the workspace name",
    "repositoryRule": "The name of the repository rule that instantiates this dependency",
    "sha256": "The sha256 attribute for http_archive",
    "stripPrefix": "The strip_prefix attribute for http_archive",
    "urls": "The urls string list",
    "workspaceSnippet": "The workspaceSnippet string list",
})

def proto_dependency_info_to_struct(info):
    return struct(
        buildFile = info.build_file,
        label = str(info.label),
        name = info.name,
        repositoryRule = info.repositoryRule,
        sha256 = info.sha256,
        stripPrefix = info.stripPrefix,
        urls = info.urls,
        workspaceSnippet = info.workspaceSnippet,
    )

def _proto_dependency_impl(ctx):
    return [
        ProtoDependencyInfo(
            build_file = ctx.attr.build_file,
            deps = depset(direct = [dep[ProtoDependencyInfo] for dep in ctx.attr.deps], transitive = [dep[ProtoDependencyInfo].deps for dep in ctx.attr.deps]),
            label = ctx.label,
            name = ctx.attr.name,
            repositoryRule = ctx.attr.repository_rule,
            sha256 = ctx.attr.sha256,
            stripPrefix = ctx.attr.strip_prefix,
            urls = ctx.attr.urls,
            workspaceSnippet = ctx.attr.workspace_snippet,
        ),
    ]

proto_dependency = rule(
    implementation = _proto_dependency_impl,
    attrs = {
        "build_file": attr.string(
            doc = "The build_file attribute for http_archive",
        ),
        "workspace_snippet": attr.string(
            doc = "The starlark code snippet for the WORKSPACE needed when using this dependency",
        ),
        "deps": attr.label_list(
            doc = "Additional transitive dependencies",
            providers = [ProtoDependencyInfo],
        ),
        "repository_rule": attr.string(
            doc = "The repository rule that instantiates this dependency",
            values = ["http_archive", "http_file", "bind", "go_repository"],
        ),
        "sha256": attr.string(
            doc = "The sha256 attribute for http_archive",
        ),
        "strip_prefix": attr.string(
            doc = "The strip_prefix attribute for http_archive",
        ),
        "urls": attr.string_list(
            doc = "The strip_prefix attribute for http_archive",
        ),
    },
)
