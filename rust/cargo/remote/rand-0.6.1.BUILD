"""
cargo-raze crate build file.

DO NOT EDIT! Replaced on runs of cargo-raze
"""

package(default_visibility = [
    # Public for visibility by "@raze__crate__version//" targets.
    #
    # Prefer access through "//rust/cargo", which limits external
    # visibility to explicit Cargo.toml dependencies.
    "//visibility:public",
])

licenses([
    "notice",  # "MIT,Apache-2.0"
])

load(
    "@io_bazel_rules_rust//rust:rust.bzl",
    "rust_binary",
    "rust_library",
    "rust_test",
)

# Unsupported target "build-script-build" with type "custom-build" omitted
# Unsupported target "distributions" with type "bench" omitted
# Unsupported target "generators" with type "bench" omitted
# Unsupported target "misc" with type "bench" omitted
# Unsupported target "monte-carlo" with type "example" omitted
# Unsupported target "monty-hall" with type "example" omitted

rust_library(
    name = "rand",
    srcs = glob(["**/*.rs"]),
    crate_features = [
        "alloc",
        "cloudabi",
        "default",
        "fuchsia-zircon",
        "libc",
        "rand_core",
        "std",
        "winapi",
    ],
    crate_root = "src/lib.rs",
    crate_type = "lib",
    rustc_flags = [
        "--cap-lints=allow",
    ],
    version = "0.6.1",
    deps = [
        "@raze__libc__0_2_45//:libc",
        "@raze__rand_chacha__0_1_0//:rand_chacha",
        "@raze__rand_core__0_3_0//:rand_core",
        "@raze__rand_hc__0_1_0//:rand_hc",
        "@raze__rand_isaac__0_1_1//:rand_isaac",
        "@raze__rand_pcg__0_1_1//:rand_pcg",
        "@raze__rand_xorshift__0_1_0//:rand_xorshift",
    ],
)

# Unsupported target "seq" with type "bench" omitted
# Unsupported target "uniformity" with type "test" omitted
