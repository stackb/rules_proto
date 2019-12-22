load("//closure:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def closure_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.pop("verbose", 0)
    closure_deps = kwargs.pop("closure_deps", [])
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        verbose = verbose,
        visibility = visibility,
        transitive = kwargs.pop("transitive", True),
        transitivity = kwargs.pop("transitivity", {}),
    )

    closure_js_library(
        name = name,
        srcs = [name_pb],
        deps = ["@io_bazel_rules_closure//closure/protobuf:jspb"] + closure_deps,
        visibility = visibility,
        internal_descriptors = [name_pb + "/descriptor.source.bin"],
        suppress = [
            "JSC_LATE_PROVIDE_ERROR",
            "JSC_UNDEFINED_VARIABLE",
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            "JSC_STRICT_INEXISTENT_PROPERTY",
            "JSC_POSSIBLE_INEXISTENT_PROPERTY",
            "JSC_UNRECOGNIZED_TYPE_ERROR",
        ],
    )
