"""
cargo-raze crate workspace functions

DO NOT EDIT! Replaced on runs of cargo-raze
"""
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "new_git_repository")

def _new_http_archive(name, **kwargs):
    if not native.existing_rule(name):
        http_archive(name=name, **kwargs)

def _new_git_repository(name, **kwargs):
    if not native.existing_rule(name):
        new_git_repository(name=name, **kwargs)

def raze_fetch_remote_crates():

    _new_http_archive(
        name = "raze__byteorder__1_3_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/byteorder/byteorder-1.3.2.crate",
        type = "tar.gz",
        sha256 = "a7c3dd8985a7111efc5c80b44e23ecdd8c007de8ade3b96595387e812b957cf5",
        strip_prefix = "byteorder-1.3.2",
        build_file = Label("//rust/raze/remote:byteorder-1.3.2.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__bytes__0_4_12",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/bytes/bytes-0.4.12.crate",
        type = "tar.gz",
        sha256 = "206fdffcfa2df7cbe15601ef46c813fce0965eb3286db6b56c583b814b51c81c",
        strip_prefix = "bytes-0.4.12",
        build_file = Label("//rust/raze/remote:bytes-0.4.12.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__cc__1_0_37",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/cc/cc-1.0.37.crate",
        type = "tar.gz",
        sha256 = "39f75544d7bbaf57560d2168f28fd649ff9c76153874db88bdbdfd839b1a7e7d",
        strip_prefix = "cc-1.0.37",
        build_file = Label("//rust/raze/remote:cc-1.0.37.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__cfg_if__0_1_9",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/cfg-if/cfg-if-0.1.9.crate",
        type = "tar.gz",
        sha256 = "b486ce3ccf7ffd79fdeb678eac06a9e6c09fc88d33836340becb8fffe87c5e33",
        strip_prefix = "cfg-if-0.1.9",
        build_file = Label("//rust/raze/remote:cfg-if-0.1.9.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__cmake__0_1_40",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/cmake/cmake-0.1.40.crate",
        type = "tar.gz",
        sha256 = "2ca4386c8954b76a8415b63959337d940d724b336cabd3afe189c2b51a7e1ff0",
        strip_prefix = "cmake-0.1.40",
        build_file = Label("//rust/raze/remote:cmake-0.1.40.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__futures__0_1_28",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/futures/futures-0.1.28.crate",
        type = "tar.gz",
        sha256 = "45dc39533a6cae6da2b56da48edae506bb767ec07370f86f70fc062e9d435869",
        strip_prefix = "futures-0.1.28",
        build_file = Label("//rust/raze/remote:futures-0.1.28.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__grpcio__0_4_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/grpcio/grpcio-0.4.4.crate",
        type = "tar.gz",
        sha256 = "c02fb3c9c44615973814c838f75d7898695d4d4b97a3e8cf52e9ccca30664b6f",
        strip_prefix = "grpcio-0.4.4",
        build_file = Label("//rust/raze/remote:grpcio-0.4.4.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__grpcio_compiler__0_4_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/grpcio-compiler/grpcio-compiler-0.4.3.crate",
        type = "tar.gz",
        sha256 = "373a14f0f994d4c235770f4bb5558be00626844db130a82a70142b8fc5996fc3",
        strip_prefix = "grpcio-compiler-0.4.3",
        build_file = Label("//rust/raze/remote:grpcio-compiler-0.4.3.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__grpcio_sys__0_4_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/grpcio-sys/grpcio-sys-0.4.4.crate",
        type = "tar.gz",
        sha256 = "9d8d3b6d1a70b9dcb2545d1aff5b2c74652cb635f6ab6426be8fd201e9566b7e",
        strip_prefix = "grpcio-sys-0.4.4",
        build_file = Label("//rust/raze/remote:grpcio-sys-0.4.4.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__iovec__0_1_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/iovec/iovec-0.1.2.crate",
        type = "tar.gz",
        sha256 = "dbe6e417e7d0975db6512b90796e8ce223145ac4e33c377e4a42882a0e88bb08",
        strip_prefix = "iovec-0.1.2",
        build_file = Label("//rust/raze/remote:iovec-0.1.2.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__itoa__0_4_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/itoa/itoa-0.4.4.crate",
        type = "tar.gz",
        sha256 = "501266b7edd0174f8530248f87f99c88fbe60ca4ef3dd486835b8d8d53136f7f",
        strip_prefix = "itoa-0.4.4",
        build_file = Label("//rust/raze/remote:itoa-0.4.4.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__libc__0_2_59",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/libc/libc-0.2.59.crate",
        type = "tar.gz",
        sha256 = "3262021842bf00fe07dbd6cf34ff25c99d7a7ebef8deea84db72be3ea3bb0aff",
        strip_prefix = "libc-0.2.59",
        build_file = Label("//rust/raze/remote:libc-0.2.59.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__log__0_4_6",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/log/log-0.4.6.crate",
        type = "tar.gz",
        sha256 = "c84ec4b527950aa83a329754b01dbe3f58361d1c5efacd1f6d68c494d08a17c6",
        strip_prefix = "log-0.4.6",
        build_file = Label("//rust/raze/remote:log-0.4.6.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__pkg_config__0_3_14",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/pkg-config/pkg-config-0.3.14.crate",
        type = "tar.gz",
        sha256 = "676e8eb2b1b4c9043511a9b7bea0915320d7e502b0a079fb03f9635a5252b18c",
        strip_prefix = "pkg-config-0.3.14",
        build_file = Label("//rust/raze/remote:pkg-config-0.3.14.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__protobuf__2_7_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/protobuf/protobuf-2.7.0.crate",
        type = "tar.gz",
        sha256 = "5f00e4a3cb64ecfeac2c0a73c74c68ae3439d7a6bead3870be56ad5dd2620a6f",
        strip_prefix = "protobuf-2.7.0",
        build_file = Label("//rust/raze/remote:protobuf-2.7.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__protobuf_codegen__2_7_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/protobuf-codegen/protobuf-codegen-2.7.0.crate",
        type = "tar.gz",
        sha256 = "d2c6e555166cdb646306f599da020e01548e9f4d6ec2fd39802c6db2347cbd3e",
        strip_prefix = "protobuf-codegen-2.7.0",
        build_file = Label("//rust/raze/remote:protobuf-codegen-2.7.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__ryu__1_0_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/ryu/ryu-1.0.0.crate",
        type = "tar.gz",
        sha256 = "c92464b447c0ee8c4fb3824ecc8383b81717b9f1e74ba2e72540aef7b9f82997",
        strip_prefix = "ryu-1.0.0",
        build_file = Label("//rust/raze/remote:ryu-1.0.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__serde__1_0_94",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/serde/serde-1.0.94.crate",
        type = "tar.gz",
        sha256 = "076a696fdea89c19d3baed462576b8f6d663064414b5c793642da8dfeb99475b",
        strip_prefix = "serde-1.0.94",
        build_file = Label("//rust/raze/remote:serde-1.0.94.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__serde_json__1_0_40",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/serde_json/serde_json-1.0.40.crate",
        type = "tar.gz",
        sha256 = "051c49229f282f7c6f3813f8286cc1e3323e8051823fce42c7ea80fe13521704",
        strip_prefix = "serde_json-1.0.40",
        build_file = Label("//rust/raze/remote:serde_json-1.0.40.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__winapi__0_2_8",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi/winapi-0.2.8.crate",
        type = "tar.gz",
        sha256 = "167dc9d6949a9b857f3451275e911c3f44255842c1f7a76f33c55103a909087a",
        strip_prefix = "winapi-0.2.8",
        build_file = Label("//rust/raze/remote:winapi-0.2.8.BUILD.bazel")
    )

