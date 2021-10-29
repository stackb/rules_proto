"grpc_closure_js_library.bzl provides a closure_js_library for grpc files."

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def grpc_closure_js_library(**kwargs):
    """grpc_closure_js_library is a thin wrapper over closure_js_library.

    Args:
        **kwargs: keyword arguments passed to the closure_js_library.  Additional supressions and dependencies are added.
    """
    suppress = kwargs.pop("suppress", [])
    suppress.append(
        "JSC_MISSING_REQUIRE_TYPE",
    )

    deps = kwargs.pop("deps", [])
    deps += [
        "@com_google_javascript_closure_library//closure/goog/promise",
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
