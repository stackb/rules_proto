load("//:deps.bzl", 
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
