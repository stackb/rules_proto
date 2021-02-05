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
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/{{ .Ref }}.tar.gz"],
    sha256 = "{{ .Sha256 }}",
    strip_prefix = "rules_proto_grpc-{{ .Ref }}",
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
