load(":proto_compile.bzl", "proto_compile")
load(":proto_compile_gencopy.bzl", "proto_compile_gencopy_test", "proto_compile_gencopy_run")


def proto_compiled_sources(**kwargs):
    name = kwargs.pop("name")
    srcs = kwargs.pop("srcs", [])
    name_update = name + ".update"
    name_test = name + "_test"

    proto_compile(
        name = name,
        srcs = srcs,
        **kwargs,
    )

    proto_compile_gencopy_test(
        name = name_test,
        srcs = srcs,
        deps = [name],
        mode = "check",
        update_target_label_name = name_update,
    )

    proto_compile_gencopy_run(
        name = name_update,
        deps = [name],
        mode = "update",
        update_target_label_name = name_update,
    )
