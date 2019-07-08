load(
    "//:deps.bzl",
    "io_bazel_rules_go",
)

def gogo_deps(**kwargs):
    # Same as rules_go as rules_go is already loading gogo protobuf
    io_bazel_rules_go(**kwargs)
    native.register_toolchains(str(Label("//protobuf:protoc_toolchain")))

def gogo_proto_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogo_grpc_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogo_proto_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogo_grpc_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogotypes_proto_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogotypes_grpc_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogotypes_proto_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogotypes_grpc_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogoslick_proto_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogoslick_grpc_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogoslick_proto_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogoslick_grpc_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogofast_proto_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogofast_grpc_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogofast_proto_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogofast_grpc_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogofaster_proto_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogofaster_grpc_compile(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogofaster_proto_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)

def gogofaster_grpc_library(**kwargs): # Kept for backwards compatibility
    gogo_deps(**kwargs)
