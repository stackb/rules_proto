# Toolchains

The standard toolchain uses `@com_google_protobuf//:protoc`:

```py
register_toolchains("@build_stack_rules_proto//toolchain:standard")
```

To define an alternative, prepare a toolchain of type
`@build_stack_rules_proto//toolchain:protoc` and register that instead. See
[//toolchain:BUILD.bazel](/toolchain/BUILD.bazel) for an example.
