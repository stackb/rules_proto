<div align="center">
    <img width="200" height="200" src="https://raw.githubusercontent.com/rules-proto-grpc/rules_proto_grpc/master/internal/resources/logo.svg">
    <h1>Protobuf and gRPC rules for <a href="https://bazel.build">Bazel</a></h1>
</div>

<div align="center">
    <a href="https://bazel.build">Bazel</a> rules for building <a href="https://developers.google.com/protocol-buffers">Protocol Buffers</a> ± <a href="https://grpc.io/">gRPC</a> code and libraries from <a href="https://docs.bazel.build/versions/master/be/protocol-buffer.html#proto_library">proto_library</a> targets<br><br>
    <a href="https://buildkite.com/bazel/rules-proto-grpc-rules-proto-grpc"><img src="https://badge.buildkite.com/a0c88e60f21c85a8bb53a8c73175aebd64f50a0d4bacbdb038.svg?branch=master"></a>
    <a href="https://github.com/rules-proto-grpc/rules_proto_grpc/actions"><img src="https://github.com/rules-proto-grpc/rules_proto_grpc/workflows/CI/badge.svg"></a>
    <img src="https://img.shields.io/github/license/rules-proto-grpc/rules_proto_grpc.svg">
</div>


## Announcements 📣

- **2020/10/11**: [Version 2.0.0 has been released](https://github.com/rules-proto-grpc/rules_proto_grpc/releases/tag/2.0.0)
  with updated Protobuf and gRPC versions. For some languages this may not be a drop-in replacement
  and it may be necessary to update your WORKSPACE file due to changes in dependencies; please see
  the above linked release notes for details or the language specific rules pages. If you discover
  any problems with the new release, please
  [open a new issue](https://github.com/rules-proto-grpc/rules_proto_grpc/issues/new).


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
    - [Go (gogoprotobuf)](/github.com/gogo/protobuf/README.md)
    - [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway/README.md)
    - [gRPC-Web](/github.com/grpc/grpc-web/README.md)
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
- [Migration](#migration)
- [Developers](#developers)
    - [Code Layout](#code-layout)
    - [Rule Generation](#rule-generation)
    - [How-it-works](#how-it-works)
    - [Developing Custom Plugins](#developing-custom-plugins)
- [License](#license)
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

- `{lang}_grpc_library`: Provides a language-specific library from the
  generated Protobuf and gRPC `protoc` plugins outputs, along with necessary
  dependencies. e.g for C++ this provides a Bazel native `cpp_library` created
  from the generated `*.pb.cc`, `*.grpc.pb.cc`, `*.pb.h` and `*.grpc.pb.h`
  files, with the Protobuf and gRPC libraries linked. For languages that do not
  have a 'library' concept, this rule may not exist.

Therefore, if you are solely interested in the generated source code artifacts,
use the `{lang}_{proto|grpc}_compile` rules. Otherwise, if you want a
ready-to-go library, use the `{lang}_{proto|grpc}_library` rules.

These rules are derived from the excellent [stackb/rules_proto](https://github.com/stackb/rules_proto)
and add aspect-based compilation to all languages, allowing for all
`proto_library` options to be expressed in user code.


## Installation

Add `rules_proto_grpc` your `WORKSPACE` file:

```starlark
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "rules_proto_grpc",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/2.0.0.tar.gz"],
    sha256 = "d771584bbff98698e7cb3cb31c132ee206a972569f4dc8b65acbdd934d156b33",
    strip_prefix = "rules_proto_grpc-2.0.0",
)

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_toolchains", "rules_proto_grpc_repos")
rules_proto_grpc_toolchains()
rules_proto_grpc_repos()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
rules_proto_dependencies()
rules_proto_toolchains()
```

It is recommended that you use the tagged releases for stable rules. Master is
intended to be 'ready-to-use', but may be unstable at certain periods. To be
notified of new releases, you can use GitHub's 'Watch Releases Only' on the
repository.

**Important**: You will also need to follow instructions in the
language-specific `README.md` for additional workspace dependencies that may be
required.


## Rules

| Language | Rule | Description
| ---: | :--- | :--- |
| [Android](/android) | [android_proto_compile](/android#android_proto_compile) | Generates an Android protobuf `.jar` artifact ([example](/example/android/android_proto_compile)) |
| [Android](/android) | [android_grpc_compile](/android#android_grpc_compile) | Generates Android protobuf+gRPC `.jar` artifacts ([example](/example/android/android_grpc_compile)) |
| [Android](/android) | [android_proto_library](/android#android_proto_library) | Generates an Android protobuf library using `android_library` from `rules_android` ([example](/example/android/android_proto_library)) |
| [Android](/android) | [android_grpc_library](/android#android_grpc_library) | Generates Android protobuf+gRPC library using `android_library` from `rules_android` ([example](/example/android/android_grpc_library)) |
| [Closure](/closure) | [closure_proto_compile](/closure#closure_proto_compile) | Generates Closure protobuf `.js` files ([example](/example/closure/closure_proto_compile)) |
| [Closure](/closure) | [closure_proto_library](/closure#closure_proto_library) | Generates a Closure library with compiled protobuf `.js` files using `closure_js_library` from `rules_closure` ([example](/example/closure/closure_proto_library)) |
| [C++](/cpp) | [cpp_proto_compile](/cpp#cpp_proto_compile) | Generates C++ protobuf `.h` & `.cc` artifacts ([example](/example/cpp/cpp_proto_compile)) |
| [C++](/cpp) | [cpp_grpc_compile](/cpp#cpp_grpc_compile) | Generates C++ protobuf+gRPC `.h` & `.cc` artifacts ([example](/example/cpp/cpp_grpc_compile)) |
| [C++](/cpp) | [cpp_proto_library](/cpp#cpp_proto_library) | Generates a C++ protobuf library using `cc_library`, with dependencies linked ([example](/example/cpp/cpp_proto_library)) |
| [C++](/cpp) | [cpp_grpc_library](/cpp#cpp_grpc_library) | Generates a C++ protobuf+gRPC library using `cc_library`, with dependencies linked ([example](/example/cpp/cpp_grpc_library)) |
| [C#](/csharp) | [csharp_proto_compile](/csharp#csharp_proto_compile) | Generates C# protobuf `.cs` artifacts ([example](/example/csharp/csharp_proto_compile)) |
| [C#](/csharp) | [csharp_grpc_compile](/csharp#csharp_grpc_compile) | Generates C# protobuf+gRPC `.cs` artifacts ([example](/example/csharp/csharp_grpc_compile)) |
| [C#](/csharp) | [csharp_proto_library](/csharp#csharp_proto_library) | Generates a C# protobuf library using `core_library` from `rules_dotnet`. Note that the library name must end in `.dll` ([example](/example/csharp/csharp_proto_library)) |
| [C#](/csharp) | [csharp_grpc_library](/csharp#csharp_grpc_library) | Generates a C# protobuf+gRPC library using `core_library` from `rules_dotnet`. Note that the library name must end in `.dll` ([example](/example/csharp/csharp_grpc_library)) |
| [D](/d) | [d_proto_compile](/d#d_proto_compile) | Generates D protobuf `.d` artifacts ([example](/example/d/d_proto_compile)) |
| [D](/d) | [d_proto_library](/d#d_proto_library) | Generates a D protobuf library using `d_library` from `rules_d` ([example](/example/d/d_proto_library)) |
| [Go](/go) | [go_proto_compile](/go#go_proto_compile) | Generates Go protobuf `.go` artifacts ([example](/example/go/go_proto_compile)) |
| [Go](/go) | [go_grpc_compile](/go#go_grpc_compile) | Generates Go protobuf+gRPC `.go` artifacts ([example](/example/go/go_grpc_compile)) |
| [Go](/go) | [go_proto_library](/go#go_proto_library) | Generates a Go protobuf library using `go_library` from `rules_go` ([example](/example/go/go_proto_library)) |
| [Go](/go) | [go_grpc_library](/go#go_grpc_library) | Generates a Go protobuf+gRPC library using `go_library` from `rules_go` ([example](/example/go/go_grpc_library)) |
| [Java](/java) | [java_proto_compile](/java#java_proto_compile) | Generates a Java protobuf srcjar artifact ([example](/example/java/java_proto_compile)) |
| [Java](/java) | [java_grpc_compile](/java#java_grpc_compile) | Generates a Java protobuf+gRPC srcjar artifact ([example](/example/java/java_grpc_compile)) |
| [Java](/java) | [java_proto_library](/java#java_proto_library) | Generates a Java protobuf library using `java_library` ([example](/example/java/java_proto_library)) |
| [Java](/java) | [java_grpc_library](/java#java_grpc_library) | Generates a Java protobuf+gRPC library using `java_library` ([example](/example/java/java_grpc_library)) |
| [Node.js](/nodejs) | [nodejs_proto_compile](/nodejs#nodejs_proto_compile) | Generates Node.js protobuf `.js` artifacts ([example](/example/nodejs/nodejs_proto_compile)) |
| [Node.js](/nodejs) | [nodejs_grpc_compile](/nodejs#nodejs_grpc_compile) | Generates Node.js protobuf+gRPC `.js` artifacts ([example](/example/nodejs/nodejs_grpc_compile)) |
| [Node.js](/nodejs) | [nodejs_proto_library](/nodejs#nodejs_proto_library) | Generates a Node.js protobuf library using `js_library` from `rules_nodejs` ([example](/example/nodejs/nodejs_proto_library)) |
| [Node.js](/nodejs) | [nodejs_grpc_library](/nodejs#nodejs_grpc_library) | Generates a Node.js protobuf+gRPC library using `js_library` from `rules_nodejs` ([example](/example/nodejs/nodejs_grpc_library)) |
| [Objective-C](/objc) | [objc_proto_compile](/objc#objc_proto_compile) | Generates Objective-C protobuf `.m` & `.h` artifacts ([example](/example/objc/objc_proto_compile)) |
| [Objective-C](/objc) | [objc_grpc_compile](/objc#objc_grpc_compile) | Generates Objective-C protobuf+gRPC `.m` & `.h` artifacts ([example](/example/objc/objc_grpc_compile)) |
| [Objective-C](/objc) | [objc_proto_library](/objc#objc_proto_library) | Generates an Objective-C protobuf library using `objc_library` ([example](/example/objc/objc_proto_library)) |
| [PHP](/php) | [php_proto_compile](/php#php_proto_compile) | Generates PHP protobuf `.php` artifacts ([example](/example/php/php_proto_compile)) |
| [PHP](/php) | [php_grpc_compile](/php#php_grpc_compile) | Generates PHP protobuf+gRPC `.php` artifacts ([example](/example/php/php_grpc_compile)) |
| [Python](/python) | [python_proto_compile](/python#python_proto_compile) | Generates Python protobuf `.py` artifacts ([example](/example/python/python_proto_compile)) |
| [Python](/python) | [python_grpc_compile](/python#python_grpc_compile) | Generates Python protobuf+gRPC `.py` artifacts ([example](/example/python/python_grpc_compile)) |
| [Python](/python) | [python_grpclib_compile](/python#python_grpclib_compile) | Generates Python protobuf+grpclib `.py` artifacts (supports Python 3 only) ([example](/example/python/python_grpclib_compile)) |
| [Python](/python) | [python_proto_library](/python#python_proto_library) | Generates a Python protobuf library using `py_library` from `rules_python` ([example](/example/python/python_proto_library)) |
| [Python](/python) | [python_grpc_library](/python#python_grpc_library) | Generates a Python protobuf+gRPC library using `py_library` from `rules_python` ([example](/example/python/python_grpc_library)) |
| [Python](/python) | [python_grpclib_library](/python#python_grpclib_library) | Generates a Python protobuf+grpclib library using `py_library` from `rules_python` (supports Python 3 only) ([example](/example/python/python_grpclib_library)) |
| [Ruby](/ruby) | [ruby_proto_compile](/ruby#ruby_proto_compile) | Generates Ruby protobuf `.rb` artifacts ([example](/example/ruby/ruby_proto_compile)) |
| [Ruby](/ruby) | [ruby_grpc_compile](/ruby#ruby_grpc_compile) | Generates Ruby protobuf+gRPC `.rb` artifacts ([example](/example/ruby/ruby_grpc_compile)) |
| [Ruby](/ruby) | [ruby_proto_library](/ruby#ruby_proto_library) | Generates a Ruby protobuf library using `ruby_library` from `rules_ruby` ([example](/example/ruby/ruby_proto_library)) |
| [Ruby](/ruby) | [ruby_grpc_library](/ruby#ruby_grpc_library) | Generates a Ruby protobuf+gRPC library using `ruby_library` from `rules_ruby` ([example](/example/ruby/ruby_grpc_library)) |
| [Rust](/rust) | [rust_proto_compile](/rust#rust_proto_compile) | Generates Rust protobuf `.rs` artifacts ([example](/example/rust/rust_proto_compile)) |
| [Rust](/rust) | [rust_grpc_compile](/rust#rust_grpc_compile) | Generates Rust protobuf+gRPC `.rs` artifacts ([example](/example/rust/rust_grpc_compile)) |
| [Rust](/rust) | [rust_proto_library](/rust#rust_proto_library) | Generates a Rust protobuf library using `rust_library` from `rules_rust` ([example](/example/rust/rust_proto_library)) |
| [Rust](/rust) | [rust_grpc_library](/rust#rust_grpc_library) | Generates a Rust protobuf+gRPC library using `rust_library` from `rules_rust` ([example](/example/rust/rust_grpc_library)) |
| [Scala](/scala) | [scala_proto_compile](/scala#scala_proto_compile) | Generates a Scala protobuf `.jar` artifact ([example](/example/scala/scala_proto_compile)) |
| [Scala](/scala) | [scala_grpc_compile](/scala#scala_grpc_compile) | Generates Scala protobuf+gRPC `.jar` artifacts ([example](/example/scala/scala_grpc_compile)) |
| [Scala](/scala) | [scala_proto_library](/scala#scala_proto_library) | Generates a Scala protobuf library using `scala_library` from `rules_scala` ([example](/example/scala/scala_proto_library)) |
| [Scala](/scala) | [scala_grpc_library](/scala#scala_grpc_library) | Generates a Scala protobuf+gRPC library using `scala_library` from `rules_scala` ([example](/example/scala/scala_grpc_library)) |
| [Swift](/swift) | [swift_proto_compile](/swift#swift_proto_compile) | Generates Swift protobuf `.swift` artifacts ([example](/example/swift/swift_proto_compile)) |
| [Swift](/swift) | [swift_grpc_compile](/swift#swift_grpc_compile) | Generates Swift protobuf+gRPC `.swift` artifacts ([example](/example/swift/swift_grpc_compile)) |
| [Swift](/swift) | [swift_proto_library](/swift#swift_proto_library) | Generates a Swift protobuf library using `swift_library` from `rules_swift` ([example](/example/swift/swift_proto_library)) |
| [Swift](/swift) | [swift_grpc_library](/swift#swift_grpc_library) | Generates a Swift protobuf+gRPC library using `swift_library` from `rules_swift` ([example](/example/swift/swift_grpc_library)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogo_proto_compile](/github.com/gogo/protobuf#gogo_proto_compile) | Generates gogo protobuf `.go` artifacts ([example](/example/github.com/gogo/protobuf/gogo_proto_compile)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogo_grpc_compile](/github.com/gogo/protobuf#gogo_grpc_compile) | Generates gogo protobuf+gRPC `.go` artifacts ([example](/example/github.com/gogo/protobuf/gogo_grpc_compile)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogo_proto_library](/github.com/gogo/protobuf#gogo_proto_library) | Generates a Go gogo protobuf library using `go_library` from `rules_go` ([example](/example/github.com/gogo/protobuf/gogo_proto_library)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogo_grpc_library](/github.com/gogo/protobuf#gogo_grpc_library) | Generates a Go gogo protobuf+gRPC library using `go_library` from `rules_go` ([example](/example/github.com/gogo/protobuf/gogo_grpc_library)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofast_proto_compile](/github.com/gogo/protobuf#gogofast_proto_compile) | Generates gogofast protobuf `.go` artifacts ([example](/example/github.com/gogo/protobuf/gogofast_proto_compile)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofast_grpc_compile](/github.com/gogo/protobuf#gogofast_grpc_compile) | Generates gogofast protobuf+gRPC `.go` artifacts ([example](/example/github.com/gogo/protobuf/gogofast_grpc_compile)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofast_proto_library](/github.com/gogo/protobuf#gogofast_proto_library) | Generates a Go gogofast protobuf library using `go_library` from `rules_go` ([example](/example/github.com/gogo/protobuf/gogofast_proto_library)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofast_grpc_library](/github.com/gogo/protobuf#gogofast_grpc_library) | Generates a Go gogofast protobuf+gRPC library using `go_library` from `rules_go` ([example](/example/github.com/gogo/protobuf/gogofast_grpc_library)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofaster_proto_compile](/github.com/gogo/protobuf#gogofaster_proto_compile) | Generates gogofaster protobuf `.go` artifacts ([example](/example/github.com/gogo/protobuf/gogofaster_proto_compile)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofaster_grpc_compile](/github.com/gogo/protobuf#gogofaster_grpc_compile) | Generates gogofaster protobuf+gRPC `.go` artifacts ([example](/example/github.com/gogo/protobuf/gogofaster_grpc_compile)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofaster_proto_library](/github.com/gogo/protobuf#gogofaster_proto_library) | Generates a Go gogofaster protobuf library using `go_library` from `rules_go` ([example](/example/github.com/gogo/protobuf/gogofaster_proto_library)) |
| [Go (gogoprotobuf)](/github.com/gogo/protobuf) | [gogofaster_grpc_library](/github.com/gogo/protobuf#gogofaster_grpc_library) | Generates a Go gogofaster protobuf+gRPC library using `go_library` from `rules_go` ([example](/example/github.com/gogo/protobuf/gogofaster_grpc_library)) |
| [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_grpc_compile](/github.com/grpc-ecosystem/grpc-gateway#gateway_grpc_compile) | Generates grpc-gateway `.go` files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_grpc_compile)) |
| [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_swagger_compile](/github.com/grpc-ecosystem/grpc-gateway#gateway_swagger_compile) | Generates grpc-gateway swagger `.json` files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_swagger_compile)) |
| [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_grpc_library](/github.com/grpc-ecosystem/grpc-gateway#gateway_grpc_library) | Generates grpc-gateway library files ([example](/example/github.com/grpc-ecosystem/grpc-gateway/gateway_grpc_library)) |
| [gRPC-Web](/github.com/grpc/grpc-web) | [closure_grpc_compile](/github.com/grpc/grpc-web#closure_grpc_compile) | Generates Closure *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/closure_grpc_compile)) |
| [gRPC-Web](/github.com/grpc/grpc-web) | [commonjs_grpc_compile](/github.com/grpc/grpc-web#commonjs_grpc_compile) | Generates CommonJS *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/commonjs_grpc_compile)) |
| [gRPC-Web](/github.com/grpc/grpc-web) | [commonjs_dts_grpc_compile](/github.com/grpc/grpc-web#commonjs_dts_grpc_compile) | Generates commonjs_dts *.js protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/commonjs_dts_grpc_compile)) |
| [gRPC-Web](/github.com/grpc/grpc-web) | [ts_grpc_compile](/github.com/grpc/grpc-web#ts_grpc_compile) | Generates CommonJS *.ts protobuf+gRPC files ([example](/example/github.com/grpc/grpc-web/ts_grpc_compile)) |
| [gRPC-Web](/github.com/grpc/grpc-web) | [closure_grpc_library](/github.com/grpc/grpc-web#closure_grpc_library) | Generates protobuf closure library *.js files ([example](/example/github.com/grpc/grpc-web/closure_grpc_library)) |

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


**Step 2**: Write a `BAZEL.build` file with a [`proto_library`](https://docs.bazel.build/versions/master/be/protocol-buffer.html#proto_library)
rule:

```starlark
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

```starlark
# BUILD.bazel
load("@rules_proto_grpc//cpp:defs.bzl", "cpp_proto_compile")

cpp_proto_compile(
    name = "cpp_thing_proto",
    deps = [":thing_proto"],
)
```

But wait, before we can build this, we need to load the dependencies necessary
for this rule (from [cpp/README.md](/cpp/README.md)):


**Step 4**: Load the workspace macro corresponding to the build rule.

```starlark
# WORKSPACE
load("@rules_proto_grpc//cpp:repositories.bzl", "cpp_repos")

cpp_repos()
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

```starlark
# BUILD.bazel
load("@rules_proto_grpc//cpp:defs.bzl", "cpp_proto_library")

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


## Migration

For users migrating from the [stackb/rules_proto](https://github.com/stackb/rules_proto)
rules, please see the help at [MIGRATION.md](/docs/MIGRATION.md)


## Developers

### Code Layout

Each language `{lang}` has a top-level subdirectory that contains:

1. `{lang}/README.md`: Generated documentation for the language rules.

1. `{lang}/repositories.bzl`: Macro functions that declare repository rule
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
  shared by all rules and is the heart of `rules_proto_grpc`; it calls `protoc`
  with a given list of plugins and generates output files.

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
   finds all the directly and transitively required proto files, along with
   their options.

2. At each node visited by the aspect, `protoc` is invoked with the relevant
   plugins and options to generate the desired outputs. The aspect uses only
   the generated proto descriptors from the `ProtoInfo` providers.

3. Once the aspect stage is complete, all generated outputs are optionally
   gathered into a final output tree.

4. For `{lang}_{proto|grpc}_library` rules, the generated outputs are then
   aggregated into a language-specific library. e.g a `.so` file for C++.


### Developing Custom Plugins

Generally, follow the pattern seen in the multiple language examples in this
repository.  The basic idea is:

1. Load the plugin rule: `load("@rules_proto_grpc//:plugin.bzl", "proto_plugin")`.
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

```starlark
load("@rules_proto_grpc//:plugin.bzl", "ProtoPluginInfo")
load(
    "@rules_proto_grpc//:aspect.bzl",
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
    toolchains = ["@rules_proto_grpc//protobuf:toolchain_type"],
)

# Create compile rule to apply aspect
_rule = rule(
    implementation = proto_compile_impl,
    attrs = dict(
        proto_compile_attrs,
        deps = attr.label_list(
            mandatory = True,
            providers = [ProtoInfo, ProtoLibraryAspectNodeInfo],
            aspects = [example_aspect],
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


## License

This project is derived from [stackb/rules_proto](https://github.com/stackb/rules_proto)
under the [Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0) license and
this project therefore maintains the terms of that license. An overview of the
changes can be found at [MIGRATION.md](/docs/MIGRATION.md).


## Contributing

Contributions are very welcome. Please see [CONTRIBUTING.md](/docs/CONTRIBUTING.md)
for further details.
