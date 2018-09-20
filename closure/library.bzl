load("//closure:compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def closure_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
    )

    closure_js_library(
        name = name,
        srcs = [name_pb],
        deps = ["@io_bazel_rules_closure//closure/protobuf:jspb"],
        visibility = visibility,
        internal_descriptors = [name_pb + "/descriptor.source.bin"],
        lenient = True,
    )
