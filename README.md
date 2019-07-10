# Protobuf and gRPC rules for Bazel

Bazel rules for building Protocol Buffers +/- gRPC :rocket:


## Contents:

- [Overview](#overview)
- [Installation](#installation)
- [Rules](#rules)
    - [Android](/android/README.md)
    - [Closure](/closure/README.md)
    - [C++](/cpp/README.md)
    - [C#](/csharp/README.md)
    - [D](/d/README.md)
    - [Go](/go/README.md)
    - [Go (gogoprotobuf)](/gogo/README.md)
    - [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway/README.md)
    - [gRPC-Web](/github.com/grpc/grpc-web/README.md)
    - [grpc.js](/github.com/stackb/grpc.js/README.md)
    - [Java](/java/README.md)
    - [Node.js](/nodejs/README.md)
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
    urls = ["https://github.com/stackb/rules_proto/archive/89cb1da964299d7c8901657b341b7a9ff7d83e39.tar.gz"],
    sha256 = "395408a3dc9c3db2b5c200b8722a13a60898c861633b99e6e250186adffd1370",
    strip_prefix = "rules_proto-89cb1da964299d7c8901657b341b7a9ff7d83e39",
)
```

**Important**: Follow instructions in the language-specific `README.md` for
additional workspace dependencies that may be required.


## Rules

| Status | Language | Rule | Description
| ---    | ---: | :--- | :--- |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Android](/android) | [android_proto_compile](/android#android_proto_compile) | Generates android protobuf artifacts ([example](/example/android/android_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Android](/android) | [android_grpc_compile](/android#android_grpc_compile) | Generates android protobuf+gRPC artifacts ([example](/example/android/android_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Android](/android) | [android_proto_library](/android#android_proto_library) | Generates android protobuf library ([example](/example/android/android_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Android](/android) | [android_grpc_library](/android#android_grpc_library) | Generates android protobuf+gRPC library ([example](/example/android/android_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Closure](/closure) | [closure_proto_compile](/closure#closure_proto_compile) | Generates closure *.js protobuf+gRPC files ([example](/example/closure/closure_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Closure](/closure) | [closure_proto_library](/closure#closure_proto_library) | Generates a closure_library with compiled protobuf *.js files ([example](/example/closure/closure_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [C++](/cpp) | [cpp_proto_compile](/cpp#cpp_proto_compile) | Generates *.h,*.cc protobuf artifacts ([example](/example/cpp/cpp_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [C++](/cpp) | [cpp_grpc_compile](/cpp#cpp_grpc_compile) | Generates *.h,*.cc protobuf+gRPC artifacts ([example](/example/cpp/cpp_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [C++](/cpp) | [cpp_proto_library](/cpp#cpp_proto_library) | Generates *.h,*.cc protobuf library ([example](/example/cpp/cpp_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [C++](/cpp) | [cpp_grpc_library](/cpp#cpp_grpc_library) | Generates *.h,*.cc protobuf+gRPC library ([example](/example/cpp/cpp_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [C#](/csharp) | [csharp_proto_compile](/csharp#csharp_proto_compile) | Generates csharp protobuf artifacts ([example](/example/csharp/csharp_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [C#](/csharp) | [csharp_grpc_compile](/csharp#csharp_grpc_compile) | Generates csharp protobuf+gRPC artifacts ([example](/example/csharp/csharp_grpc_compile)) |
| - | [C#](/csharp) | [csharp_proto_library](/csharp#csharp_proto_library) | Generates csharp protobuf library ([example](/example/csharp/csharp_proto_library)) |
| - | [C#](/csharp) | [csharp_grpc_library](/csharp#csharp_grpc_library) | Generates csharp protobuf+gRPC library ([example](/example/csharp/csharp_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [D](/d) | [d_proto_compile](/d#d_proto_compile) | Generates d protobuf artifacts ([example](/example/d/d_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [D](/d) | [d_proto_library](/d#d_proto_library) | Generates d protobuf library ([example](/example/d/d_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go](/go) | [go_proto_compile](/go#go_proto_compile) | Generates *.go protobuf artifacts ([example](/example/go/go_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go](/go) | [go_grpc_compile](/go#go_grpc_compile) | Generates *.go protobuf+gRPC artifacts ([example](/example/go/go_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go](/go) | [go_proto_library](/go#go_proto_library) | Generates *.go protobuf library ([example](/example/go/go_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go](/go) | [go_grpc_library](/go#go_grpc_library) | Generates *.go protobuf+gRPC library ([example](/example/go/go_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Java](/java) | [java_proto_compile](/java#java_proto_compile) | Generates a srcjar with protobuf *.java files ([example](/example/java/java_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Java](/java) | [java_grpc_compile](/java#java_grpc_compile) | Generates a srcjar with protobuf+gRPC *.java files ([example](/example/java/java_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Java](/java) | [java_proto_library](/java#java_proto_library) | Generates a jar with compiled protobuf *.class files ([example](/example/java/java_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Java](/java) | [java_grpc_library](/java#java_grpc_library) | Generates a jar with compiled protobuf+gRPC *.class files ([example](/example/java/java_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Node.js](/nodejs) | [nodejs_proto_compile](/nodejs#nodejs_proto_compile) | Generates Node.js *.js protobuf artifacts ([example](/example/nodejs/nodejs_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Node.js](/nodejs) | [nodejs_grpc_compile](/nodejs#nodejs_grpc_compile) | Generates Node.js *.js protobuf+gRPC artifacts ([example](/example/nodejs/nodejs_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Objective-C](/objc) | [objc_proto_compile](/objc#objc_proto_compile) | Generates objc protobuf artifacts ([example](/example/objc/objc_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Objective-C](/objc) | [objc_grpc_compile](/objc#objc_grpc_compile) | Generates objc protobuf+gRPC artifacts ([example](/example/objc/objc_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Objective-C](/objc) | [objc_proto_library](/objc#objc_proto_library) | Generates objc protobuf library ([example](/example/objc/objc_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [PHP](/php) | [php_proto_compile](/php#php_proto_compile) | Generates php protobuf artifacts ([example](/example/php/php_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [PHP](/php) | [php_grpc_compile](/php#php_grpc_compile) | Generates php protobuf+gRPC artifacts ([example](/example/php/php_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Python](/python) | [python_proto_compile](/python#python_proto_compile) | Generates *.py protobuf artifacts ([example](/example/python/python_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Python](/python) | [python_grpc_compile](/python#python_grpc_compile) | Generates *.py protobuf+gRPC artifacts ([example](/example/python/python_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Python](/python) | [python_proto_library](/python#python_proto_library) | Generates *.py protobuf library ([example](/example/python/python_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Python](/python) | [python_grpc_library](/python#python_grpc_library) | Generates *.py protobuf+gRPC library ([example](/example/python/python_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Ruby](/ruby) | [ruby_proto_compile](/ruby#ruby_proto_compile) | Generates *.ruby protobuf artifacts ([example](/example/ruby/ruby_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Ruby](/ruby) | [ruby_grpc_compile](/ruby#ruby_grpc_compile) | Generates *.ruby protobuf+gRPC artifacts ([example](/example/ruby/ruby_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Ruby](/ruby) | [ruby_proto_library](/ruby#ruby_proto_library) | Generates *.rb protobuf library ([example](/example/ruby/ruby_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Ruby](/ruby) | [ruby_grpc_library](/ruby#ruby_grpc_library) | Generates *.rb protobuf+gRPC library ([example](/example/ruby/ruby_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Rust](/rust) | [rust_proto_compile](/rust#rust_proto_compile) | Generates rust protobuf artifacts ([example](/example/rust/rust_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Rust](/rust) | [rust_grpc_compile](/rust#rust_grpc_compile) | Generates rust protobuf+gRPC artifacts ([example](/example/rust/rust_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Rust](/rust) | [rust_proto_library](/rust#rust_proto_library) | Generates rust protobuf library ([example](/example/rust/rust_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Rust](/rust) | [rust_grpc_library](/rust#rust_grpc_library) | Generates rust protobuf+gRPC library ([example](/example/rust/rust_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Scala](/scala) | [scala_proto_compile](/scala#scala_proto_compile) | Generates *.scala protobuf artifacts ([example](/example/scala/scala_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Scala](/scala) | [scala_grpc_compile](/scala#scala_grpc_compile) | Generates *.scala protobuf+gRPC artifacts ([example](/example/scala/scala_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Scala](/scala) | [scala_proto_library](/scala#scala_proto_library) | Generates *.scala protobuf library ([example](/example/scala/scala_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Scala](/scala) | [scala_grpc_library](/scala#scala_grpc_library) | Generates *.scala protobuf+gRPC library ([example](/example/scala/scala_grpc_library)) |
| - | [Swift](/swift) | [swift_proto_compile](/swift#swift_proto_compile) | Generates swift protobuf artifacts ([example](/example/swift/swift_proto_compile)) |
| - | [Swift](/swift) | [swift_grpc_compile](/swift#swift_grpc_compile) | Generates swift protobuf+gRPC artifacts ([example](/example/swift/swift_grpc_compile)) |
| - | [Swift](/swift) | [swift_proto_library](/swift#swift_proto_library) | Generates swift protobuf library ([example](/example/swift/swift_proto_library)) |
| - | [Swift](/swift) | [swift_grpc_library](/swift#swift_grpc_library) | Generates swift protobuf+gRPC library ([example](/example/swift/swift_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogo_proto_compile](/github.com/gogo/protobuf#gogo_proto_compile) | Generates gogo protobuf artifacts ([example](/example/github.com/gogo/protobuf/gogo_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogo_grpc_compile](/github.com/gogo/protobuf#gogo_grpc_compile) | Generates gogo protobuf+gRPC artifacts ([example](/example/github.com/gogo/protobuf/gogo_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogo_proto_library](/github.com/gogo/protobuf#gogo_proto_library) | Generates gogo protobuf library ([example](/example/github.com/gogo/protobuf/gogo_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogo_grpc_library](/github.com/gogo/protobuf#gogo_grpc_library) | Generates gogo protobuf+gRPC library ([example](/example/github.com/gogo/protobuf/gogo_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofast_proto_compile](/github.com/gogo/protobuf#gogofast_proto_compile) | Generates gogofast protobuf artifacts ([example](/example/github.com/gogo/protobuf/gogofast_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofast_grpc_compile](/github.com/gogo/protobuf#gogofast_grpc_compile) | Generates gogofast protobuf+gRPC artifacts ([example](/example/github.com/gogo/protobuf/gogofast_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofast_proto_library](/github.com/gogo/protobuf#gogofast_proto_library) | Generates gogofast protobuf library ([example](/example/github.com/gogo/protobuf/gogofast_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofast_grpc_library](/github.com/gogo/protobuf#gogofast_grpc_library) | Generates gogofast protobuf+gRPC library ([example](/example/github.com/gogo/protobuf/gogofast_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofaster_proto_compile](/github.com/gogo/protobuf#gogofaster_proto_compile) | Generates gogofaster protobuf artifacts ([example](/example/github.com/gogo/protobuf/gogofaster_proto_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofaster_grpc_compile](/github.com/gogo/protobuf#gogofaster_grpc_compile) | Generates gogofaster protobuf+gRPC artifacts ([example](/example/github.com/gogo/protobuf/gogofaster_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofaster_proto_library](/github.com/gogo/protobuf#gogofaster_proto_library) | Generates gogofaster protobuf library ([example](/example/github.com/gogo/protobuf/gogofaster_proto_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofaster_grpc_library](/github.com/gogo/protobuf#gogofaster_grpc_library) | Generates gogofaster protobuf+gRPC library ([example](/example/github.com/gogo/protobuf/gogofaster_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_grpc_compile](/github.com/grpc-ecosystem/grpc-gateway#gateway_grpc_compile) | Generates grpc-gateway *.go files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_swagger_compile](/github.com/grpc-ecosystem/grpc-gateway#gateway_swagger_compile) | Generates grpc-gateway swagger *.json files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_swagger_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_grpc_library](/github.com/grpc-ecosystem/grpc-gateway#gateway_grpc_library) | Generates grpc-gateway library files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc.js](/github.com/stackb/grpc.js) | [grpcjs_grpc_compile](/github.com/stackb/grpc.js#grpcjs_grpc_compile) | Generates protobuf closure grpc *.js files ([example](/example/github.com/stackb/grpc.js/grpcjs_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [grpc.js](/github.com/stackb/grpc.js) | [grpcjs_grpc_library](/github.com/stackb/grpc.js#grpcjs_grpc_library) | Generates protobuf closure library *.js files ([example](/example/github.com/stackb/grpc.js/grpcjs_grpc_library)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gRPC-Web](/github.com/grpc/grpc-web) | [closure_grpc_compile](/github.com/grpc/grpc-web#closure_grpc_compile) | Generates a closure *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/closure_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gRPC-Web](/github.com/grpc/grpc-web) | [commonjs_grpc_compile](/github.com/grpc/grpc-web#commonjs_grpc_compile) | Generates a commonjs *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/commonjs_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gRPC-Web](/github.com/grpc/grpc-web) | [commonjs_dts_grpc_compile](/github.com/grpc/grpc-web#commonjs_dts_grpc_compile) | Generates a commonjs_dts *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/commonjs_dts_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gRPC-Web](/github.com/grpc/grpc-web) | [ts_grpc_compile](/github.com/grpc/grpc-web#ts_grpc_compile) | Generates a commonjs *.ts protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/ts_grpc_compile)) |
| [![Build Status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/rules-proto) | [gRPC-Web](/github.com/grpc/grpc-web) | [closure_grpc_library](/github.com/grpc/grpc-web#closure_grpc_library) | Generates protobuf closure library *.js files ([example](/example/github.com/grpc/grpc-web/closure_grpc_library)) |

## Example Usage

These steps walk through the steps to go from a raw `.proto` file to a C++
library:


**Step 1**: Write a Protocol Buffer file (example: `thing.proto`):

```proto
syntax = "proto3";

package example;

import "google/protobuf/any.proto";

message Thing {
    string name = 1;
    google.protobuf.Any payload = 2;
}
```


**Step 2**: Write a `BAZEL.build` file with a native [`proto_library`](https://docs.bazel.build/versions/master/be/protocol-buffer.html#proto_library)
rule:

```python
proto_library(
    name = "thing_proto",
    srcs = ["thing.proto"],
    deps = ["@com_google_protobuf//:any_proto"],
)
```

In this example we have a dependency on a well-known type `any.proto`, hence the
`proto_library` to `proto_library` dependency (`"@com_google_protobuf//:any_proto"`)


**Step 3**: Add a `cpp_proto_compile` rule (substitute `cpp_` for the language
of your choice).

> NOTE: In this example `thing.proto` does not include service definitions
(gRPC).  For protos with services, use the `cpp_grpc_compile` rule instead.

```python
# BUILD.bazel
load("@build_stack_rules_proto//cpp:defs.bzl", "cpp_proto_compile")

cpp_proto_compile(
    name = "cpp_thing_proto",
    deps = [":thing_proto"],
)
```

But wait, before we can build this, we need to load the dependencies necessary
for this rule (from [cpp/README.md](/cpp/README.md)):


**Step 4**: Load the workspace macro corresponding to the build rule.

```python
# WORKSPACE
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_deps")

cpp_deps()
```

We're now ready to build the rule:


**Step 5**: Build it!

```sh
$ bazel build //example/proto:cpp_thing_proto
Target //example/proto:cpp_thing_proto up-to-date:
  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.h
  bazel-genfiles/example/proto/cpp_thing_proto/example/proto/thing.pb.cc
```

If we were only interested in the generated file artifacts, the
`cpp_grpc_compile` rule would be fine. However, for convenience we'd rather
have the outputs compiled into an `*.so` file. To do that, let's change the
rule from `cpp_proto_compile` to `cpp_proto_library`:

```python
# BUILD.bazel
load("@build_stack_rules_proto//cpp:defs.bzl", "cpp_proto_library")

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


## Developers

### Code Layout

Each language `{lang}` has a top-level subdirectory that contains:

1. `{lang}/README.md`: Generated documentation for the language rules.

1. `{lang}/deps.bzl`: Macro functions that declare repository rule
   dependencies for that language.

2. `{lang}/{rule}.bzl`: Rule implementations of the form
   `{lang}_{kind}_{type}`, where `kind` is one of `proto|grpc` and `type` is one
   of `compile|library`.

3. `{lang}/BUILD.bazel`: `proto_plugin()` declarations for the available
   plugins for the language.

4. `example/{lang}/{rule}/`: Generated `WORKSPACE` and `BUILD.bazel`
   demonstrating standalone usage of the rules.

5. `{lang}/example/routeguide/`: Example routeguide example implementation, if
   possible.


The repository root directory contains the base rule defintions:

* `plugin.bzl`: A build rule that defines the name, tool binary, and options for
  a particular proto plugin.

* `aspect.bzl`: Contains the implementation of the compilation aspect. This is
  shared by all rules and is the heart of `rules_proto`; it calls `protoc` with
  a given list of plugins and generates output files.

Additional protoc plugins and their rules are scoped to the github repository
name where the plugin resides.


### Rule Generation

To help maintain consistency of the rule implementations and documentation, all
of the rule implementations are generated by the tool `//tools/rulegen`. Changes
in the main `README.md` should be placed in `tools/rulegen/README.header.md` or
`tools/rulegen/README.footer.md`. Changes to generated rules should be put in
the source files (example: `tools/rulegen/java.go`).


### How-it-works

Briefly, here's how the rules work:

1. Using the `proto_library` graph, an aspect walks through the [`ProtoInfo`](https://docs.bazel.build/versions/master/skylark/lib/ProtoInfo.html)
   providers on the `deps` attribute to `{lang}_{proto|grpc}_compile`. This
   finds all the directly and transitively required proto files., along with
   their options.

2. At each node visited by the aspect, `protoc` is invoked with the relevant
   plugins and options to generate the desired outputs.

3. Once the aspect stage is complete, all generated outputs are optionally
   gathered into a final output tree.

4. For `{lang}_{proto|grpc}_library` rules, the generated outputs are then
   aggregated into a language-specific library. e.g a `.so` file for C++.


### Developing Custom Plugins

Generally, follow the pattern seen in the multiple language examples in this
repository.  The basic idea is:

1. Load the plugin rule: `load("@build_stack_rules_proto//:plugin.bzl", "proto_plugin")`.
2. Define the rule, giving it a `name`, `options` (not mandatory), `tool`, and
   `outputs`.
3. `tool` is a label that refers to the binary executable for the plugin itself.
4. Choose your output type (pick one!):
    - `outputs`: a list of strings patterns that predicts the pattern of files
      generated by the plugin. For plugins that produce one output file per
      input proto file
    - `out`: the name of a single output file generated by the plugin.
    - `output_directory`: Set to true if your plugin generates files in a
      non-predictable way. e.g. if the output paths depend on the service names.
5. Create a compilation rule and aspect using the following template:

```python
load("@build_stack_rules_proto//:plugin.bzl", "ProtoPluginInfo")
load(
    "@build_stack_rules_proto//:aspect.bzl",
    "ProtoLibraryAspectNodeInfo",
    "proto_compile_aspect_attrs",
    "proto_compile_aspect_impl",
    "proto_compile_attrs",
    "proto_compile_impl",
)

# Create aspect
example_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    provides = [ProtoLibraryAspectNodeInfo],
    attr_aspects = ["deps"],
    attrs = dict(
        proto_compile_aspect_attrs,
        _plugins = attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                Label("//<LABEL OF YOUR PLUGIN>"),
            ],
        ),
        _prefix = attr.string(
            doc = "String used to disambiguate aspects when generating outputs",
            default = "example_aspect",
        )
    ),
    toolchains = ["@build_stack_rules_proto//protobuf:toolchain_type"],
)

# Create compile rule to apply aspect
_rule = rule(
    implementation = proto_compile_impl,
    attrs = dict(
        proto_compile_attrs,
        deps = attr.label_list(
            mandatory = True,
            providers = [ProtoInfo, ProtoLibraryAspectNodeInfo],
            aspects = [example_compile],
        ),
    ),
)

# Create macro for converting attrs and passing to compile
def example_compile(**kwargs):
    _rule(
        verbose_string = "{}".format(kwargs.get("verbose", 0)),
        merge_directories = True,
        **{k: v for k, v in kwargs.items() if k != "merge_directories"}
    )

```


## Contributing

Contributions welcome; please create Issues or PRs.
