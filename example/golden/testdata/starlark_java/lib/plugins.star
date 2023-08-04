"""starlark plugin definitions"""

def _configure_java(ctx):
    """_configure_java prepares the PluginConfiguration for the builtin protoc java plugin.

    Args:
        ctx (protoc.PluginContext): The context object.
    Returns:
        config (PluginConfiguration): The configured PluginConfiguration object.
    """

    # replicating the following golang in starlark:
    #
    # srcjar := path.Join(ctx.Rel, ctx.ProtoLibrary.BaseName()+".srcjar")
    # return &protoc.PluginConfiguration{
    # 	Label:   label.New("build_stack_rules_proto", "plugin/builtin", "java"),
    # 	Outputs: []string{srcjar},
    # 	Out:     srcjar,
    # 	Options: ctx.PluginConfig.GetOptions(),
    # }
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
    name = "java",
    configure = _configure_java,
)
