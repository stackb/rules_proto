load("//:deps", "grpc_deps")

def cpp_grpc_compile_deps():
    grpc_deps()

def cpp_grpc_library_deps():
    cpp_grpc_compile_deps()