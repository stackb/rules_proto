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

# Unsupported target "clock" with type "test" omitted
# Unsupported target "deadline" with type "test" omitted
# Unsupported target "delay" with type "test" omitted
# Unsupported target "hammer" with type "test" omitted
# Unsupported target "interval" with type "test" omitted
# Unsupported target "queue" with type "test" omitted
# Unsupported target "timeout" with type "test" omitted

rust_library(
    name = "tokio_timer",
    crate_root = "src/lib.rs",
    crate_type = "lib",
    srcs = glob(["**/*.rs"]),
    deps = [
        "@raze__crossbeam_utils__0_5_0//:crossbeam_utils",
        "@raze__futures__0_1_24//:futures",
        "@raze__slab__0_4_1//:slab",
        "@raze__tokio_executor__0_1_4//:tokio_executor",
    ],
    rustc_flags = [
        "--cap-lints allow",
        "--target=x86_64-unknown-linux-gnu",
    ],
    crate_features = [
    ],
)

