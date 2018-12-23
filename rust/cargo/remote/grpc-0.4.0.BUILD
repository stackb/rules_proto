"""
cargo-raze crate build file.

DO NOT EDIT! Replaced on runs of cargo-raze
"""

package(default_visibility = ["//visibility:public"])

licenses([
    "notice",  # "MIT,Apache-2.0"
])

load(
    "@io_bazel_rules_rust//rust:rust.bzl",
    "rust_binary",
    "rust_library",
    "rust_test",
)

# Unsupported target "client" with type "test" omitted

rust_library(
    name = "grpc",
    srcs = glob(["**/*.rs"]),
    crate_features = [
    ],
    crate_root = "src/lib.rs",
    crate_type = "lib",
    rustc_flags = [
        "--cap-lints allow",
        "--target=x86_64-unknown-linux-gnu",
    ],
    deps = [
        "@raze__base64__0_9_3//:base64",
        "@raze__bytes__0_4_10//:bytes",
        "@raze__futures__0_1_24//:futures",
        "@raze__futures_cpupool__0_1_8//:futures_cpupool",
        "@raze__httpbis__0_6_1//:httpbis",
        "@raze__log__0_4_5//:log",
        "@raze__protobuf__1_6_0//:protobuf",
        "@raze__tls_api__0_1_20//:tls_api",
        "@raze__tls_api_stub__0_1_20//:tls_api_stub",
        "@raze__tokio_core__0_1_17//:tokio_core",
        "@raze__tokio_io__0_1_8//:tokio_io",
        "@raze__tokio_tls_api__0_1_20//:tokio_tls_api",
    ],
)

# Unsupported target "server" with type "test" omitted
# Unsupported target "simple" with type "test" omitted
