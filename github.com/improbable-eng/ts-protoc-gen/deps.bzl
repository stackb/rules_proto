load(
    "//:deps.bzl",
    "bazel_gazelle",
    "bazel_skylib",
    "build_bazel_rules_nodejs",
    "build_bazel_rules_typescript",
    "com_google_protobuf",
    "io_bazel_rules_go",
    "io_bazel_rules_webtesting",
    "org_pubref_rules_node",
    "ts_protoc_gen",
)
load(
    "//closure:deps.bzl",
    "io_bazel_rules_closure",
)

def ts_proto_compile(**kwargs):
    bazel_gazelle(**kwargs)
    bazel_skylib(**kwargs)
    org_pubref_rules_node(**kwargs)
    build_bazel_rules_nodejs(**kwargs)
    build_bazel_rules_typescript(**kwargs)
    com_google_protobuf(**kwargs)
    io_bazel_rules_closure(**kwargs)
    io_bazel_rules_go(**kwargs)
    io_bazel_rules_go(**kwargs)
    io_bazel_rules_webtesting(**kwargs)
    ts_protoc_gen(**kwargs)

def ts_grpc_compile(**kwargs):
    ts_proto_compile(**kwargs)

def ts_proto_library(**kwargs):
    ts_proto_compile(**kwargs)

def ts_grpc_library(**kwargs):
    ts_grpc_compile(**kwargs)
    ts_proto_library(**kwargs)
