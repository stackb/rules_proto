"""proto_compiled_sources provides the proto_compiled_sources macro which contains three rules.

1.  The main target of the macro is the 'proto_compile' rule (e.g. //proto:foo_py_proto_compile).  
    This rule runs 'protoc' and generates output files.

2. A 'proto_compile_gencopy_run' rule (e.g. bazel run //proto:foo_py_proto_compile.update) is used 
   to copy the generated sources back into the workspace package.

3. A 'proto_compile_gencopy_test' rule (e.g. bazel test //proto:foo_py_proto_compile_test) is used 
   to compare the generated sources with those in workspace package, preventing drift between the 
   source of truth (proto file) and the workspace copy (presumably under source control).  
   This is be used to force developers to update the generated files when making proto file changes. 
"""

load(":proto_compile.bzl", "proto_compile")
load(":proto_compile_gencopy.bzl", "proto_compile_gencopy_run", "proto_compile_gencopy_test")

def proto_compiled_sources(**kwargs):
    """proto_compiled_sources macro.

    Args:
        **kwargs: the kwargs dict for the 'proto_compile' rule.
    """
    name = kwargs.pop("name")
    srcs = kwargs.pop("srcs", [])
    tags = kwargs.pop("tags", [])
    protoc = kwargs.pop("protoc", None)
    name_update = name + ".update"
    name_test = name + "_test"

    proto_compile(
        name = name,
        srcs = srcs,
        protoc = protoc,
        **kwargs
    )

    proto_compile_gencopy_run(
        name = name_update,
        deps = [name],
        mode = "update",
        update_target_label_name = name_update,
        tags = tags,
    )

    proto_compile_gencopy_test(
        name = name_test,
        srcs = srcs,
        deps = [name],
        mode = "check",
        update_target_label_name = name_update,
        tags = tags,
    )
