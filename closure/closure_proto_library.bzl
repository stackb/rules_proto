load("//closure:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def closure_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    closure_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create closure library
    closure_js_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = ["@io_bazel_rules_closure//closure/protobuf:jspb"],
        visibility = kwargs.get("visibility"),
        suppress = [
            "JSC_LATE_PROVIDE_ERROR",
            "JSC_UNDEFINED_VARIABLE",
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            "JSC_STRICT_INEXISTENT_PROPERTY",
            "JSC_POSSIBLE_INEXISTENT_PROPERTY",
            "JSC_UNRECOGNIZED_TYPE_ERROR",
        ],
    )
