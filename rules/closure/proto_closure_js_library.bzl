"proto_closure_js_library.bzl provides a closure_js_library for proto files."

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def proto_closure_js_library(**kwargs):
    suppress = kwargs.pop("suppress", [])
    suppress += [
        "JSC_LATE_PROVIDE_ERROR",
        "JSC_UNDEFINED_VARIABLE",
        "JSC_IMPLICITLY_NULLABLE_JSDOC",
        "JSC_STRICT_INEXISTENT_PROPERTY",
        "JSC_POSSIBLE_INEXISTENT_PROPERTY",
        "JSC_UNRECOGNIZED_TYPE_ERROR",
        "JSC_DEPRECATED_PROP_REASON",
        "JSC_MISSING_REQUIRE_TYPE_IN_PROVIDES_FILE",
    ]

    deps = kwargs.pop("deps", [])
    deps.append("@io_bazel_rules_closure//closure/protobuf:jspb")

    closure_js_library(
        deps = deps,
        suppress = suppress,
        **kwargs
    )
