load(
    "//:deps.bzl",
    "com_github_scalapb_scalapb",
    "com_google_protobuf",
    "com_thesamet_scalapb_scalapb_json4s",
    "io_bazel_rules_go",
    "io_bazel_rules_scala",
    "org_json4s_json4s_ast_2_12",
    "org_json4s_json4s_jackson_2_12",
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

    # This one actually only needed for routeguide example
    com_thesamet_scalapb_scalapb_json4s(**kwargs)
    org_json4s_json4s_jackson_2_12(**kwargs)
    org_json4s_json4s_ast_2_12(**kwargs)
