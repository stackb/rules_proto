ProtoDependencyInfo = provider(fields = {
    "name": "The proto dependency name (should correspond to the workspace name",
    "repository_rule": "The name of the repository rule that instantiates this dependency",
    "urls": "The urls string list",
    "sha256": "The sha256 attribute for http_archive",
    "strip_prefix": "The strip_prefix attribute for http_archive",
})

def proto_dependency_info_to_struct(info):
    return struct(
        name = info.name,
        repository_rule = info.repository_rule,
        urls = info.urls,
        sha256 = info.sha256,
        strip_prefix = info.strip_prefix,
    )

def _proto_dependency_impl(ctx):
    return [
        ProtoDependencyInfo(
            name = ctx.attr.name,
            repository_rule = ctx.attr.repository_rule,
            urls = ctx.attr.urls,
            sha256 = ctx.attr.sha256,
            strip_prefix = ctx.attr.strip_prefix,
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
