load(
    "//:deps.bzl",
    "io_bazel_rules_closure"
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def closure_proto_compile(**kwargs):
    protobuf(**kwargs)

def closure_proto_library(**kwargs):
    closure_proto_compile(**kwargs)
    io_bazel_rules_closure(**kwargs)
