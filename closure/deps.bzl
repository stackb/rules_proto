load("//:deps.bzl", 
    "com_google_protobuf", 
    "io_bazel_rules_closure",
)

def closure_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)

def closure_proto_library(**kwargs):
    closure_proto_compile(**kwargs)
    io_bazel_rules_closure(**kwargs)
