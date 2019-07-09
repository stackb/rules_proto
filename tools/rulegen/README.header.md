# Protobuf and gRPC rules for Bazel

Bazel rules for building Protocol Buffers +/- gRPC :rocket:


## Contents:

- [Overview](#overview)
- [Installation](#installation)
- [Rules](#rules)
    - [Android](/android/README.md)
    - [Closure](/closure/README.md)
    - [C++](/cpp/README.md)
    - [C Sharp](/csharp/README.md)
    - [D](/d/README.md)
    - [Go](/go/README.md)
    - [GoGo](/gogo/README.md)
    - [gRPC Gateway](/github.com/grpc-ecosystem/grpc-gateway/README.md)
    - [gRPC Web](/github.com/grpc/grpc-web/README.md)
    - [gRPC.js](/github.com/stackb/grpc.js/README.md)
    - [Java](/java/README.md)
    - [NodeJs](/nodejs/README.md)
    - [Objective-C](/objc/README.md)
    - [PHP](/php/README.md)
    - [Python](/python/README.md)
    - [Ruby](/ruby/README.md)
    - [Rust](/rust/README.md)
    - [Scala](/scala/README.md)
    - [Swift](/swift/README.md)
- [Example Usage](#example-usage)
- [Developers](#developers)
    - [Code Layout](#code-layout)
    - [Rule Generation](#rule-generation)
    - [How-it-works](#how-it-works)
    - [Developing Custom Plugins](#developing-custom-plugins)
- [Contributing](#contributing)


## Overview

These rules provide [Protocol Buffers (Protobuf)](https://developers.google.com/protocol-buffers/)
and [gRPC](https://grpc.io/) rules for a range of languages and services.

Each supported language (`{lang}` below) is generally split into four rule
flavours:

- `{lang}_proto_compile`: Provides generated files from the Protobuf `protoc`
  plugin for the language. e.g for C++ this provides the generated `*.pb.cc`
  and `*.pb.h` files.

- `{lang}_proto_library`: Provides a language-specific library from the
  generated Protobuf `protoc` plugin outputs, along with necessary
  dependencies. e.g for C++ this provides a Bazel native `cpp_library` created
  from the generated `*.pb.cc` and `*pb.h` files, with the Protobuf library
  linked. For languages that do not have a 'library' concept, this rule may not
  exist.

- `{lang}_grpc_compile`: Provides generated files from both the Protobuf and
  gRPC `protoc` plugins for the language. e.g for C++ this provides the
  generated `*.pb.cc`, `*.grpc.pb.cc`, `*.pb.h` and `*.grpc.pb.h` files.

- `{lang}_proto_library`: Provides a language-specific library from the
  generated Protobuf and gRPC `protoc` plugins outputs, along with necessary
  dependencies. e.g for C++ this provides a Bazel native `cpp_library` created
  from the generated `*.pb.cc`, `*.grpc.pb.cc`, `*.pb.h` and `*.grpc.pb.h`
  files, with the Protobuf and gRPC libraries linked. For languages that do not
  have a 'library' concept, this rule may not exist.

Therefore, if you are solely interested in the generated source code artifacts,
use the `{lang}_{proto|grpc}_compile` rules. Otherwise, if you want a
ready-to-go library, use the `{lang}_{proto|grpc}_library` rules.

These rules are the successor to [rules_protobuf](https://github.com/pubref/rules_protobuf) and
are in a pre-release status. The primary goals are:

1. Interoperate with the native [`proto_library`](https://docs.bazel.build/versions/master/be/protocol-buffer.html#proto_library)
   rules and other proto support in the Bazel ecosystem as much as possible.

2. Provide a `proto_plugin` rule to support custom protoc plugins.

3. Minimal dependency loading. Proto rules should not pull in more dependencies
   than they absolutely need.

> NOTE: In general, try to use the native proto library rules when possible to
minimize external dependencies in your project. Add `rules_proto` when you have
more complex proto requirements such as when dealing with multiple output
languages, gRPC, unsupported (native) language support, or custom proto plugins.


## Installation

Add `rules_proto` your `WORKSPACE` file:

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    urls = ["https://github.com/stackb/rules_proto/archive/{{ .Ref }}.tar.gz"],
    sha256 = "{{ .Sha256 }}",
    strip_prefix = "rules_proto-{{ .Ref }}",
)
```

**Important**: Follow instructions in the language-specific `README.md` for
additional workspace dependencies that may be required.
