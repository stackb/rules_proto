"proto_rust_library.bzl provides a rust_library for proto files."

load("@rules_rust//rust:defs.bzl", "rust_library")

def proto_rust_library(**kwargs):
    rust_library(**kwargs)
