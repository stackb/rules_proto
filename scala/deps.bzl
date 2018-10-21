load("//:deps.bzl", 
    "com_google_protobuf",
    "io_bazel_rules_go",
    "io_bazel_rules_scala",
    "com_github_scalapb_scalapb",
)

def scala_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)
    io_bazel_rules_go(**kwargs)
    io_bazel_rules_scala(**kwargs)
    com_github_scalapb_scalapb(**kwargs)

def scala_grpc_compile(**kwargs):
    scala_proto_compile(**kwargs)

def scala_proto_library(**kwargs):
    scala_proto_compile(**kwargs)

def scala_grpc_library(**kwargs):
    scala_grpc_compile(**kwargs)
    scala_proto_library(**kwargs)
