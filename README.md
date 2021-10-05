# `rules_proto`

Bazel starlark rules for building protocol buffers +/- gRPC :sparkles:.

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
2. A [gazelle](https://github.com/bazelbuild/bazel-gazelle/) extension that
   generates rules based on the content of your `.proto` files.
3. A repository rule that runs gazelle in an external workspace.
4. Example setups for a variety of languages.

# Table of Contents
  - [Getting Started](#getting-started)
    - [Repository Rule](#repository-rule)
    - [Workspace Boilerplate](#workspace-boilerplate)
    - [Gazelle Setup](#gazelle-setup)
    - [Gazelle Configuration](#gazelle-configuration)
      - [Build Directives](#build-directives)
      - [YAML Config File](#yaml-configuration)
    - [Running Gazelle](#running-gazelle)
  - [Core Rules](#core-rules)
    - [proto_compile](#proto_compile)
      - [Deep dive on the mappings attribute](#the-mappings-attribute)
    - [proto_plugin](#proto_compile)
    - [proto_repository](#proto-repository)
  - [Plugin Implementations](#plugin-implementations)
  - [Rule Implementations](#rule-implementations)
  - [History of this repository](#history)

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

> This prepares `protoc` for the `proto_compile` rule.  For simple setups,
> consider `@build_stack_rules_proto//toolchain:prebuilt` to skip compilation of
> the tool.

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

The gazelle extension can be configured using "build directives" and/or a YAML file.

### Build Directives
Gazelle is configured by special comments in BUILD files called *directives*.

> Gazelle works by visiting each package in your workspace; configuration is
> done "on the way in" whereas actual rule generation is done "on the way out".
> Gazelle configuration of a subdirectory inherits that from its
> parents.  As such, directives placed in the root `BUILD.bazel` file apply to
> the entire workspace.

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
  interface. The extension maintains a registry of these actors (the gazelle
  extension ships with a number of them out of the box, but you can also write
  your own).  The core responsibility a `protoc.Plugin` implementation is to to
  *predict* the files that a protoc plugin tool will generate for an individual
  `proto_library` rule.  The implemention has full read access to the
  `protoc.File`s in the `proto_library` to be able to predict *if* a file will
  be generated and *where* it will appear in the filesystem (specifically,
  relative to the bazel execution root during a `proto_compile` action). 
- `gazelle:proto_rule proto_compile implementation
  stackb:rules_proto:proto_compile` associates the name `proto_compile` with a
  piece of golang code that implements the `protoc.LanguageRule` interface.  The
  extension maintains a registry of rule implementations. Similarly, the
  extension ships with a bunch of them out of the box, but you can write your
  own custom rules as well.  The core responsibility a `protoc.LanguageRule`
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

### YAML Configuration

You can also configure the extension using a YAML file.  This is semantically similar to adding gazelle directives to the root `BUILD` file; the YAML configuration applies to all downstream packages.  The equivalent YAML config for the above directives looks like:

```yaml
plugins:
  - name: cpp
    implementation: builtin:cpp
  - name: protoc-gen-grpc-cpp
    implementation: grpc:grpc:cpp
rules:
  - name: proto_compile
    implementation: stackb:rules_proto:proto_compile
    visibility:
      -  //visibility:public
  - name: proto_cc_library
    implementation: stackb:rules_proto:proto_cc_library
    visibility:
      -  //visibility:public
    deps:
      - "@com_google_protobuf//:protobuf"
  - name: grpc_cc_library
    implementation: stackb:rules_proto:grpc_cc_library
    visibility:
      -  //visibility:public
    deps:
      - "@com_github_grpc_grpc//:grpc++"
      - "@com_github_grpc_grpc//:grpc++_reflection"
languages:
  - name: "cpp"
    plugins:
      - cpp
      - protoc-gen-grpc-cpp
    rules:
      - proto_compile
      - proto_cc_library
      - grpc_cc_library
```

> A yaml config is particularly useful in conjunction with the `proto_repository` rule, for example to apply a set of custom plugins over 
> the googleapis/googleapis repo.

To use this in a gazelle rule, specify `-proto_language_config_file` in `args`:

```python
gazelle(
    name = "gazelle",
    gazelle = ":gazelle-protobuf",
    args = [
        "-proto_language_config_file=example/config.yaml",
    ],
)
```


## Running Gazelle

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
These are nearly always very thin wrappers for the "real" rule.  For example,
here's the implementation in `proto_cc_library.bzl`:

```python
load("@rules_cc//cc:defs.bzl", "cc_library")

def proto_cc_library(**kwargs):
    cc_library(**kwargs)
```

An implementation detail of gazelle itself is that two different language
extensions should not *claim* the same load namespace, so in order to prevent
potential conflicts with other possible gazelle extensions, using the name
`@rules_cc//cc:defs.bzl%cc_library` is undesirable.

## Core Rules

The heart of `stackb/rules_proto` contains two build rules and one repository rule:

| Rule               | Description                                      |
| ------------------ | ------------------------------------------------ |
| `proto_compile`    | Executes the `protoc` tool.                      |
| `proto_plugin`     | Provides `protoc` plugin-specific configuration. |
| `proto_repository` | Generate BUILD files for an external repository. |

## proto_compile

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

### The `mappings` attribute

Consider the following rule within the package `example/thing`:

```python
proto_compile(
    name = "thing_go_compile",
    mappings = {"thing.pb.go": "github.com/stackb/rules_proto/example/thing/thing.pb.go"},
    outputs = ["thing.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/grpc/grpc-go:protoc-gen-go"],
    proto = "thing_proto",
)
```

This rule is declaring that a file `bazel-bin/example/thing/thing.pb.go` will be
output when the action is run. When we `bazel build
//example/thing:thing_go_compile`, the file is indeed created.

Let's temporarily comment out the `mappings` attribute and rebuild:

```python
proto_compile(
    name = "thing_go_compile",
    # mappings = {"thing.pb.go": "github.com/stackb/rules_proto/example/thing/thing.pb.go"},
    outputs = ["thing.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/grpc/grpc-go:protoc-gen-go"],
    proto = "thing_proto",
)
```

```sh
$ bazel build //example/thing:thing_go_compile
ERROR: /github.com/stackb/rules_proto/example/thing/BUILD.bazel:54:14: output 'example/thing/thing.pb.go' was not created
```

What happened?  Let's add a debugging attribute `verbose = True` on the rule: this will print debugging information and show the bazel sandbox before and after the `protoc` tool is invoked:

```python
proto_compile(
    name = "thing_go_compile",
    # mappings = {"thing.pb.go": "github.com/stackb/rules_proto/example/thing/thing.pb.go"},
    outputs = ["thing.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/grpc/grpc-go:protoc-gen-go"],
    proto = "thing_proto",
    verbose = True,
)
```

```sh
$ bazel build //example/thing:thing_go_compile
##### SANDBOX BEFORE RUNNING PROTOC
./bazel-out/host/bin/external/com_google_protobuf/protoc
./bazel-out/darwin-opt-exec-2B5CBBC6/bin/external/com_github_golang_protobuf/protoc-gen-go/protoc-gen-go_/protoc-gen-go
./bazel-out/darwin-fastbuild/bin/example/thing/thing_proto-descriptor-set.proto.bin
./bazel-out/darwin-fastbuild/bin/external/com_google_protobuf/timestamp_proto-descriptor-set.proto.bin

##### SANDBOX AFTER RUNNING PROTOC
./bazel-out/darwin-fastbuild/bin/github.com/stackb/rules_proto/example/thing/thing.pb.go
```

So, the file was created, but not in the location we wanted.  In this case the
`protoc-gen-go` plugin is not "playing nice" with Bazel.  Because this
`thing.proto` has `option go_package =
"github.com/stackb/rules_proto/example/thing;thing";`, the output location is no
longer based on the `package`.  This is a problem, because Bazel semantics
disallow declaring a File outside its package boundary.  As a result, we need to
do a `mv
./bazel-out/darwin-fastbuild/bin/github.com/stackb/rules_proto/example/thing/thing.pb.go
./bazel-out/darwin-fastbuild/bin/example/thing/thing.pb.go` to relocate the
file into its expected location before the action terminates.

Therefore, the `mappings` attribute is a dict that maps file locations `{ want:
got }` relative to the action execution root.  It is required when the actual
output location does not match the desired location.  This can occur if the
proto `package` statement does not match the Bazel package path, or in special
circumstances specific to the plugin itself (like `go_package`).

## proto_repository

From an implementation standpoint, this is very similar to the `go_repository` rule.  Both can download external files and then run the gazelle generator over the downloaded files.  Example:

```python
proto_repository(
    name = "proto_googleapis",
    build_directives = [
        "gazelle:resolve proto google/api/http.proto //google/api:http_proto",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    proto_language_config_file = "//example:config.yaml",
    strip_prefix = "googleapis-02710fa0ea5312d79d7fb986c9c9823fb41049a9",
    type = "zip",
    urls = ["https://codeload.github.com/googleapis/googleapis/zip/02710fa0ea5312d79d7fb986c9c9823fb41049a9"],
)
```

Takeaways: 

- The `urls`, `strip_prefix` and `type` behave similarly to `http_archive`.
- `build_file_proto_mode` is the same the `go_repository` attribute of the same name; additionally the value `file` is permitted which generates a `proto_library` for each individual proto file.
- `build_file_generation` is the same the `go_repository` attribute of the same name; additionally the value `clean` is supported.  For example, googleapis already has a set of BUILD files; the `clean` mode will remove all the existing build files before generating new ones.
- `build_directives` is the same as `go_repository`.  Resolve directives in this case are needed because the gazelle `language/proto` extension hardcodes a proto import like `google/api/http.proto` to resolve to the `@go_googleapis` workspace; here we want to make sure that http.proto resolves to the same external workspace.
- `proto_language_config_file` is an optional label pointing to a valid `config.yaml` file to configure this extension.

With this sample configuration, the following build command succeeds:

```bash
$ bazel build @proto_googleapis//google/api:annotations_cc_library
Target @proto_googleapis//google/api:annotations_cc_library up-to-date:
  bazel-bin/external/proto_googleapis/google/api/libannotations_cc_library.a
  bazel-bin/external/proto_googleapis/google/api/libannotations_cc_library.so
```

## Plugin Implementations

The plugin name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/plugin_name convention.

| Plugin                                                | Link     |
| ----------------------------------------------------- | -------- |
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
| `golang:protobuf:protoc-gen-go`                       | [link]() |
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
| --------------------------------------------- | -------- |
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

# History

The original rules_proto was <https://github.com/pubref/rules_proto>.  This was
redesigned around the `proto_library` rule and moved to
<https://github.com/stackb/rules_proto>.  

An initial aspect implementation was prototyped and considered, but not merged.
As a result, the ruleset was forked to
<https://github.com/rules-proto-grpc/rules_proto_grpc> and primarily used
aspect-based compilation for quite a while before that was again deprecated and
removed.

Maintaining `stackb/rules_proto` and its polyglot set of languages in its
original v0-v1 form became a full-time (unpaid) job.  The main issue is that the
`{LANG}_{PROTO|GRPC}_library` rules are **tightly bound to a specific set of
dependencies**.  As such, rules_proto users are tightly bound to the specific
labels named by those rules.  This is a problem for the maintainer as one must
keep the dependencies current.  It is also a problem for rule consumers: *it
becomes increasingly difficult to upgrade as the dependencies as assumptions and
dependencies drift*.

With `stackb/rules_proto` in its `v2` gazelle-based form, it is largely
**dependency-free**: other than gazelle and the `protoc` toolchain, *there are no
dependencies that you cannot fully control in your own workspace via the gazelle
configuration*.

The gazelle based design also makes things *much* simpler and powerful, because
the **content of the proto files is the source of truth**.  Due to the fact that
Bazel does not permit reading/interpreting a file during the scope of an action,
it is impossible to make a decision about what to do.  A prime example of this
is the `go_package` option.  If the `go_package` option is present, the location
of the output file for `protoc-gen-go` is completely different.  As a result,
the information about the go_package metadata ultimately needs to be duplicated
so that the build system can know about it.

The gazelle-based approach moves all that messy interpretation and evaluation
into a precompiled state; as a result, the work that needs to be done in the
action itself is dramatically simplified.

Furthermore, by parsing the proto files it is easy to support complex custom
plugins that do things like:

- Emit no files (only assert/lint).
- Emit a file only if a specific enum constant is found.  With the previous
  design, this was near impossible.  With the `v2` design, the `protoc.Plugin`
  implementation can trivially perform that evaluation because it is handed the
  complete proto AST during gazelle evaluation.

