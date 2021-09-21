"grpc_closure_js_library.bzl provides a closure_js_library for grpc files."

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def grpc_closure_js_library(**kwargs):
    suppress = kwargs.pop("suppress", [])
    suppress += [
        "JSC_MISSING_REQUIRE_TYPE",
        "reportUnknownTypes",  # TODO: external/com_github_stackb_grpc_js/js/grpc/transport/fetch/observer.js:181
    ]

    deps = kwargs.pop("deps", [])
    deps += [
        "@io_bazel_rules_closure//closure/library/promise",
        "@com_github_stackb_grpc_js//js/grpc/stream/observer:call",
        "@com_github_stackb_grpc_js//js/grpc",
        "@com_github_stackb_grpc_js//js/grpc:api",
        "@com_github_stackb_grpc_js//js/grpc:options",
    ]

    closure_js_library(
        deps = deps,
        suppress = suppress,
        **kwargs
    )
