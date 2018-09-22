"""
cargo-raze crate build file.

DO NOT EDIT! Replaced on runs of cargo-raze
"""
package(default_visibility = ["//visibility:public"])

licenses([
  "notice", # "MIT"
])

load(
    "@io_bazel_rules_rust//rust:rust.bzl",
    "rust_library",
    "rust_binary",
    "rust_test",
)


rust_library(
    name = "tokio_udp",
    crate_root = "src/lib.rs",
    crate_type = "lib",
    srcs = glob(["**/*.rs"]),
    deps = [
        "@raze__bytes__0_4_10//:bytes",
        "@raze__futures__0_1_24//:futures",
        "@raze__log__0_4_5//:log",
        "@raze__mio__0_6_16//:mio",
        "@raze__tokio_codec__0_1_0//:tokio_codec",
        "@raze__tokio_io__0_1_8//:tokio_io",
        "@raze__tokio_reactor__0_1_5//:tokio_reactor",
    ],
    rustc_flags = [
        "--cap-lints allow",
        "--target=x86_64-unknown-linux-gnu",
    ],
    crate_features = [
    ],
)

# Unsupported target "udp" with type "test" omitted
