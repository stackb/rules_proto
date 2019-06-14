load("//github.com/grpc/grpc-web:closure_grpc_compile.bzl", "closure_grpc_compile")
load("//closure:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def closure_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    name_pb_grpc = kwargs.get("name") + "_pb_grpc"
    closure_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )
    closure_grpc_compile(
        name = name_pb_grpc,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create closure library
    closure_js_library(
        name = kwargs.get("name"),
        srcs = [name_pb, name_pb_grpc],
        deps = [
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:abstractclientbase",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:clientreadablestream",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:grpcwebclientbase",
            "@com_github_grpc_grpc_web//javascript/net/grpc/web:error",
            "@io_bazel_rules_closure//closure/library",
            "@io_bazel_rules_closure//closure/protobuf:jspb",
        ],
        suppress = [
            "JSC_LATE_PROVIDE_ERROR",
            "JSC_UNDEFINED_VARIABLE",
            "JSC_IMPLICITLY_NULLABLE_JSDOC",
            "JSC_STRICT_INEXISTENT_PROPERTY",
            "JSC_POSSIBLE_INEXISTENT_PROPERTY",
            "JSC_UNRECOGNIZED_TYPE_ERROR",
            "JSC_UNUSED_PRIVATE_PROPERTY",
            "JSC_EXTRA_REQUIRE_WARNING",
            "JSC_INVALID_INTERFACE_MEMBER_DECLARATION",
        ],
        visibility = kwargs.get("visibility"),
    )
