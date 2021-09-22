# `rules_proto`

Bazel starlark rules for building protocol buffers +/- gRPC :sparkles:.

<img src="data:image/svg+xml,%3C!--%20Created%20by%20Jose%20Rivera%20--%3E%0A%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20xmlns%3Axlink%3D%22http%3A%2F%2Fwww.w3.org%2F1999%2Fxlink%22%20version%3D%221.1%22%20x%3D%220px%22%20y%3D%220px%22%20viewBox%3D%220%200%20100%20100%22%20style%3D%22enable-background%3Anew%200%200%20100%20100%3B%22%20xml%3Aspace%3D%22preserve%22%3E%0A%09%3Cg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M68.9%2C52.4L68.9%2C52.4l5.3%2C3.4L68.9%2C52.4z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cpath%20d%3D%22M90.9%2C69.5l0-0.2l0%2C0.8L90.9%2C69.5z%22%20%2F%3E%0A%09%09%3Cpath%20d%3D%22M16.2%2C34.4L15.7%2C34l-0.2-0.2v0l0.2%2C0.1L16.2%2C34.4L16.2%2C34.4L16.2%2C34.4z%20M16.2%2C34.4L15.7%2C34l-0.2-0.2v0l0.2%2C0.1L16.2%2C34.4%20%20%20L16.2%2C34.4L16.2%2C34.4z%20M15.2%2C33.5L15.2%2C33.5l0.2%2C0.2L15.2%2C33.5z%20M16.2%2C34.4L15.7%2C34l-0.2-0.2v0l0.2%2C0.1L16.2%2C34.4L16.2%2C34.4%20%20%20L16.2%2C34.4z%20M16.2%2C34.4L15.7%2C34h0L16.2%2C34.4L16.2%2C34.4L16.2%2C34.4z%22%20%2F%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M16.2%2C34.4L15.7%2C34l-0.2-0.2v0l0.2%2C0.1L16.2%2C34.4L16.2%2C34.4L16.2%2C34.4z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cpath%20d%3D%22M52.4%2C31.2L52.4%2C31.2L52.4%2C31.2L52.4%2C31.2z%22%20%2F%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M68%2C32.8L68%2C32.8l0.7%2C0.6L68%2C32.8z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M94.7%2C82.9l-2.9%2C2.9l-0.9-15.8l0-0.8l-0.2-0.2l-10.1-9.1l4.6-6.4l-0.9-6.7l-9.5-8.3l-0.6-0.5l-5.5-4.7L68%2C32.8v0%20%20%20%20L56.2%2C22.5l-0.4-0.3l-22.9%2C5l-11.7%2C2.6l-0.6%2C0.1l-0.2%2C0.2l-4.2%2C2.7l0.8-10L13%2C28.7l1.3-11.2l0-0.1L13%2C9.4L12%2C19.6L8.3%2C34.4%20%20%20%20l-8%2C16.2l2.4-1.4l3.7%2C1.9l1-4.6h0l13.4-7.8l0.1-0.1l0.3%2C0.1l13.4%2C2.2L29%2C43.6l-0.3%2C0.2l0.1%2C0.7l2.9%2C14.2l0.2%2C0.8l0.7%2C0.6l8%2C6.3%20%20%20%20l-0.2-5.7l-7.8-1.6l1-15.2l9.6-1.5l1.2-0.2l0.2%2C0.8l0.2%2C0.1l0.4%2C0.2l0.4%2C0.2l0.3%2C0.1l22.2%2C9.5L72%2C63.5l12.8%2C4.4L91%2C86.7l0.1%2C0.4%20%20%20%20l0.8%2C0.3l7.8%2C3.2L94.7%2C82.9z%20M6.7%2C45.9L6.7%2C45.9l-4.5%2C2.6l5.2-10.4L9.9%2C33l0%2C0L8.6%2C38l-0.2%2C0.7L6.7%2C45.9z%20M19.7%2C38.4l-12.1%2C7%20%20%20%20l1.7-6.9l3.6-1.3L10.7%2C33l0.7-2.9l3.7%2C3.4l0%2C0l0.2%2C0.2l0.1%2C0.1l0%2C0l0.2%2C0.2l0.5%2C0.5l0.1%2C0.1l3.8%2C3.5l0.1%2C0.1L19.7%2C38.4z%20%20%20%20%20M21.3%2C37.6v-7L34%2C27.8l7.8%2C4.5L21.3%2C37.6z%20M45.6%2C42.4l-0.3-0.1l-0.1-0.4L45%2C41.1l-0.8-4l-0.2-0.8l-0.7-3.5l8.1-0.7L45.6%2C42.4z%20%20%20%20%20M52.5%2C31.2L52.5%2C31.2L52.5%2C31.2L52.5%2C31.2z%20M52.8%2C31.6L52.8%2C31.6l0.1-0.4h0l3-7.9l12%2C10.4l-1.9%2C10.7L52.8%2C31.6z%20M74.1%2C55.8%20%20%20%20l-5.3-3.4v0l0%2C0l-2.2-7l7.5-6.2l9.4%2C8.1l0.8%2C6.1l-4.4%2C6.1L74.1%2C55.8z%20M85.8%2C68.3l4.3%2C1.5l0.8%2C14L85.8%2C68.3z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M15.2%2C33.5L15.2%2C33.5l0.2%2C0.2L15.2%2C33.5z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M15.2%2C33.5L15.2%2C33.5l0.2%2C0.2L15.2%2C33.5z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpolygon%20points%3D%2252.5%2C31.2%2052.5%2C31.2%2052.4%2C31.2%20%20%20%22%20%2F%3E%0A%09%09%09%3Cpolygon%20points%3D%2252.5%2C31.2%2052.5%2C31.2%2052.4%2C31.2%20%20%20%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%3C%2Fg%3E%0A%3C%2Fsvg%3E%0A"/>

<table border="0"><tr>
<td><img src="https://bazel.build/images/bazel-icon.svg" height="180"/></td>
<td><img src="https://github.com/pubref/rules_protobuf/blob/master/images/wtfcat.png" height="180"/></td>
<td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" height="180"/></td>
</tr><tr>
<td>Bazel</td>
<td>rules_proto</td>
<td>gRPC</td>
</tr></table>

`@build_stack_rules_proto` provides:

1. Rules for driving the `protoc` tool within a bazel workspace.
2. A [gazelle] extension that generates rules based on the content of your
   `.proto` files.
3. Example setups for a variety of languages.

# Getting Started

## Repository rule

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Commit: 5bbf4640487c6f3167b79c9277a377d80ec7ba3d
# Date: 2021-09-22 03:32:09 +0000 UTC
# URL: https://github.com/stackb/rules_proto/commit/5bbf4640487c6f3167b79c9277a377d80ec7ba3d
#
# Bump npm google-protobuf back up to 3.18.0 (#179)
# Size: 306266 (306 kB)
http_archive(
    name = "build_stack_rules_proto",
    sha256 = "50bc1d9c5b8436d75ee09fb1386621835c04ca26f8a0020946d672e4427d7eba",
    strip_prefix = "rules_proto-5bbf4640487c6f3167b79c9277a377d80ec7ba3d",
    urls = ["https://github.com/stackb/rules_proto/archive/5bbf4640487c6f3167b79c9277a377d80ec7ba3d.tar.gz"],
)
```

## `WORKSPACE` Boilerplate

```python
register_toolchains("@build_stack_rules_proto//toolchain:standard")
```

> This prepares `protoc` for the `proto_compile` rule.  For simple setups, consider `@build_stack_rules_proto//toolchain:prebuilt` to skip compilation of the tool.

---

```python
load("@build_stack_rules_proto//deps:core_deps.bzl", "core_deps")

core_deps()
```

> This brings in `@io_bazel_rules_go`, `@bazel_gazelle`, and `@rules_proto` if
> you don't already have them.

> The gazelle extension and associated golang dependencies are optional; you can
> write `proto_compile` and other rules by hand (but but why would you not want
> to?).

---

```python
load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(version = "1.16.2")
```

> Standard biolerplate for `@io_bazel_rules_go`.

---

```python
load( "@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
```

> Standard boilerplate for `@bazel_gazelle`.

---

```python
load("//:go_deps.bzl", "gazelle_protobuf_extension_go_deps")

gazelle_protobuf_extension_go_deps()
```

> This brings in `@com_github_emicklei_proto`.
> [github.com/emicklei/proto](https://github.com/emicklei/proto) is used by the
> gazelle extension to parse proto files.

---

```python
load("//deps:protobuf_core_deps.bzl", "protobuf_core_deps")

protobuf_core_deps()
```

> This brings in `@com_google_protobuf` and friends if you don't already have them.

## Gazelle Setup

```python
load("@bazel_gazelle//:def.bzl", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle-protobuf",
    languages = [
        "@bazel_gazelle//language/go",
        "@bazel_gazelle//language/proto",
        # must be after the proto extension (order matters)
        "@build_stack_rules_proto//language/protobuf",
    ],
)

gazelle(
    name = "gazelle",
    gazelle = ":gazelle-protobuf",
)
```

> The gazelle setup is typically placed in the root `BUILD.bazel` file.

---

## Gazelle Configuration

Gazelle is configured by special comments in BUILD files called `directives`.

> Gazelle works by visiting each package in your workspace; configuration is
> done "on the way in" whereas actual rule generation is done "on the way out".
> Therefore, gazelle configuration if a subdirectory inherits that from its
> parents.  As such, directives placed in the root `BUILD.bazel` file apply the
> entire workspace.

```python
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# gazelle:proto_plugin cpp implementation builtin:cpp
# gazelle:proto_plugin protoc-gen-grpc-cpp implementation grpc:grpc:cpp

# gazelle:proto_rule proto_cc_library implementation stackb:rules_proto:proto_cc_library
# gazelle:proto_rule proto_cc_library deps @com_google_protobuf//:protobuf
# gazelle:proto_rule proto_cc_library visibility //visibility:public
# gazelle:proto_rule grpc_cc_library implementation stackb:rules_proto:grpc_cc_library
# gazelle:proto_rule grpc_cc_library deps @com_github_grpc_grpc//:grpc++
# gazelle:proto_rule grpc_cc_library deps @com_github_grpc_grpc//:grpc++_reflection
# gazelle:proto_rule grpc_cc_library visibility //visibility:public

# gazelle:proto_language cpp plugin cpp
# gazelle:proto_language cpp plugin protoc-gen-grpc-cpp
# gazelle:proto_language cpp rule proto_compile
# gazelle:proto_language cpp rule proto_cc_library
# gazelle:proto_language cpp rule grpc_cc_library
```

Let's unpack this a bit:

- `gazelle:proto_plugin cpp implementation builtin:cpp` associates the name
  `cpp` with a piece of golang code that implements the `protoc.Plugin`
  interface. The extension maintains a registry of these actors. The gazelle
  extension ships with a number of them out of the box, but you can also write
  your own.  The core responsibility a `protoc.Plugin` implementation is to to
  *predict* the files that a protoc plugin tool will generate for an individual
  `proto_library` rule.  The implemention has full read access to the
  `protoc.File`s in the `proto_library` to be able to predict *if* a file will
  be generated and *where* it will appear in the filesystem (specifically, the
  bazel execution root). 
- `gazelle:proto_rule proto_compile implementation
  stackb:rules_proto:proto_compile` associates the name `proto_compile` with a
  piece of golang code that implements the `protoc.Rule` interface.  The
  extension maintains a registry of rule implementations. Similarly, the
  extension ships with a bunch of them out of the box, but you can write your
  own custom rules as well.  The core responsibility a `protoc.Rule`
  implementation is construct a gazelle `rule.Rule` based upon a `proto_library`
  rule and the set of plugins that are configured with it.
- `gazelle:proto_language cpp plugin cpp` instantiates a `protoc.LanguageConfig`
  having the name `cpp` and adds the `cpp` plugin to it.  The language
  configuration bundles bundles plugins and rules together.
- `gazelle:proto_rule grpc_cc_library deps @com_github_grpc_grpc//:grpc++`
  configures the rule such that all generated rules will have that dependency.

> **+/- intent modifiers**.  Although not pictured in this example, many of the
> directives take an *intent modifier* to turn configuration on/off.  For
> example, if you wanted to suppress the grpc c++ plugin in the package
> `//proto/javaapi`, put a directive like `gazelle:proto_language cpp rule
> -grpc_cc_library` in `proto/javaapi/BUILD.bazel` (note the `-` symbol
> preceding the name).  To suppress the language entirely, use
> `gazelle:proto_language cpp enabled false`.

## Gazelle Execution

Now that we have the `WORKSPACE` setup and gazelle configured, we can run
gazelle:

```sh
$ bazel run //:gazelle -- proto/
```

To restrict gazelle to only a particular subdirectory `example/routeguide/`:

```sh
$ bazel run //:gazelle -- example/routeguide/
```

Gazelle should now have generated something like the following:

```python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules/cc:grpc_cc_library.bzl", "grpc_cc_library")
load("@build_stack_rules_proto//rules/cc:proto_cc_library.bzl", "proto_cc_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

proto_library(
    name = "routeguide_proto",
    srcs = ["routeguide.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "routeguide_cpp_compile",
    outputs = [
        "routeguide.grpc.pb.cc",
        "routeguide.grpc.pb.h",
        "routeguide.pb.cc",
        "routeguide.pb.h",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/builtin:cpp",
        "@build_stack_rules_proto//plugin/grpc/grpc:protoc-gen-grpc-cpp",
    ],
    proto = "routeguide_proto",
)

proto_cc_library(
    name = "routeguide_cc_library",
    srcs = ["routeguide.pb.cc"],
    hdrs = ["routeguide.pb.h"],
    visibility = ["//visibility:public"],
    deps = ["@com_google_protobuf//:protobuf"],
)

grpc_cc_library(
    name = "routeguide_grpc_cc_library",
    srcs = ["routeguide.grpc.pb.cc"],
    hdrs = ["routeguide.grpc.pb.h"],
    visibility = ["//visibility:public"],
    deps = [
        ":routeguide_cc_library",
        "@com_github_grpc_grpc//:grpc++",
        "@com_github_grpc_grpc//:grpc++_reflection",
    ],
)
```

Regarding rules like
`@build_stack_rules_proto//rules/cc:proto_cc_library.bzl%proto_cc_library"`.
These are nearly always extremely thin wrappers for the "real" rule.  For
example, here's the implementation in `proto_cc_library.bzl`:

```python
load("@rules_cc//cc:defs.bzl", "cc_library")

def proto_cc_library(**kwargs):
    cc_library(**kwargs)
```

An implementation detail of gazelle itself is that two different language
extensions should not *claim* the same load namespace, so in order to prevent
potential conflicts with other possible gazelle extensions, using the name
`@rules_cc//cc:defs.bzl%cc_library` is undesirable.

You can use a directive like `gazelle:proto_rule label //:my_fancy_cc_library.bzl` 


## Core Rules

The heart of `stackb/rules_proto` contains two rules:

| Rule            | Description                                      |
|-----------------|--------------------------------------------------|
| `proto_compile` | Executes the `protoc` tool.                      |
| `proto_plugin`  | Provides `protoc` plugin-specific configuration. |

Example:

```python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

proto_library(
    name = "thing_proto",
    srcs = ["thing.proto"],
    deps = ["@com_google_protobuf//:timestamp_proto"],
)

proto_plugin(name = "cpp")

proto_compile(
    name = "person_cpp_compile",
    outputs = [
        "person.pb.cc",
        "person.pb.h",
    ],
    plugins = [":cpp"],
    proto = "person_proto",
)
```

Takeaways:

- A `proto_library` rule forms the basis for other language-specific derived
  rules.
- `proto_library` is provided by
  [bazelbuild/rules_proto](https://github.com/bazelbuild/rules_proto).
- `proto_plugin` provides the plugin tool and other options and configuration.
  This example is the simplest case where the plugin is builtin to `protoc`
  itself; no separate plugin tool is required.
- A `proto_compile` rule references a single `proto_library` target. 
- The `plugins` attribute is a list of labels to `proto_plugin` targets.
- The `outputs` attribute names the files that will be generated by the protoc
  invocation.
- The
  [proto](https://github.com/bazelbuild/bazel-gazelle/blob/master/language/proto/lang.go)
  extension provided by [bazel-gazelle] is responsible for generating
  `proto_library`.

## Plugin Implementations

The plugin name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/plugin_name convention.

| Plugin                                                | Link     |
|-------------------------------------------------------|----------|
| `builtin:cpp`                                         | [link]() |
| `builtin:csharp`                                      | [link]() |
| `builtin:java`                                        | [link]() |
| `builtin:js:closure`                                  | [link]() |
| `builtin:js:common`                                   | [link]() |
| `builtin:objc`                                        | [link]() |
| `builtin:php`                                         | [link]() |
| `builtin:python`                                      | [link]() |
| `builtin:ruby`                                        | [link]() |
| `grpc:grpc:cpp`                                       | [link]() |
| `grpc:grpc:protoc-gen-grpc-python`                    | [link]() |
| `grpc:grpc-go:protoc-gen-go`                          | [link]() |
| `grpc:grpc-go:protoc-gen-go-grpc`                     | [link]() |
| `grpc:grpc-java:protoc-gen-grpc-java`                 | [link]() |
| `grpc:grpc-node:protoc-gen-grpc-node`                 | [link]() |
| `gogo:protobuf:protoc-gen-combo`                      | [link]() |
| `gogo:protobuf:protoc-gen-gogo`                       | [link]() |
| `gogo:protobuf:protoc-gen-gogofast`                   | [link]() |
| `gogo:protobuf:protoc-gen-gogofaster`                 | [link]() |
| `gogo:protobuf:protoc-gen-gogoslick`                  | [link]() |
| `gogo:protobuf:protoc-gen-gogotypes`                  | [link]() |
| `gogo:protobuf:protoc-gen-gostring`                   | [link]() |
| `grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway` | [link]() |
| `scalapb:scalapb:protoc-gen-scala`                    | [link]() |
| `stackb:grpc.js:protoc-gen-grpc-js`                   | [link]() |
| `stephenh:ts-proto:protoc-gen-ts-proto`               | [link]() |
| `bazelbuild:rules_proto:proto_descriptor_set`         | [repo]() |

## Rule Implementations

The rule name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/rule_name convention.

| Plugin                                        | Link     |
|-----------------------------------------------|----------|
| `stackb:rules_proto:grpc_cc_library`          | [link]() |
| `stackb:rules_proto:grpc_closure_js_library`  | [link]() |
| `stackb:rules_proto:grpc_java_library`        | [link]() |
| `stackb:rules_proto:grpc_nodejs_library`      | [link]() |
| `stackb:rules_proto:grpc_py_library`          | [link]() |
| `stackb:rules_proto:grpc_scala_library`       | [link]() |
| `stackb:rules_proto:proto_cc_library`         | [link]() |
| `stackb:rules_proto:proto_closure_js_library` | [link]() |
| `stackb:rules_proto:proto_compile`            | [link]() |
| `stackb:rules_proto:proto_compiled_sources`   | [link]() |
| `stackb:rules_proto:proto_descriptor_set`     | [link]() |
| `stackb:rules_proto:proto_go_library`         | [link]() |
| `stackb:rules_proto:proto_java_library`       | [link]() |
| `stackb:rules_proto:proto_nodejs_library`     | [link]() |
| `stackb:rules_proto:proto_py_library`         | [link]() |
| `stackb:rules_proto:proto_scala_library`      | [link]() |
| `bazelbuild:rules_scala:scala_proto_library`  | [link]() |

Please consult the `example/` directory and unit tests for more additional
detail.