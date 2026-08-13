# Toolchains

The standard toolchain uses `@com_google_protobuf//:protoc`:

```py
register_toolchains("@build_stack_rules_proto//toolchain:standard")
```

To define an alternative, prepare a toolchain of type
`@build_stack_rules_proto//toolchain:protoc` and register that instead. See
[//toolchain:BUILD.bazel](/toolchain/BUILD.bazel) for an example.

## `--incompatible_enable_proto_toolchain_resolution`

When the
[`--incompatible_enable_proto_toolchain_resolution`](https://github.com/bazelbuild/rules_proto/discussions/213)
flag is enabled, `protoc` is instead taken from the registered proto toolchain
(`@rules_proto//proto:toolchain_type`, provided by the `proto_toolchain` rule),
which is the same toolchain `proto_library` itself uses.  The
`@build_stack_rules_proto//toolchain:protoc` toolchain is then optional and is
only consulted if no proto toolchain is registered.

This makes it possible to use a prebuilt `protoc` -- rather than building it
from source -- via a ruleset such as
[toolchains_protoc](https://github.com/aspect-build/toolchains_protoc):

```py
# MODULE.bazel
bazel_dep(name = "toolchains_protoc", version = "0.5.0")
bazel_dep(name = "build_stack_rules_proto", version = "4.1.0")

protoc = use_extension("@toolchains_protoc//protoc:extensions.bzl", "protoc")
protoc.toolchain(
    google_protobuf = "com_google_protobuf",
    version = "v29.0",
)
use_repo(protoc, "com_google_protobuf")
```

```py
# .bazelrc
common --incompatible_enable_proto_toolchain_resolution
```

With this configuration, `register_toolchains("@build_stack_rules_proto//toolchain:standard")`
is not required (registering it is harmless: the proto toolchain takes
precedence while the flag is enabled).

Note that the `protoc` attribute of the `proto_compile` rule still takes
precedence over both toolchains.
