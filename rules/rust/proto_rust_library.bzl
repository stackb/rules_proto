"proto_rust_library.bzl provides a rust_library for proto files."

load("@rules_rust//rust:defs.bzl", "rust_library")
load(":rust_proto_lib.bzl", "rust_proto_lib")

def proto_rust_library(**kwargs):
    name = kwargs.get("name", "")
    name_lib = name + "_lib"
    compilations = kwargs.pop("compilations")

    srcs = kwargs.pop("srcs", [])
    srcs.append(name_lib)
    kwargs.setdefault("srcs", srcs)

    rust_proto_lib(
        name = name_lib,
        compilations = compilations,
        externs = ["protobuf"],
    )

    rust_library(**kwargs)
