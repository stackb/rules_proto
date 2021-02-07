ProtoDependencyInfo = provider(fields = {
    "name": "The proto dependency name (should correspond to the workspace name",
    "label": "The proto dependency label",
    "repositoryRule": "The name of the repository rule that instantiates this dependency",
    "urls": "The urls string list",
    "sha256": "The sha256 attribute for http_archive",
    "stripPrefix": "The strip_prefix attribute for http_archive",
    "deps": "The deps of this dependency",
})

def proto_dependency_info_to_struct(info):
    return struct(
        name = info.name,
        label = str(info.label),
        repositoryRule = info.repositoryRule,
        urls = info.urls,
        sha256 = info.sha256,
        stripPrefix = info.stripPrefix,
    )

def _proto_dependency_impl(ctx):
    return [
        ProtoDependencyInfo(
            name = ctx.attr.name,
            label = ctx.label,
            repositoryRule = ctx.attr.repository_rule,
            urls = ctx.attr.urls,
            sha256 = ctx.attr.sha256,
            stripPrefix = ctx.attr.strip_prefix,
            deps = depset(direct = [dep[ProtoDependencyInfo] for dep in ctx.attr.deps], transitive = [dep[ProtoDependencyInfo].deps for dep in ctx.attr.deps]),
        ),
    ]

proto_dependency = rule(
    implementation = _proto_dependency_impl,
    attrs = {
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
        "deps": attr.label_list(
            doc = "Additional transitive dependencies",
            providers = [ProtoDependencyInfo],
        ),
    },
)
