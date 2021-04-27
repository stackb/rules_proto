load(":proto_compile.bzl", "proto_compile")

def proto_compiled_sources(**kwargs):
    name = kwargs.pop("name")
    name_update = name + ".update"
    proto_compile(
        name = name,
        **kwargs)
