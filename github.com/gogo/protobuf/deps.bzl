load(
    "//:deps.bzl",
    "io_bazel_rules_go",
)

# Same as rules_go as rules_go is already loading gogo protobuf

def gogo_proto_compile(**kwargs):
    io_bazel_rules_go(**kwargs)

def gogo_grpc_compile(**kwargs):
    gogo_proto_compile(**kwargs)

def gogo_proto_library(**kwargs):
    gogo_proto_compile(**kwargs)

def gogo_grpc_library(**kwargs):
    gogo_grpc_compile(**kwargs)
    gogo_proto_library(**kwargs)

def gogotypes_proto_compile(**kwargs):
    gogo_proto_compile(**kwargs)

def gogotypes_grpc_compile(**kwargs):
    gogo_grpc_compile(**kwargs)

def gogotypes_proto_library(**kwargs):
    gogo_proto_library(**kwargs)

def gogotypes_grpc_library(**kwargs):
    gogo_grpc_library(**kwargs)

def gogoslick_proto_compile(**kwargs):
    gogo_proto_compile(**kwargs)

def gogoslick_grpc_compile(**kwargs):
    gogo_grpc_compile(**kwargs)

def gogoslick_proto_library(**kwargs):
    gogo_proto_library(**kwargs)

def gogoslick_grpc_library(**kwargs):
    gogo_grpc_library(**kwargs)

def gogofast_proto_compile(**kwargs):
    gogo_proto_compile(**kwargs)

def gogofast_grpc_compile(**kwargs):
    gogo_grpc_compile(**kwargs)

def gogofast_proto_library(**kwargs):
    gogo_proto_library(**kwargs)

def gogofast_grpc_library(**kwargs):
    gogo_grpc_library(**kwargs)

def gogofaster_proto_compile(**kwargs):
    gogo_proto_compile(**kwargs)

def gogofaster_grpc_compile(**kwargs):
    gogo_grpc_compile(**kwargs)

def gogofaster_proto_library(**kwargs):
    gogo_proto_library(**kwargs)

def gogofaster_grpc_library(**kwargs):
    gogo_grpc_library(**kwargs)
