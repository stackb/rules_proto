"proto_compiled_source_update.bzl provides the proto_compiled_source_update rule."

load(":proto_compile_gencopy.bzl", "proto_compile_gencopy_run")

def proto_compiled_source_update(**kwargs):
    name = kwargs.pop("name")
    deps = kwargs.pop("deps", [])
    name_update = name

    proto_compile_gencopy_run(
        name = name_update,
        deps = deps,
        mode = "update",
        update_target_label_name = name,
    )
