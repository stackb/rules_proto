# scala rules

This package registers two rules:

| Name                  | Implementation                           |
| --------------------- | ---------------------------------------- |
| `proto_scala_library` | `stackb:rules_proto:proto_scala_library` |
| `grpc_scala_library`  | `stackb:rules_proto:grpc_scala_library`  |

These rules have the following characteristics:

- They generate
  `@build_stack_rules_proto//rules/scala:{proto|grpc}_scala_library`. Both are
  thin wrappers around `@io_bazel_rules_scala//scala:scala.bzl%scala_library`.
- The rule suffix is `_scala_library`.
- `proto_scala_library` is generated if the protos have no services, otherwise
  `grpc_scala_library` is emitted.
- They merge on `srcs` and resolve on `deps`.
- They provide a littany of symbols onto the global resolver, to be
  theoretically consumed by a separate `scala` gazelle extension (see
  `provideScalaImports` for details).

Example:

```
gazelle:proto_rule proto_scala_library implementation stackb:rules_proto:proto_scala_library
gazelle:proto_rule proto_scala_library deps @maven//:com_google_protobuf_protobuf_java
gazelle:proto_rule proto_scala_library deps @maven//:com_thesamet_scalapb_lenses_2_12
gazelle:proto_rule proto_scala_library deps @maven//:com_thesamet_scalapb_scalapb_runtime_2_12
gazelle:proto_rule proto_scala_library options --noresolve=scalapb/scalapb.proto
gazelle:proto_rule proto_scala_library options --nooutput=package_scala.srcjar
gazelle:proto_rule proto_scala_library visibility //visibility:public
```

The above configuration declares a `proto_scala_library` rule that will
statically include the thee named deps, and have public visibility.

Consider a proto package that declares "package options" (other proto files in
this package are ignored for this example):

```proto
syntax = "proto2";

package example.proto;

import "thirdparty/protobuf/scalapb/scalapb.proto";

option (scalapb.options) = {
    scope: PACKAGE
    preserve_unknown_fields: false
};
```

In the typical case, the following `proto_library` and `proto_compile` rules
will be generated:

```
proto_library(
    name = "package_proto",
    srcs = ["package.proto"],
    deps = ["//thirdparty/protobuf/scalapb:scalapb_proto"],
)

proto_compile(
    name = "package_scala_compile",
    outputs = ["package_scala.srcjar"],
    plugins = ["@build_stack_rules_proto//rules/scala:protoc-gen-scala"],
    proto = "package_proto",
    visibility = ["//visibility:public"],
)
```

However, in this case we do not want to emit a `proto_scala_library`. Why?
Because the `package_scala.srcjar` in this case will be empty, as the scalapbc
plugin does not emit and corresponding JVM code for this degenerate case.

Therefore, the `--nooutput` means "suppress this output from the output list".
