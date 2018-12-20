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

# Unsupported target "bench-poll" with type "example" omitted

rust_library(
    name = "tokio_reactor",
    crate_root = "src/lib.rs",
    crate_type = "lib",
    srcs = glob(["**/*.rs"]),
    deps = [
        "@raze__crossbeam_utils__0_5_0//:crossbeam_utils",
        "@raze__futures__0_1_24//:futures",
        "@raze__lazy_static__1_1_0//:lazy_static",
        "@raze__log__0_4_5//:log",
        "@raze__mio__0_6_16//:mio",
        "@raze__num_cpus__1_8_0//:num_cpus",
        "@raze__parking_lot__0_6_4//:parking_lot",
        "@raze__slab__0_4_1//:slab",
        "@raze__tokio_executor__0_1_4//:tokio_executor",
        "@raze__tokio_io__0_1_8//:tokio_io",
    ],
    rustc_flags = [
        "--cap-lints allow",
        "--target=x86_64-unknown-linux-gnu",
    ],
    crate_features = [
    ],
)
