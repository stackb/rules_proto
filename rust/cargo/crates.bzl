"""
cargo-raze crate workspace functions

DO NOT EDIT! Replaced on runs of cargo-raze
"""

def raze_fetch_remote_crates():

    native.new_http_archive(
        name = "raze__arrayvec__0_4_7",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/arrayvec/arrayvec-0.4.7.crate",
        type = "tar.gz",
        strip_prefix = "arrayvec-0.4.7",
        build_file = "//rust/cargo/remote:arrayvec-0.4.7.BUILD"
    )

    native.new_http_archive(
        name = "raze__base64__0_9_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/base64/base64-0.9.3.crate",
        type = "tar.gz",
        strip_prefix = "base64-0.9.3",
        build_file = "//rust/cargo/remote:base64-0.9.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__bitflags__1_0_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/bitflags/bitflags-1.0.4.crate",
        type = "tar.gz",
        strip_prefix = "bitflags-1.0.4",
        build_file = "//rust/cargo/remote:bitflags-1.0.4.BUILD"
    )

    native.new_http_archive(
        name = "raze__byteorder__1_2_6",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/byteorder/byteorder-1.2.6.crate",
        type = "tar.gz",
        strip_prefix = "byteorder-1.2.6",
        build_file = "//rust/cargo/remote:byteorder-1.2.6.BUILD"
    )

    native.new_http_archive(
        name = "raze__bytes__0_4_10",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/bytes/bytes-0.4.10.crate",
        type = "tar.gz",
        strip_prefix = "bytes-0.4.10",
        build_file = "//rust/cargo/remote:bytes-0.4.10.BUILD"
    )

    native.new_http_archive(
        name = "raze__cfg_if__0_1_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/cfg-if/cfg-if-0.1.5.crate",
        type = "tar.gz",
        strip_prefix = "cfg-if-0.1.5",
        build_file = "//rust/cargo/remote:cfg-if-0.1.5.BUILD"
    )

    native.new_http_archive(
        name = "raze__cloudabi__0_0_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/cloudabi/cloudabi-0.0.3.crate",
        type = "tar.gz",
        strip_prefix = "cloudabi-0.0.3",
        build_file = "//rust/cargo/remote:cloudabi-0.0.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__crossbeam_deque__0_6_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam-deque/crossbeam-deque-0.6.1.crate",
        type = "tar.gz",
        strip_prefix = "crossbeam-deque-0.6.1",
        build_file = "//rust/cargo/remote:crossbeam-deque-0.6.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__crossbeam_epoch__0_5_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam-epoch/crossbeam-epoch-0.5.2.crate",
        type = "tar.gz",
        strip_prefix = "crossbeam-epoch-0.5.2",
        build_file = "//rust/cargo/remote:crossbeam-epoch-0.5.2.BUILD"
    )

    native.new_http_archive(
        name = "raze__crossbeam_utils__0_5_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/crossbeam-utils/crossbeam-utils-0.5.0.crate",
        type = "tar.gz",
        strip_prefix = "crossbeam-utils-0.5.0",
        build_file = "//rust/cargo/remote:crossbeam-utils-0.5.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__fuchsia_zircon__0_3_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/fuchsia-zircon/fuchsia-zircon-0.3.3.crate",
        type = "tar.gz",
        strip_prefix = "fuchsia-zircon-0.3.3",
        build_file = "//rust/cargo/remote:fuchsia-zircon-0.3.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__fuchsia_zircon_sys__0_3_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/fuchsia-zircon-sys/fuchsia-zircon-sys-0.3.3.crate",
        type = "tar.gz",
        strip_prefix = "fuchsia-zircon-sys-0.3.3",
        build_file = "//rust/cargo/remote:fuchsia-zircon-sys-0.3.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__futures__0_1_24",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/futures/futures-0.1.24.crate",
        type = "tar.gz",
        strip_prefix = "futures-0.1.24",
        build_file = "//rust/cargo/remote:futures-0.1.24.BUILD"
    )

    native.new_http_archive(
        name = "raze__futures_cpupool__0_1_8",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/futures-cpupool/futures-cpupool-0.1.8.crate",
        type = "tar.gz",
        strip_prefix = "futures-cpupool-0.1.8",
        build_file = "//rust/cargo/remote:futures-cpupool-0.1.8.BUILD"
    )

    native.new_http_archive(
        name = "raze__grpc__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/grpc/grpc-0.4.0.crate",
        type = "tar.gz",
        strip_prefix = "grpc-0.4.0",
        build_file = "//rust/cargo/remote:grpc-0.4.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__grpc_compiler__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/grpc-compiler/grpc-compiler-0.4.0.crate",
        type = "tar.gz",
        strip_prefix = "grpc-compiler-0.4.0",
        build_file = "//rust/cargo/remote:grpc-compiler-0.4.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__httpbis__0_6_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/httpbis/httpbis-0.6.1.crate",
        type = "tar.gz",
        strip_prefix = "httpbis-0.6.1",
        build_file = "//rust/cargo/remote:httpbis-0.6.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__iovec__0_1_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/iovec/iovec-0.1.2.crate",
        type = "tar.gz",
        strip_prefix = "iovec-0.1.2",
        build_file = "//rust/cargo/remote:iovec-0.1.2.BUILD"
    )

    native.new_http_archive(
        name = "raze__kernel32_sys__0_2_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/kernel32-sys/kernel32-sys-0.2.2.crate",
        type = "tar.gz",
        strip_prefix = "kernel32-sys-0.2.2",
        build_file = "//rust/cargo/remote:kernel32-sys-0.2.2.BUILD"
    )

    native.new_http_archive(
        name = "raze__lazy_static__1_1_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/lazy_static/lazy_static-1.1.0.crate",
        type = "tar.gz",
        strip_prefix = "lazy_static-1.1.0",
        build_file = "//rust/cargo/remote:lazy_static-1.1.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__lazycell__1_2_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/lazycell/lazycell-1.2.0.crate",
        type = "tar.gz",
        strip_prefix = "lazycell-1.2.0",
        build_file = "//rust/cargo/remote:lazycell-1.2.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__libc__0_2_43",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/libc/libc-0.2.43.crate",
        type = "tar.gz",
        strip_prefix = "libc-0.2.43",
        build_file = "//rust/cargo/remote:libc-0.2.43.BUILD"
    )

    native.new_http_archive(
        name = "raze__lock_api__0_1_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/lock_api/lock_api-0.1.3.crate",
        type = "tar.gz",
        strip_prefix = "lock_api-0.1.3",
        build_file = "//rust/cargo/remote:lock_api-0.1.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__log__0_3_9",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/log/log-0.3.9.crate",
        type = "tar.gz",
        strip_prefix = "log-0.3.9",
        build_file = "//rust/cargo/remote:log-0.3.9.BUILD"
    )

    native.new_http_archive(
        name = "raze__log__0_4_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/log/log-0.4.5.crate",
        type = "tar.gz",
        strip_prefix = "log-0.4.5",
        build_file = "//rust/cargo/remote:log-0.4.5.BUILD"
    )

    native.new_http_archive(
        name = "raze__memoffset__0_2_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/memoffset/memoffset-0.2.1.crate",
        type = "tar.gz",
        strip_prefix = "memoffset-0.2.1",
        build_file = "//rust/cargo/remote:memoffset-0.2.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__mio__0_6_16",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/mio/mio-0.6.16.crate",
        type = "tar.gz",
        strip_prefix = "mio-0.6.16",
        build_file = "//rust/cargo/remote:mio-0.6.16.BUILD"
    )

    native.new_http_archive(
        name = "raze__mio_uds__0_6_7",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/mio-uds/mio-uds-0.6.7.crate",
        type = "tar.gz",
        strip_prefix = "mio-uds-0.6.7",
        build_file = "//rust/cargo/remote:mio-uds-0.6.7.BUILD"
    )

    native.new_http_archive(
        name = "raze__miow__0_2_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/miow/miow-0.2.1.crate",
        type = "tar.gz",
        strip_prefix = "miow-0.2.1",
        build_file = "//rust/cargo/remote:miow-0.2.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__net2__0_2_33",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/net2/net2-0.2.33.crate",
        type = "tar.gz",
        strip_prefix = "net2-0.2.33",
        build_file = "//rust/cargo/remote:net2-0.2.33.BUILD"
    )

    native.new_http_archive(
        name = "raze__nodrop__0_1_12",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/nodrop/nodrop-0.1.12.crate",
        type = "tar.gz",
        strip_prefix = "nodrop-0.1.12",
        build_file = "//rust/cargo/remote:nodrop-0.1.12.BUILD"
    )

    native.new_http_archive(
        name = "raze__num_cpus__1_8_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/num_cpus/num_cpus-1.8.0.crate",
        type = "tar.gz",
        strip_prefix = "num_cpus-1.8.0",
        build_file = "//rust/cargo/remote:num_cpus-1.8.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__owning_ref__0_3_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/owning_ref/owning_ref-0.3.3.crate",
        type = "tar.gz",
        strip_prefix = "owning_ref-0.3.3",
        build_file = "//rust/cargo/remote:owning_ref-0.3.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__parking_lot__0_6_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/parking_lot/parking_lot-0.6.4.crate",
        type = "tar.gz",
        strip_prefix = "parking_lot-0.6.4",
        build_file = "//rust/cargo/remote:parking_lot-0.6.4.BUILD"
    )

    native.new_http_archive(
        name = "raze__parking_lot_core__0_3_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/parking_lot_core/parking_lot_core-0.3.1.crate",
        type = "tar.gz",
        strip_prefix = "parking_lot_core-0.3.1",
        build_file = "//rust/cargo/remote:parking_lot_core-0.3.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__protobuf__1_6_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/protobuf/protobuf-1.6.0.crate",
        type = "tar.gz",
        strip_prefix = "protobuf-1.6.0",
        build_file = "//rust/cargo/remote:protobuf-1.6.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__protobuf_codegen__1_6_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/protobuf-codegen/protobuf-codegen-1.6.0.crate",
        type = "tar.gz",
        strip_prefix = "protobuf-codegen-1.6.0",
        build_file = "//rust/cargo/remote:protobuf-codegen-1.6.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__rand__0_5_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand/rand-0.5.5.crate",
        type = "tar.gz",
        strip_prefix = "rand-0.5.5",
        build_file = "//rust/cargo/remote:rand-0.5.5.BUILD"
    )

    native.new_http_archive(
        name = "raze__rand_core__0_2_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rand_core/rand_core-0.2.1.crate",
        type = "tar.gz",
        strip_prefix = "rand_core-0.2.1",
        build_file = "//rust/cargo/remote:rand_core-0.2.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__rustc_version__0_2_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/rustc_version/rustc_version-0.2.3.crate",
        type = "tar.gz",
        strip_prefix = "rustc_version-0.2.3",
        build_file = "//rust/cargo/remote:rustc_version-0.2.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__safemem__0_3_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/safemem/safemem-0.3.0.crate",
        type = "tar.gz",
        strip_prefix = "safemem-0.3.0",
        build_file = "//rust/cargo/remote:safemem-0.3.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__scoped_tls__0_1_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/scoped-tls/scoped-tls-0.1.2.crate",
        type = "tar.gz",
        strip_prefix = "scoped-tls-0.1.2",
        build_file = "//rust/cargo/remote:scoped-tls-0.1.2.BUILD"
    )

    native.new_http_archive(
        name = "raze__scopeguard__0_3_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/scopeguard/scopeguard-0.3.3.crate",
        type = "tar.gz",
        strip_prefix = "scopeguard-0.3.3",
        build_file = "//rust/cargo/remote:scopeguard-0.3.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__semver__0_9_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/semver/semver-0.9.0.crate",
        type = "tar.gz",
        strip_prefix = "semver-0.9.0",
        build_file = "//rust/cargo/remote:semver-0.9.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__semver_parser__0_7_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/semver-parser/semver-parser-0.7.0.crate",
        type = "tar.gz",
        strip_prefix = "semver-parser-0.7.0",
        build_file = "//rust/cargo/remote:semver-parser-0.7.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__slab__0_3_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/slab/slab-0.3.0.crate",
        type = "tar.gz",
        strip_prefix = "slab-0.3.0",
        build_file = "//rust/cargo/remote:slab-0.3.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__slab__0_4_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/slab/slab-0.4.1.crate",
        type = "tar.gz",
        strip_prefix = "slab-0.4.1",
        build_file = "//rust/cargo/remote:slab-0.4.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__smallvec__0_6_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/smallvec/smallvec-0.6.5.crate",
        type = "tar.gz",
        strip_prefix = "smallvec-0.6.5",
        build_file = "//rust/cargo/remote:smallvec-0.6.5.BUILD"
    )

    native.new_http_archive(
        name = "raze__stable_deref_trait__1_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/stable_deref_trait/stable_deref_trait-1.1.1.crate",
        type = "tar.gz",
        strip_prefix = "stable_deref_trait-1.1.1",
        build_file = "//rust/cargo/remote:stable_deref_trait-1.1.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__tls_api__0_1_20",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tls-api/tls-api-0.1.20.crate",
        type = "tar.gz",
        strip_prefix = "tls-api-0.1.20",
        build_file = "//rust/cargo/remote:tls-api-0.1.20.BUILD"
    )

    native.new_http_archive(
        name = "raze__tls_api_stub__0_1_20",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tls-api-stub/tls-api-stub-0.1.20.crate",
        type = "tar.gz",
        strip_prefix = "tls-api-stub-0.1.20",
        build_file = "//rust/cargo/remote:tls-api-stub-0.1.20.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio__0_1_8",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio/tokio-0.1.8.crate",
        type = "tar.gz",
        strip_prefix = "tokio-0.1.8",
        build_file = "//rust/cargo/remote:tokio-0.1.8.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_codec__0_1_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-codec/tokio-codec-0.1.0.crate",
        type = "tar.gz",
        strip_prefix = "tokio-codec-0.1.0",
        build_file = "//rust/cargo/remote:tokio-codec-0.1.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_core__0_1_17",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-core/tokio-core-0.1.17.crate",
        type = "tar.gz",
        strip_prefix = "tokio-core-0.1.17",
        build_file = "//rust/cargo/remote:tokio-core-0.1.17.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_current_thread__0_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-current-thread/tokio-current-thread-0.1.1.crate",
        type = "tar.gz",
        strip_prefix = "tokio-current-thread-0.1.1",
        build_file = "//rust/cargo/remote:tokio-current-thread-0.1.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_executor__0_1_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-executor/tokio-executor-0.1.4.crate",
        type = "tar.gz",
        strip_prefix = "tokio-executor-0.1.4",
        build_file = "//rust/cargo/remote:tokio-executor-0.1.4.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_fs__0_1_3",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-fs/tokio-fs-0.1.3.crate",
        type = "tar.gz",
        strip_prefix = "tokio-fs-0.1.3",
        build_file = "//rust/cargo/remote:tokio-fs-0.1.3.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_io__0_1_8",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-io/tokio-io-0.1.8.crate",
        type = "tar.gz",
        strip_prefix = "tokio-io-0.1.8",
        build_file = "//rust/cargo/remote:tokio-io-0.1.8.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_reactor__0_1_5",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-reactor/tokio-reactor-0.1.5.crate",
        type = "tar.gz",
        strip_prefix = "tokio-reactor-0.1.5",
        build_file = "//rust/cargo/remote:tokio-reactor-0.1.5.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_tcp__0_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-tcp/tokio-tcp-0.1.1.crate",
        type = "tar.gz",
        strip_prefix = "tokio-tcp-0.1.1",
        build_file = "//rust/cargo/remote:tokio-tcp-0.1.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_threadpool__0_1_6",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-threadpool/tokio-threadpool-0.1.6.crate",
        type = "tar.gz",
        strip_prefix = "tokio-threadpool-0.1.6",
        build_file = "//rust/cargo/remote:tokio-threadpool-0.1.6.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_timer__0_1_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-timer/tokio-timer-0.1.2.crate",
        type = "tar.gz",
        strip_prefix = "tokio-timer-0.1.2",
        build_file = "//rust/cargo/remote:tokio-timer-0.1.2.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_timer__0_2_6",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-timer/tokio-timer-0.2.6.crate",
        type = "tar.gz",
        strip_prefix = "tokio-timer-0.2.6",
        build_file = "//rust/cargo/remote:tokio-timer-0.2.6.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_tls_api__0_1_20",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-tls-api/tokio-tls-api-0.1.20.crate",
        type = "tar.gz",
        strip_prefix = "tokio-tls-api-0.1.20",
        build_file = "//rust/cargo/remote:tokio-tls-api-0.1.20.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_udp__0_1_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-udp/tokio-udp-0.1.2.crate",
        type = "tar.gz",
        strip_prefix = "tokio-udp-0.1.2",
        build_file = "//rust/cargo/remote:tokio-udp-0.1.2.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_uds__0_1_7",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-uds/tokio-uds-0.1.7.crate",
        type = "tar.gz",
        strip_prefix = "tokio-uds-0.1.7",
        build_file = "//rust/cargo/remote:tokio-uds-0.1.7.BUILD"
    )

    native.new_http_archive(
        name = "raze__tokio_uds__0_2_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/tokio-uds/tokio-uds-0.2.1.crate",
        type = "tar.gz",
        strip_prefix = "tokio-uds-0.2.1",
        build_file = "//rust/cargo/remote:tokio-uds-0.2.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__unix_socket__0_5_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/unix_socket/unix_socket-0.5.0.crate",
        type = "tar.gz",
        strip_prefix = "unix_socket-0.5.0",
        build_file = "//rust/cargo/remote:unix_socket-0.5.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__unreachable__1_0_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/unreachable/unreachable-1.0.0.crate",
        type = "tar.gz",
        strip_prefix = "unreachable-1.0.0",
        build_file = "//rust/cargo/remote:unreachable-1.0.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__version_check__0_1_4",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/version_check/version_check-0.1.4.crate",
        type = "tar.gz",
        strip_prefix = "version_check-0.1.4",
        build_file = "//rust/cargo/remote:version_check-0.1.4.BUILD"
    )

    native.new_http_archive(
        name = "raze__void__1_0_2",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/void/void-1.0.2.crate",
        type = "tar.gz",
        strip_prefix = "void-1.0.2",
        build_file = "//rust/cargo/remote:void-1.0.2.BUILD"
    )

    native.new_http_archive(
        name = "raze__winapi__0_2_8",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi/winapi-0.2.8.crate",
        type = "tar.gz",
        strip_prefix = "winapi-0.2.8",
        build_file = "//rust/cargo/remote:winapi-0.2.8.BUILD"
    )

    native.new_http_archive(
        name = "raze__winapi__0_3_6",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi/winapi-0.3.6.crate",
        type = "tar.gz",
        strip_prefix = "winapi-0.3.6",
        build_file = "//rust/cargo/remote:winapi-0.3.6.BUILD"
    )

    native.new_http_archive(
        name = "raze__winapi_build__0_1_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi-build/winapi-build-0.1.1.crate",
        type = "tar.gz",
        strip_prefix = "winapi-build-0.1.1",
        build_file = "//rust/cargo/remote:winapi-build-0.1.1.BUILD"
    )

    native.new_http_archive(
        name = "raze__winapi_i686_pc_windows_gnu__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi-i686-pc-windows-gnu/winapi-i686-pc-windows-gnu-0.4.0.crate",
        type = "tar.gz",
        strip_prefix = "winapi-i686-pc-windows-gnu-0.4.0",
        build_file = "//rust/cargo/remote:winapi-i686-pc-windows-gnu-0.4.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__winapi_x86_64_pc_windows_gnu__0_4_0",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/winapi-x86_64-pc-windows-gnu/winapi-x86_64-pc-windows-gnu-0.4.0.crate",
        type = "tar.gz",
        strip_prefix = "winapi-x86_64-pc-windows-gnu-0.4.0",
        build_file = "//rust/cargo/remote:winapi-x86_64-pc-windows-gnu-0.4.0.BUILD"
    )

    native.new_http_archive(
        name = "raze__ws2_32_sys__0_2_1",
        url = "https://crates-io.s3-us-west-1.amazonaws.com/crates/ws2_32-sys/ws2_32-sys-0.2.1.crate",
        type = "tar.gz",
        strip_prefix = "ws2_32-sys-0.2.1",
        build_file = "//rust/cargo/remote:ws2_32-sys-0.2.1.BUILD"
    )

