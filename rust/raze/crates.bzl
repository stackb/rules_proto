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
        name = "raze__argon2rs__0_2_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/argon2rs/argon2rs-0.2.5.crate",
        type = "tar.gz",
        sha256 = "3f67b0b6a86dae6e67ff4ca2b6201396074996379fba2b92ff649126f37cb392",
        strip_prefix = "argon2rs-0.2.5",
        build_file = Label("//rust/raze/remote:argon2rs-0.2.5.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__arrayvec__0_4_10",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/arrayvec/arrayvec-0.4.10.crate",
        type = "tar.gz",
        sha256 = "92c7fb76bc8826a8b33b4ee5bb07a247a81e76764ab4d55e8f73e3a4d8808c71",
        strip_prefix = "arrayvec-0.4.10",
        build_file = Label("//rust/raze/remote:arrayvec-0.4.10.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__autocfg__0_1_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/autocfg/autocfg-0.1.4.crate",
        type = "tar.gz",
        sha256 = "0e49efa51329a5fd37e7c79db4621af617cd4e3e5bc224939808d076077077bf",
        strip_prefix = "autocfg-0.1.4",
        build_file = Label("//rust/raze/remote:autocfg-0.1.4.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__backtrace__0_3_32",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/backtrace/backtrace-0.3.32.crate",
        type = "tar.gz",
        sha256 = "18b50f5258d1a9ad8396d2d345827875de4261b158124d4c819d9b351454fae5",
        strip_prefix = "backtrace-0.3.32",
        build_file = Label("//rust/raze/remote:backtrace-0.3.32.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__backtrace_sys__0_1_30",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/backtrace-sys/backtrace-sys-0.1.30.crate",
        type = "tar.gz",
        sha256 = "5b3a000b9c543553af61bc01cbfc403b04b5caa9e421033866f2e98061eb3e61",
        strip_prefix = "backtrace-sys-0.1.30",
        build_file = Label("//rust/raze/remote:backtrace-sys-0.1.30.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__bitflags__1_1_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/bitflags/bitflags-1.1.0.crate",
        type = "tar.gz",
        sha256 = "3d155346769a6855b86399e9bc3814ab343cd3d62c7e985113d46a0ec3c281fd",
        strip_prefix = "bitflags-1.1.0",
        build_file = Label("//rust/raze/remote:bitflags-1.1.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__blake2_rfc__0_2_18",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/blake2-rfc/blake2-rfc-0.2.18.crate",
        type = "tar.gz",
        sha256 = "5d6d530bdd2d52966a6d03b7a964add7ae1a288d25214066fd4b600f0f796400",
        strip_prefix = "blake2-rfc-0.2.18",
        build_file = Label("//rust/raze/remote:blake2-rfc-0.2.18.BUILD.bazel")
    )

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
        name = "raze__chrono__0_4_7",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/chrono/chrono-0.4.7.crate",
        type = "tar.gz",
        sha256 = "77d81f58b7301084de3b958691458a53c3f7e0b1d702f77e550b6a88e3a88abe",
        strip_prefix = "chrono-0.4.7",
        build_file = Label("//rust/raze/remote:chrono-0.4.7.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__cloudabi__0_0_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/cloudabi/cloudabi-0.0.3.crate",
        type = "tar.gz",
        sha256 = "ddfc5b9aa5d4507acaf872de71051dfd0e309860e88966e1051e462a077aac4f",
        strip_prefix = "cloudabi-0.0.3",
        build_file = Label("//rust/raze/remote:cloudabi-0.0.3.BUILD.bazel")
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
        name = "raze__constant_time_eq__0_1_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/constant_time_eq/constant_time_eq-0.1.3.crate",
        type = "tar.gz",
        sha256 = "8ff012e225ce166d4422e0e78419d901719760f62ae2b7969ca6b564d1b54a9e",
        strip_prefix = "constant_time_eq-0.1.3",
        build_file = Label("//rust/raze/remote:constant_time_eq-0.1.3.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__crossbeam__0_2_12",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam/crossbeam-0.2.12.crate",
        type = "tar.gz",
        sha256 = "bd66663db5a988098a89599d4857919b3acf7f61402e61365acfd3919857b9be",
        strip_prefix = "crossbeam-0.2.12",
        build_file = Label("//rust/raze/remote:crossbeam-0.2.12.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__crossbeam__0_6_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam/crossbeam-0.6.0.crate",
        type = "tar.gz",
        sha256 = "ad4c7ea749d9fb09e23c5cb17e3b70650860553a0e2744e38446b1803bf7db94",
        strip_prefix = "crossbeam-0.6.0",
        build_file = Label("//rust/raze/remote:crossbeam-0.6.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__crossbeam_channel__0_3_8",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam-channel/crossbeam-channel-0.3.8.crate",
        type = "tar.gz",
        sha256 = "0f0ed1a4de2235cabda8558ff5840bffb97fcb64c97827f354a451307df5f72b",
        strip_prefix = "crossbeam-channel-0.3.8",
        build_file = Label("//rust/raze/remote:crossbeam-channel-0.3.8.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__crossbeam_deque__0_6_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam-deque/crossbeam-deque-0.6.3.crate",
        type = "tar.gz",
        sha256 = "05e44b8cf3e1a625844d1750e1f7820da46044ff6d28f4d43e455ba3e5bb2c13",
        strip_prefix = "crossbeam-deque-0.6.3",
        build_file = Label("//rust/raze/remote:crossbeam-deque-0.6.3.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__crossbeam_epoch__0_7_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam-epoch/crossbeam-epoch-0.7.1.crate",
        type = "tar.gz",
        sha256 = "04c9e3102cc2d69cd681412141b390abd55a362afc1540965dad0ad4d34280b4",
        strip_prefix = "crossbeam-epoch-0.7.1",
        build_file = Label("//rust/raze/remote:crossbeam-epoch-0.7.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__crossbeam_utils__0_6_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam-utils/crossbeam-utils-0.6.5.crate",
        type = "tar.gz",
        sha256 = "f8306fcef4a7b563b76b7dd949ca48f52bc1141aa067d2ea09565f3e2652aa5c",
        strip_prefix = "crossbeam-utils-0.6.5",
        build_file = Label("//rust/raze/remote:crossbeam-utils-0.6.5.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__dirs__1_0_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/dirs/dirs-1.0.5.crate",
        type = "tar.gz",
        sha256 = "3fd78930633bd1c6e35c4b42b1df7b0cbc6bc191146e512bb3bedf243fcc3901",
        strip_prefix = "dirs-1.0.5",
        build_file = Label("//rust/raze/remote:dirs-1.0.5.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__failure__0_1_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/failure/failure-0.1.5.crate",
        type = "tar.gz",
        sha256 = "795bd83d3abeb9220f257e597aa0080a508b27533824adf336529648f6abf7e2",
        strip_prefix = "failure-0.1.5",
        build_file = Label("//rust/raze/remote:failure-0.1.5.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__failure_derive__0_1_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/failure_derive/failure_derive-0.1.5.crate",
        type = "tar.gz",
        sha256 = "ea1063915fd7ef4309e222a5a07cf9c319fb9c7836b1f89b85458672dbb127e1",
        strip_prefix = "failure_derive-0.1.5",
        build_file = Label("//rust/raze/remote:failure_derive-0.1.5.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__fuchsia_cprng__0_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/fuchsia-cprng/fuchsia-cprng-0.1.1.crate",
        type = "tar.gz",
        sha256 = "a06f77d526c1a601b7c4cdd98f54b5eaabffc14d5f2f0296febdc7f357c6d3ba",
        strip_prefix = "fuchsia-cprng-0.1.1",
        build_file = Label("//rust/raze/remote:fuchsia-cprng-0.1.1.BUILD.bazel")
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
        name = "raze__isatty__0_1_9",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/isatty/isatty-0.1.9.crate",
        type = "tar.gz",
        sha256 = "e31a8281fc93ec9693494da65fbf28c0c2aa60a2eaec25dc58e2f31952e95edc",
        strip_prefix = "isatty-0.1.9",
        build_file = Label("//rust/raze/remote:isatty-0.1.9.BUILD.bazel")
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
        name = "raze__lazy_static__1_3_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/lazy_static/lazy_static-1.3.0.crate",
        type = "tar.gz",
        sha256 = "bc5729f27f159ddd61f4df6228e827e86643d4d3e7c32183cb30a1c08f604a14",
        strip_prefix = "lazy_static-1.3.0",
        build_file = Label("//rust/raze/remote:lazy_static-1.3.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__libc__0_2_58",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/libc/libc-0.2.58.crate",
        type = "tar.gz",
        sha256 = "6281b86796ba5e4366000be6e9e18bf35580adf9e63fbe2294aadb587613a319",
        strip_prefix = "libc-0.2.58",
        build_file = Label("//rust/raze/remote:libc-0.2.58.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__lock_api__0_1_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/lock_api/lock_api-0.1.5.crate",
        type = "tar.gz",
        sha256 = "62ebf1391f6acad60e5c8b43706dde4582df75c06698ab44511d15016bc2442c",
        strip_prefix = "lock_api-0.1.5",
        build_file = Label("//rust/raze/remote:lock_api-0.1.5.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__log__0_3_9",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/log/log-0.3.9.crate",
        type = "tar.gz",
        sha256 = "e19e8d5c34a3e0e2223db8e060f9e8264aeeb5c5fc64a4ee9965c062211c024b",
        strip_prefix = "log-0.3.9",
        build_file = Label("//rust/raze/remote:log-0.3.9.BUILD.bazel")
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
        name = "raze__memoffset__0_2_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/memoffset/memoffset-0.2.1.crate",
        type = "tar.gz",
        sha256 = "0f9dc261e2b62d7a622bf416ea3c5245cdd5d9a7fcc428c0d06804dfce1775b3",
        strip_prefix = "memoffset-0.2.1",
        build_file = Label("//rust/raze/remote:memoffset-0.2.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__nodrop__0_1_13",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/nodrop/nodrop-0.1.13.crate",
        type = "tar.gz",
        sha256 = "2f9667ddcc6cc8a43afc9b7917599d7216aa09c463919ea32c59ed6cac8bc945",
        strip_prefix = "nodrop-0.1.13",
        build_file = Label("//rust/raze/remote:nodrop-0.1.13.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__num_integer__0_1_41",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/num-integer/num-integer-0.1.41.crate",
        type = "tar.gz",
        sha256 = "b85e541ef8255f6cf42bbfe4ef361305c6c135d10919ecc26126c4e5ae94bc09",
        strip_prefix = "num-integer-0.1.41",
        build_file = Label("//rust/raze/remote:num-integer-0.1.41.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__num_traits__0_2_8",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/num-traits/num-traits-0.2.8.crate",
        type = "tar.gz",
        sha256 = "6ba9a427cfca2be13aa6f6403b0b7e7368fe982bfa16fccc450ce74c46cd9b32",
        strip_prefix = "num-traits-0.2.8",
        build_file = Label("//rust/raze/remote:num-traits-0.2.8.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__num_cpus__1_10_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/num_cpus/num_cpus-1.10.1.crate",
        type = "tar.gz",
        sha256 = "bcef43580c035376c0705c42792c294b66974abbfd2789b511784023f71f3273",
        strip_prefix = "num_cpus-1.10.1",
        build_file = Label("//rust/raze/remote:num_cpus-1.10.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__owning_ref__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/owning_ref/owning_ref-0.4.0.crate",
        type = "tar.gz",
        sha256 = "49a4b8ea2179e6a2e27411d3bca09ca6dd630821cf6894c6c7c8467a8ee7ef13",
        strip_prefix = "owning_ref-0.4.0",
        build_file = Label("//rust/raze/remote:owning_ref-0.4.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__parking_lot__0_7_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/parking_lot/parking_lot-0.7.1.crate",
        type = "tar.gz",
        sha256 = "ab41b4aed082705d1056416ae4468b6ea99d52599ecf3169b00088d43113e337",
        strip_prefix = "parking_lot-0.7.1",
        build_file = Label("//rust/raze/remote:parking_lot-0.7.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__parking_lot_core__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/parking_lot_core/parking_lot_core-0.4.0.crate",
        type = "tar.gz",
        sha256 = "94c8c7923936b28d546dfd14d4472eaf34c99b14e1c973a32b3e6d4eb04298c9",
        strip_prefix = "parking_lot_core-0.4.0",
        build_file = Label("//rust/raze/remote:parking_lot_core-0.4.0.BUILD.bazel")
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
        name = "raze__proc_macro2__0_4_30",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/proc-macro2/proc-macro2-0.4.30.crate",
        type = "tar.gz",
        sha256 = "cf3d2011ab5c909338f7887f4fc896d35932e29146c12c8d01da6b22a80ba759",
        strip_prefix = "proc-macro2-0.4.30",
        build_file = Label("//rust/raze/remote:proc-macro2-0.4.30.BUILD.bazel")
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
        name = "raze__quote__0_6_12",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/quote/quote-0.6.12.crate",
        type = "tar.gz",
        sha256 = "faf4799c5d274f3868a4aae320a0a182cbd2baee377b378f080e16a23e9d80db",
        strip_prefix = "quote-0.6.12",
        build_file = Label("//rust/raze/remote:quote-0.6.12.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand__0_6_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand/rand-0.6.5.crate",
        type = "tar.gz",
        sha256 = "6d71dacdc3c88c1fde3885a3be3fbab9f35724e6ce99467f7d9c5026132184ca",
        strip_prefix = "rand-0.6.5",
        build_file = Label("//rust/raze/remote:rand-0.6.5.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_chacha__0_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_chacha/rand_chacha-0.1.1.crate",
        type = "tar.gz",
        sha256 = "556d3a1ca6600bfcbab7c7c91ccb085ac7fbbcd70e008a98742e7847f4f7bcef",
        strip_prefix = "rand_chacha-0.1.1",
        build_file = Label("//rust/raze/remote:rand_chacha-0.1.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_core__0_3_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_core/rand_core-0.3.1.crate",
        type = "tar.gz",
        sha256 = "7a6fdeb83b075e8266dcc8762c22776f6877a63111121f5f8c7411e5be7eed4b",
        strip_prefix = "rand_core-0.3.1",
        build_file = Label("//rust/raze/remote:rand_core-0.3.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_core__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_core/rand_core-0.4.0.crate",
        type = "tar.gz",
        sha256 = "d0e7a549d590831370895ab7ba4ea0c1b6b011d106b5ff2da6eee112615e6dc0",
        strip_prefix = "rand_core-0.4.0",
        build_file = Label("//rust/raze/remote:rand_core-0.4.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_hc__0_1_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_hc/rand_hc-0.1.0.crate",
        type = "tar.gz",
        sha256 = "7b40677c7be09ae76218dc623efbf7b18e34bced3f38883af07bb75630a21bc4",
        strip_prefix = "rand_hc-0.1.0",
        build_file = Label("//rust/raze/remote:rand_hc-0.1.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_isaac__0_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_isaac/rand_isaac-0.1.1.crate",
        type = "tar.gz",
        sha256 = "ded997c9d5f13925be2a6fd7e66bf1872597f759fd9dd93513dd7e92e5a5ee08",
        strip_prefix = "rand_isaac-0.1.1",
        build_file = Label("//rust/raze/remote:rand_isaac-0.1.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_jitter__0_1_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_jitter/rand_jitter-0.1.4.crate",
        type = "tar.gz",
        sha256 = "1166d5c91dc97b88d1decc3285bb0a99ed84b05cfd0bc2341bdf2d43fc41e39b",
        strip_prefix = "rand_jitter-0.1.4",
        build_file = Label("//rust/raze/remote:rand_jitter-0.1.4.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_os__0_1_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_os/rand_os-0.1.3.crate",
        type = "tar.gz",
        sha256 = "7b75f676a1e053fc562eafbb47838d67c84801e38fc1ba459e8f180deabd5071",
        strip_prefix = "rand_os-0.1.3",
        build_file = Label("//rust/raze/remote:rand_os-0.1.3.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_pcg__0_1_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_pcg/rand_pcg-0.1.2.crate",
        type = "tar.gz",
        sha256 = "abf9b09b01790cfe0364f52bf32995ea3c39f4d2dd011eac241d2914146d0b44",
        strip_prefix = "rand_pcg-0.1.2",
        build_file = Label("//rust/raze/remote:rand_pcg-0.1.2.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rand_xorshift__0_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_xorshift/rand_xorshift-0.1.1.crate",
        type = "tar.gz",
        sha256 = "cbf7e9e623549b0e21f6e97cf8ecf247c1a8fd2e8a992ae265314300b2455d5c",
        strip_prefix = "rand_xorshift-0.1.1",
        build_file = Label("//rust/raze/remote:rand_xorshift-0.1.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rdrand__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rdrand/rdrand-0.4.0.crate",
        type = "tar.gz",
        sha256 = "678054eb77286b51581ba43620cc911abf02758c91f93f479767aed0f90458b2",
        strip_prefix = "rdrand-0.4.0",
        build_file = Label("//rust/raze/remote:rdrand-0.4.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__redox_syscall__0_1_55",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/redox_syscall/redox_syscall-0.1.55.crate",
        type = "tar.gz",
        sha256 = "cbd093aba707641a1ebf784490f00ceb889b21953774de1d5fba05d1ee1d283c",
        strip_prefix = "redox_syscall-0.1.55",
        build_file = Label("//rust/raze/remote:redox_syscall-0.1.55.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__redox_users__0_3_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/redox_users/redox_users-0.3.0.crate",
        type = "tar.gz",
        sha256 = "3fe5204c3a17e97dde73f285d49be585df59ed84b50a872baf416e73b62c3828",
        strip_prefix = "redox_users-0.3.0",
        build_file = Label("//rust/raze/remote:redox_users-0.3.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rustc_demangle__0_1_15",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rustc-demangle/rustc-demangle-0.1.15.crate",
        type = "tar.gz",
        sha256 = "a7f4dccf6f4891ebcc0c39f9b6eb1a83b9bf5d747cb439ec6fba4f3b977038af",
        strip_prefix = "rustc-demangle-0.1.15",
        build_file = Label("//rust/raze/remote:rustc-demangle-0.1.15.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__rustc_version__0_2_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rustc_version/rustc_version-0.2.3.crate",
        type = "tar.gz",
        sha256 = "138e3e0acb6c9fb258b19b67cb8abd63c00679d2851805ea151465464fe9030a",
        strip_prefix = "rustc_version-0.2.3",
        build_file = Label("//rust/raze/remote:rustc_version-0.2.3.BUILD.bazel")
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
        name = "raze__scoped_threadpool__0_1_9",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/scoped_threadpool/scoped_threadpool-0.1.9.crate",
        type = "tar.gz",
        sha256 = "1d51f5df5af43ab3f1360b429fa5e0152ac5ce8c0bd6485cae490332e96846a8",
        strip_prefix = "scoped_threadpool-0.1.9",
        build_file = Label("//rust/raze/remote:scoped_threadpool-0.1.9.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__scopeguard__0_3_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/scopeguard/scopeguard-0.3.3.crate",
        type = "tar.gz",
        sha256 = "94258f53601af11e6a49f722422f6e3425c52b06245a5cf9bc09908b174f5e27",
        strip_prefix = "scopeguard-0.3.3",
        build_file = Label("//rust/raze/remote:scopeguard-0.3.3.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__semver__0_9_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/semver/semver-0.9.0.crate",
        type = "tar.gz",
        sha256 = "1d7eb9ef2c18661902cc47e535f9bc51b78acd254da71d375c2f6720d9a40403",
        strip_prefix = "semver-0.9.0",
        build_file = Label("//rust/raze/remote:semver-0.9.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__semver_parser__0_7_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/semver-parser/semver-parser-0.7.0.crate",
        type = "tar.gz",
        sha256 = "388a1df253eca08550bef6c72392cfe7c30914bf41df5269b68cbd6ff8f570a3",
        strip_prefix = "semver-parser-0.7.0",
        build_file = Label("//rust/raze/remote:semver-parser-0.7.0.BUILD.bazel")
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
        name = "raze__slog__2_4_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/slog/slog-2.4.1.crate",
        type = "tar.gz",
        sha256 = "1e1a2eec401952cd7b12a84ea120e2d57281329940c3f93c2bf04f462539508e",
        strip_prefix = "slog-2.4.1",
        build_file = Label("//rust/raze/remote:slog-2.4.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__slog_async__2_3_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/slog-async/slog-async-2.3.0.crate",
        type = "tar.gz",
        sha256 = "e544d16c6b230d84c866662fe55e31aacfca6ae71e6fc49ae9a311cb379bfc2f",
        strip_prefix = "slog-async-2.3.0",
        build_file = Label("//rust/raze/remote:slog-async-2.3.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__slog_scope__4_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/slog-scope/slog-scope-4.1.1.crate",
        type = "tar.gz",
        sha256 = "60c04b4726fa04595ccf2c2dad7bcd15474242c4c5e109a8a376e8a2c9b1539a",
        strip_prefix = "slog-scope-4.1.1",
        build_file = Label("//rust/raze/remote:slog-scope-4.1.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__slog_stdlog__3_0_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/slog-stdlog/slog-stdlog-3.0.2.crate",
        type = "tar.gz",
        sha256 = "ac42f8254ae996cc7d640f9410d3b048dcdf8887a10df4d5d4c44966de24c4a8",
        strip_prefix = "slog-stdlog-3.0.2",
        build_file = Label("//rust/raze/remote:slog-stdlog-3.0.2.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__slog_term__2_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/slog-term/slog-term-2.4.0.crate",
        type = "tar.gz",
        sha256 = "5951a808c40f419922ee014c15b6ae1cd34d963538b57d8a4778b9ca3fff1e0b",
        strip_prefix = "slog-term-2.4.0",
        build_file = Label("//rust/raze/remote:slog-term-2.4.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__smallvec__0_6_10",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/smallvec/smallvec-0.6.10.crate",
        type = "tar.gz",
        sha256 = "ab606a9c5e214920bb66c458cd7be8ef094f813f20fe77a54cc7dbfff220d4b7",
        strip_prefix = "smallvec-0.6.10",
        build_file = Label("//rust/raze/remote:smallvec-0.6.10.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__stable_deref_trait__1_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/stable_deref_trait/stable_deref_trait-1.1.1.crate",
        type = "tar.gz",
        sha256 = "dba1a27d3efae4351c8051072d619e3ade2820635c3958d826bfea39d59b54c8",
        strip_prefix = "stable_deref_trait-1.1.1",
        build_file = Label("//rust/raze/remote:stable_deref_trait-1.1.1.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__syn__0_15_39",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/syn/syn-0.15.39.crate",
        type = "tar.gz",
        sha256 = "b4d960b829a55e56db167e861ddb43602c003c7be0bee1d345021703fac2fb7c",
        strip_prefix = "syn-0.15.39",
        build_file = Label("//rust/raze/remote:syn-0.15.39.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__synstructure__0_10_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/synstructure/synstructure-0.10.2.crate",
        type = "tar.gz",
        sha256 = "02353edf96d6e4dc81aea2d8490a7e9db177bf8acb0e951c24940bf866cb313f",
        strip_prefix = "synstructure-0.10.2",
        build_file = Label("//rust/raze/remote:synstructure-0.10.2.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__take_mut__0_2_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/take_mut/take_mut-0.2.2.crate",
        type = "tar.gz",
        sha256 = "f764005d11ee5f36500a149ace24e00e3da98b0158b3e2d53a7495660d3f4d60",
        strip_prefix = "take_mut-0.2.2",
        build_file = Label("//rust/raze/remote:take_mut-0.2.2.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__term__0_5_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/term/term-0.5.2.crate",
        type = "tar.gz",
        sha256 = "edd106a334b7657c10b7c540a0106114feadeb4dc314513e97df481d5d966f42",
        strip_prefix = "term-0.5.2",
        build_file = Label("//rust/raze/remote:term-0.5.2.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__thread_local__0_3_6",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/thread_local/thread_local-0.3.6.crate",
        type = "tar.gz",
        sha256 = "c6b53e329000edc2b34dbe8545fd20e55a333362d0a321909685a19bd28c3f1b",
        strip_prefix = "thread_local-0.3.6",
        build_file = Label("//rust/raze/remote:thread_local-0.3.6.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__time__0_1_42",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/time/time-0.1.42.crate",
        type = "tar.gz",
        sha256 = "db8dcfca086c1143c9270ac42a2bbd8a7ee477b78ac8e45b19abfb0cbede4b6f",
        strip_prefix = "time-0.1.42",
        build_file = Label("//rust/raze/remote:time-0.1.42.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__unicode_xid__0_1_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/unicode-xid/unicode-xid-0.1.0.crate",
        type = "tar.gz",
        sha256 = "fc72304796d0818e357ead4e000d19c9c174ab23dc11093ac919054d20a6a7fc",
        strip_prefix = "unicode-xid-0.1.0",
        build_file = Label("//rust/raze/remote:unicode-xid-0.1.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__winapi__0_2_8",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi/winapi-0.2.8.crate",
        type = "tar.gz",
        sha256 = "167dc9d6949a9b857f3451275e911c3f44255842c1f7a76f33c55103a909087a",
        strip_prefix = "winapi-0.2.8",
        build_file = Label("//rust/raze/remote:winapi-0.2.8.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__winapi__0_3_7",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi/winapi-0.3.7.crate",
        type = "tar.gz",
        sha256 = "f10e386af2b13e47c89e7236a7a14a086791a2b88ebad6df9bf42040195cf770",
        strip_prefix = "winapi-0.3.7",
        build_file = Label("//rust/raze/remote:winapi-0.3.7.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__winapi_i686_pc_windows_gnu__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi-i686-pc-windows-gnu/winapi-i686-pc-windows-gnu-0.4.0.crate",
        type = "tar.gz",
        sha256 = "ac3b87c63620426dd9b991e5ce0329eff545bccbbb34f3be09ff6fb6ab51b7b6",
        strip_prefix = "winapi-i686-pc-windows-gnu-0.4.0",
        build_file = Label("//rust/raze/remote:winapi-i686-pc-windows-gnu-0.4.0.BUILD.bazel")
    )

    _new_http_archive(
        name = "raze__winapi_x86_64_pc_windows_gnu__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi-x86_64-pc-windows-gnu/winapi-x86_64-pc-windows-gnu-0.4.0.crate",
        type = "tar.gz",
        sha256 = "712e227841d057c1ee1cd2fb22fa7e5a5461ae8e48fa2ca79ec42cfc1931183f",
        strip_prefix = "winapi-x86_64-pc-windows-gnu-0.4.0",
        build_file = Label("//rust/raze/remote:winapi-x86_64-pc-windows-gnu-0.4.0.BUILD.bazel")
    )

