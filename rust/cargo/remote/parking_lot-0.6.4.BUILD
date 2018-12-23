"""
cargo-raze crate build file.

DO NOT EDIT! Replaced on runs of cargo-raze
"""

package(default_visibility = ["//visibility:public"])

licenses([
    "notice",  # "Apache-2.0,MIT"
])

load(
    "@io_bazel_rules_rust//rust:rust.bzl",
    "rust_binary",
    "rust_library",
    "rust_test",
)

rust_library(
    name = "parking_lot",
    srcs = glob(["**/*.rs"]),
    crate_features = [
        "default",
        "lock_api",
        "owning_ref",
    ],
    crate_root = "src/lib.rs",
    crate_type = "lib",
    rustc_flags = [
        "--cap-lints allow",
        "--target=x86_64-unknown-linux-gnu",
    ],
    deps = [
        "@raze__lock_api__0_1_3//:lock_api",
        "@raze__parking_lot_core__0_3_1//:parking_lot_core",
    ],
)
