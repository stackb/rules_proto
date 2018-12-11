# `rules_proto` [![Build Status](https://travis-ci.org/pubref/rules_proto.svg?branch=master)](https://travis-ci.org/stackb/rules_proto)

Bazel skylark rules for building protocol buffers +/- gRPC :sparkles:.

<table border="0"><tr>
<td><img src="https://bazel.build/images/bazel-icon.svg" height="180"/></td>
<td><img src="https://github.com/pubref/rules_protobuf/blob/master/images/wtfcat.png" height="180"/></td>
<td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" height="180"/></td>
</tr><tr>
<td>Bazel</td>
<td>rules_proto</td>
<td>gRPC</td>
</tr></table>

These rules are the successor to <https://github.com/pubref/rules_protobuf> and
are in a pre-release status.  The primary goals are:

1. Interoperate with the native `proto_library` rules and other proto support in
   the bazel ecosystem as much as possible.
2. Provide a `proto_plugin` rule to support custom protoc plugins.
3. Minimal dependency loading.  Proto rules should not pull in more dependencies
   than they absolutely need.

> NOTE: in general, try to use the native proto library rules when possible to
minimize external dependencies in your project.  Add `rules_proto` when you have
more complex proto requirements such as when dealing with multiple output
languages, gRPC, unsupported (native) language support, or custom proto plugins.

## Installation

Add rules_proto your `WORKSPACE`:

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    urls = ["https://github.com/stackb/rules_proto/archive/{{ .Ref }}.tar.gz"],
    sha256 = "{{ .Sha256 }}",
    strip_prefix = "rules_proto-{{ .Ref }}",
)
```

> Important: Follow instructions in the language-specific `README.md` for
additional workspace dependencies that may be required.

## Usage

**Step 1**: write a protocol buffer file (example: `thing.proto`):

```proto
syntax = "proto3";

package example;

import "google/protobuf/any.proto";

message Thing {
    string name = 1;
    google.protobuf.Any payload = 2;
}
```

**Step 2**: write a `BAZEL.build` file with a native `proto_library` rule:

```python
proto_library(
    name = "thing_proto",
    srcs = ["thing.proto"],
    deps = ["@com_google_protobuf//:any_proto"],
)
```

In this example we have a dependency on a well-known type `any.proto` hance the
`proto_library` to `proto_library` dependency.

**Step 3**: add a `cpp_proto_compile` rule (substitute `cpp_` for the language
of your choice).

> NOTE: In this example `thing.proto` does not include service definitions
(gRPC).  For protos with services, use the `cpp_grpc_compile` rule instead.

```python
# BUILD.bazel
load("@build_stack_rules_proto//cpp:cpp_proto_compile.bzl", "cpp_proto_compile")

cpp_proto_compile(
    name = "cpp_thing_proto",
    deps = [":thing_proto"],
)
```

But wait, before we can build this, we need to load the dependencies necessary
for this rule (from [cpp/README.md](/cpp/README.md)):

**Step 4**: load the workspace macro corresponding to the build rule.

```python
# WORKSPACE
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_proto_compile")

cpp_proto_compile()
```

Note that the workspace macro has the same name of the build rule.  Refer to the
[cpp/deps.bzl](/cpp/deps.bzl) for details on what other dependencies are loaded.

We're now ready to build the rule:

**Step 5**: build it.

```sh
$ bazel build //example/proto:cpp_thing_proto
Target //example/proto:cpp_thing_proto up-to-date:
  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.h
  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.cc
```

If we were only interested for the generated file artifacts, the
`cpp_grpc_compile` rule would be fine.  However, for convenience we'd rather
have the outputs compiled into an `*.so` file.  To do that, let's change the
rule from `cpp_proto_compile` to `cpp_proto_library`:

```python
# WORKSPACE
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_proto_library")

cpp_proto_library()
```

```python
# BUILD.bazel
load("@build_stack_rules_proto//cpp:cpp_proto_library.bzl", "cpp_proto_library")

cpp_proto_library(
    name = "cpp_thing_proto",
    deps = [":thing_proto"],
)
```

```sh
$ bazel build //example/proto:cpp_thing_proto
Target //example/proto:cpp_thing_proto up-to-date:
  bazel-bin/example/proto/libcpp_thing_proto.a
  bazel-bin/example/proto/libcpp_thing_proto.so  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.h
  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.cc
```

This way, we can use `//example/proto:cpp_thing_proto` as a dependency of any
other `cc_library` or `cc_binary` rule as per normal.

> NOTE: the `cpp_proto_library` implicitly calls `cpp_proto_compile`, and we can
access that rule by adding `_pb` at the end of the rule name, like `bazel build
//example/proto:cpp_thing_proto_pb`.

### Summary

* There are generally four rule flavors for any language `{lang}`:
`{lang}_proto_compile`, `{lang}_proto_library`, `{lang}_grpc_compile`, and
`{lang}_grpc_library`.  

* If you are solely interested in the source code artifacts, use the
  `{lang}_{proto|grpc}_compile` rule.  

* If your proto file has services, use the `{lang}_{grpc}_{compile|library}`
  rule instead.  

* Load any external dependencies needed for the rule via the
  `load("@build_stack_rules_proto//{lang}:deps.bzl",
  "{lang}_{proto|grpc}_{compile|library}")`.  
