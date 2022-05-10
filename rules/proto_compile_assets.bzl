"proto_compile_assets.bzl provides the files copy rule."

load(":proto_compile_gencopy.bzl", "proto_compile_gencopy_run")

def proto_compile_assets(**kwargs):
    """proto_compile_assets copies generated files to the source tree

    Args:
        **kwargs: the kwargs macro dict.  Should have 'name' and 'deps' attributes.
        Deps must provide ProtoCompileInfo.
    """
    name = kwargs.pop("name")
    deps = kwargs.pop("deps", [])

    proto_compile_gencopy_run(
        name = name,
        deps = deps,
        mode = "update",
        update_target_label_name = name,
    )
