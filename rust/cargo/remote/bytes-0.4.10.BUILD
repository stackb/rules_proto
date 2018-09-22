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

# Unsupported target "bytes" with type "bench" omitted

rust_library(
    name = "bytes",
    crate_root = "src/lib.rs",
    crate_type = "lib",
    srcs = glob(["**/*.rs"]),
    deps = [
        "@raze__byteorder__1_2_6//:byteorder",
        "@raze__iovec__0_1_2//:iovec",
    ],
    rustc_flags = [
        "--cap-lints allow",
        "--target=x86_64-unknown-linux-gnu",
    ],
    crate_features = [
    ],
)

# Unsupported target "test_buf" with type "test" omitted
# Unsupported target "test_buf_mut" with type "test" omitted
# Unsupported target "test_bytes" with type "test" omitted
# Unsupported target "test_chain" with type "test" omitted
# Unsupported target "test_debug" with type "test" omitted
# Unsupported target "test_from_buf" with type "test" omitted
# Unsupported target "test_iter" with type "test" omitted
# Unsupported target "test_serde" with type "test" omitted
# Unsupported target "test_take" with type "test" omitted
