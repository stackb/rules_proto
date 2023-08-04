"""starlark plugin definitions"""

def _configure_protoc_gen_java(ctx):
    """_configure_protoc_gen_java prepares the PluginConfiguration for a fictitious protoc java plugin.

    Args:
        ctx (protoc.PluginContext): The context object.
    Returns:
        config (PluginConfiguration): The configured PluginConfiguration object.
    """

    srcjar = ctx.proto_library.base_name + ".srcjar"
    if ctx.rel:
        srcjar = "/".join([ctx.rel, srcjar])

    config = protoc.PluginConfiguration(
        label = "@build_stack_rules_proto//plugin/builtin:java",
        outputs = [srcjar],
        out = srcjar,
        options = ctx.plugin_config.options,
    )

    return config

protoc.Plugin(
    name = "protoc-gen-java",
    configure = _configure_protoc_gen_java,
)
