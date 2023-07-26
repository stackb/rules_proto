# `rules_proto (v2)`

> **Warning**
> This is *not* a true fork of stackb/rules_proto -- depend on that repo, not
> this one. Changes here should flow upstream.

[![Build status](https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master)](https://buildkite.com/bazel/stackb-rules-proto)
[![Go Reference](https://pkg.go.dev/badge/github.com/stackb/rules_proto.svg)](https://pkg.go.dev/github.com/stackb/rules_proto)

Bazel starlark rules for building protocol buffers +/- gRPC :sparkles:.

<table border="0">
  <tr>
    <td><img src="https://upload.wikimedia.org/wikipedia/en/thumb/7/7d/Bazel_logo.svg/1920px-Bazel_logo.svg.png" height="120"/></td>
    <td><img src="https://user-images.githubusercontent.com/50580/141892423-5205bbfd-8487-442b-81c7-f56fa3d1f69e.jpeg" height="120"/></td>
    <td><img src="https://user-images.githubusercontent.com/50580/141900696-bfb2d42d-5d2c-46f8-bd9f-06515969f6a2.png" height="120"/></td>
    <td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" height="120"/></td>
  </tr>
  <tr>
    <td>bazel</td>
    <td>gazelle</td>
    <td>protobuf</td>
    <td>grpc</td>
  </tr>
</table>

`@build_stack_rules_proto` provides:

1. Rules for driving the `protoc` tool within a bazel workspace.
2. A [gazelle](https://github.com/bazelbuild/bazel-gazelle/) extension that
   generates rules based on the content of your `.proto` files.
3. A repository rule that runs gazelle in an external workspace.
4. Example setups for a variety of languages.

# Table of Contents

- [Getting Started](#getting-started)
  - [Workspace Boilerplate](#workspace-boilerplate)
  - [Gazelle Setup](#gazelle-setup)
  - [Gazelle Configuration](#gazelle-configuration)
    - [Build Directives](#build-directives)
    - [YAML Config File](#yaml-configuration)
  - [Running Gazelle](#running-gazelle)
- [Build Rules](#build-rules)
  - [proto_compile](#proto_compile)
  - [proto_plugin](#proto_plugin)
  - [proto_compiled_sources](#proto_compiled_sources)
  - [Deep dive on the mappings attribute](#the-output_mappings-attribute)
- [Repository Rules](#repository-rules)
  - [proto_repository](#proto_repository)
  - [proto_gazelle](#proto_gazelle)
- [Plugin Implementations](#plugin-implementations)
- [Rule Implementations](#rule-implementations)
- [Writing Custom Plugins & Rules](#writing-custom-plugins-and-rules)
- [History of this repository](#history)

# Getting Started

## `WORKSPACE` Boilerplate

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Release: v2.0.1
# TargetCommitish: master
# Date: 2022-10-20 02:38:27 +0000 UTC
# URL: https://github.com/stackb/rules_proto/releases/tag/v2.0.1
# Size: 2071295 (2.1 MB)
http_archive(
    name = "build_stack_rules_proto",
    sha256 = "ac7e2966a78660e83e1ba84a06db6eda9a7659a841b6a7fd93028cd8757afbfb",
    strip_prefix = "rules_proto-2.0.1",
    urls = ["https://github.com/stackb/rules_proto/archive/v2.0.1.tar.gz"],
)
```

```python
register_toolchains("@build_stack_rules_proto//toolchain:standard")
```

> This prepares `protoc` for the `proto_compile` rule. For simple setups,
> consider `@build_stack_rules_proto//toolchain:prebuilt` to skip compilation of
> the tool.

> **NOTE**: if you are planning on hand-writing your `BUILD.bazel` rules
> yourself (not using the gazelle build file generator), **STOP HERE**. You'll
> need to provide typical proto dependencies such as `@rules_proto` and
> `@com_google_protobuf` (use macros below if desired), but no additional core
> dependencies are needed at this point.

---

```python
load("@build_stack_rules_proto//deps:core_deps.bzl", "core_deps")

core_deps()
```

> This brings in `@io_bazel_rules_go`, `@bazel_gazelle`, and `@rules_proto` if
> you don't already have them.

> The gazelle extension and associated golang dependencies are optional; you can
> write `proto_compile` and other derived rules by hand. For gazelle support,
> carry on.

---

```python
load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(version = "1.18.2")
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
load("@build_stack_rules_proto//:go_deps.bzl", "gazelle_protobuf_extension_go_deps")

gazelle_protobuf_extension_go_deps()
```

> This brings in `@com_github_emicklei_proto`.
> [github.com/emicklei/proto](https://github.com/emicklei/proto) is used by the
> gazelle extension to parse proto files.

---

```python
load("@build_stack_rules_proto//deps:protobuf_core_deps.bzl", "protobuf_core_deps")

protobuf_core_deps()
```

> This brings in `@com_google_protobuf` and friends if you don't already have
> them.

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

The gazelle extension can be configured using "build directives" and/or a YAML
file.

### Build Directives

Gazelle is configured by special comments in BUILD files called _directives_.

> Gazelle works by visiting each package in your workspace; configuration is
> done "on the way in" whereas actual rule generation is done "on the way out".
> Gazelle configuration of a subdirectory inherits that from its parents. As
> such, directives placed in the root `BUILD.bazel` file apply to the entire
> workspace.

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
  your own). The core responsibility a `protoc.Plugin` implementation is to to
  _predict_ the files that a protoc plugin tool will generate for an individual
  `proto_library` rule. The implemention has full read access to the
  `protoc.File`s in the `proto_library` to be able to predict _if_ a file will
  be generated and _where_ it will appear in the filesystem (specifically,
  relative to the bazel execution root during a `proto_compile` action).
- `gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile`
  associates the name `proto_compile` with a piece of golang code that
  implements the `protoc.LanguageRule` interface. The extension maintains a
  registry of rule implementations. Similarly, the extension ships with a bunch
  of them out of the box, but you can write your own custom rules as well. The
  core responsibility a `protoc.LanguageRule` implementation is construct a
  gazelle `rule.Rule` based upon a `proto_library` rule and the set of plugins
  that are configured with it.
- `gazelle:proto_language cpp plugin cpp` instantiates a `protoc.LanguageConfig`
  having the name `cpp` and adds the `cpp` plugin to it. The language
  configuration bundles bundles plugins and rules together.
- `gazelle:proto_rule grpc_cc_library deps @com_github_grpc_grpc//:grpc++`
  configures the rule such that all generated rules will have that dependency.

> **+/- intent modifiers**. Although not pictured in this example, many of the
> directives take an _intent modifier_ to turn configuration on/off. For
> example, if you wanted to suppress the grpc c++ plugin in the package
> `//proto/javaapi`, put a directive like
> `gazelle:proto_language cpp rule -grpc_cc_library` in
> `proto/javaapi/BUILD.bazel` (note the `-` symbol preceding the name). To
> suppress the language entirely, use
> `gazelle:proto_language cpp enabled false`.

### YAML Configuration

You can also configure the extension using a YAML file. This is semantically
similar to adding gazelle directives to the root `BUILD` file; the YAML
configuration applies to all downstream packages. The equivalent YAML config for
the above directives looks like:

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

> A yaml config is particularly useful in conjunction with the
> `proto_repository` rule, for example to apply a set of custom plugins over the
> googleapis/googleapis repo.

To use this in a gazelle rule, specify `-proto_configs` in `args`
(comma-separated list):

```python
gazelle(
    name = "gazelle",
    gazelle = ":gazelle-protobuf",
    args = [
        "-proto_configs=example/config.yaml",
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
These are nearly always very thin wrappers for the "real" rule. For example,
here's the implementation in `proto_cc_library.bzl`:

```python
load("@rules_cc//cc:defs.bzl", "cc_library")

def proto_cc_library(**kwargs):
    cc_library(**kwargs)
```

An implementation detail of gazelle itself is that two different language
extensions should not _claim_ the same load namespace, so in order to prevent
potential conflicts with other possible gazelle extensions, using the name
`@rules_cc//cc:defs.bzl%cc_library` is undesirable.

## Build Rules

The core of `stackb/rules_proto` contains two build rules:

| Rule            | Description                                             |
|-----------------|---------------------------------------------------------|
| `proto_compile` | Executes the `protoc` tool.                             |
| `proto_plugin`  | Provides static `protoc` plugin-specific configuration. |

### proto_compile

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
- A `proto_compile` rule references a single `proto_library` target.
- The `plugins` attribute is a list of labels to `proto_plugin` targets.
- The `outputs` attribute names the files that will be generated by the protoc
  invocation.
- The
  [proto](https://github.com/bazelbuild/bazel-gazelle/blob/master/language/proto/lang.go)
  extension provided by [bazel-gazelle] is responsible for generating
  `proto_library`.

### proto_plugin

`proto_plugin` primarily provides the plugin tool executable. The example seen
above is the simplest case where the plugin is builtin to `protoc` itself; no
separate plugin tool is required. In this case the `proto_plugin` rule
degenerates into just a `name`.

It is possible to add additional plugin-specific
`name = "foo", options = ["bar"]` on the `proto_plugin` rule, but the use-case
for this is narrow. Generally it is preferred to say
`# gazelle:proto_plugin foo option bar` such that the option can be interpreted
during a gazelle run.

### proto_compiled_sources

`proto_compiled_sources` is used when you prefer to check the generated files
into source control. This may be necessary for legacy reasons, during an initial
Bazel migration, or to support better IDE integration.

The shape of a `proto_compiled_sources` rule is essentially identical to
`proto_compile` with one exception: generated source are named in the `srcs`
attribute rather than `outputs`.

For example, a `proto_compiled_sources` named `//example/thing:proto_go_sources`
is a macro that generates three rules:

1. `bazel build //example/thing:proto_go_sources` emits the generated files.
2. `bazel run //example/thing:proto_go_sources.update` copies the generated
   files back into the source package.
3. `bazel test //example/thing:proto_go_sources_test` asserts the source files
   are identical to generated files.

In this scenario, `2.` is used to build the generated files (in the `bazel-bin/`
output tree) and copy the `example/thing/thing.pb.go` back into place where it
will be committed under source control. `3.` is used to prevent drift: if a
developer modifies `thing.proto` and neglects to run the `.update` the test will
fail in CI.

### proto_compile_assets

The macro `proto_compile_assets` aggregates a list of dependencies (which
provide `ProtoCompileInfo`) into a single runnable target that copies files in
bulk.

For example, `bazel run //proto:assets` will copy all the generated `.pb.go`
files back into the source tree:

```py
load("@build_stack_rules_proto//rules:proto_compile_assets.bzl", "proto_compile_assets")

proto_compile_assets(
    name = "assets",
    deps = [,
      "//proto/api/v1:proto_go_compile",
      "//proto/api/v2:proto_go_compile",
      "//proto/api/v3:proto_go_compile",
    ],
)
```

### The `output_mappings` attribute

Consider the following rule within the package `example/thing`:

```python
proto_compile(
    name = "thing_go_compile",
    output_mappings = ["thing.pb.go=github.com/stackb/rules_proto/example/thing/thing.pb.go"],
    outputs = ["thing.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go"],
    proto = "thing_proto",
)
```

This rule is declaring that a file `bazel-bin/example/thing/thing.pb.go` will be
output when the action is run. When we
`bazel build //example/thing:thing_go_compile`, the file is indeed created.

Let's temporarily comment out the `output_mappings` attribute and rebuild:

```python
proto_compile(
    name = "thing_go_compile",
    # output_mappings = ["thing.pb.go=github.com/stackb/rules_proto/example/thing/thing.pb.go"],
    outputs = ["thing.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go"],
    proto = "thing_proto",
)
```

```sh
$ bazel build //example/thing:thing_go_compile
ERROR: /github.com/stackb/rules_proto/example/thing/BUILD.bazel:54:14: output 'example/thing/thing.pb.go' was not created
```

What happened? Let's add a debugging attribute `verbose = True` on the rule:
this will print debugging information and show the bazel sandbox before and
after the `protoc` tool is invoked:

```python
proto_compile(
    name = "thing_go_compile",
    # output_mappings = ["thing.pb.go=github.com/stackb/rules_proto/example/thing/thing.pb.go"],
    outputs = ["thing.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go"],
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

So, the file was created, but not in the location we wanted. In this case the
`protoc-gen-go` plugin is not "playing nice" with Bazel. Because this
`thing.proto` has
`option go_package = "github.com/stackb/rules_proto/example/thing;thing";`, the
output location is no longer based on the `package`. This is a problem, because
Bazel semantics disallow declaring a File outside its package boundary. As a
result, we need to do a
`mv ./bazel-out/darwin-fastbuild/bin/github.com/stackb/rules_proto/example/thing/thing.pb.go ./bazel-out/darwin-fastbuild/bin/example/thing/thing.pb.go`
to relocate the file into its expected location before the action terminates.

Therefore, the `output_mappings` attribute is a list of entries that map file
locations `want=got` relative to the action execution root. It is required when
the actual output location does not match the desired location. This can occur
if the proto `package` statement does not match the Bazel package path, or in
special circumstances specific to the plugin itself (like `go_package`).

## Repository Rules

## proto_repository

From an implementation standpoint, this is very similar to the `go_repository`
rule. Both can download external files and then run the gazelle generator over
the downloaded files. Example:

```python
proto_repository(
    name = "googleapis",
    build_directives = [
        "gazelle:resolve proto google/api/http.proto //google/api:http_proto",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    override_go_googleapis = True,
    proto_language_config_file = "//example:config.yaml",
    strip_prefix = "googleapis-02710fa0ea5312d79d7fb986c9c9823fb41049a9",
    type = "zip",
    urls = ["https://codeload.github.com/googleapis/googleapis/zip/02710fa0ea5312d79d7fb986c9c9823fb41049a9"],
)
```

Takeaways:

- The `urls`, `strip_prefix` and `type` behave similarly to `http_archive`.
- `build_file_proto_mode` is the same the `go_repository` attribute of the same
  name; additionally the value `file` is permitted which generates a
  `proto_library` for each individual proto file.
- `build_file_generation` is the same the `go_repository` attribute of the same
  name; additionally the value `clean` is supported. For example, googleapis
  already has a set of BUILD files; the `clean` mode will remove all the
  existing build files before generating new ones.
- `build_directives` is the same as `go_repository`. Resolve directives in this
  case are needed because the gazelle `language/proto` extension hardcodes a
  proto import like `google/api/http.proto` to resolve to the `@go_googleapis`
  workspace; here we want to make sure that http.proto resolves to the same
  external workspace.
- `proto_language_config_file` is an optional label pointing to a valid
  `config.yaml` file to configure this extension.
- `override_go_googleapis` is a boolean attribute that has special meaning for
  the googleapis repository. Due to the fact that the builtin gazelle "proto"
  extension has
  [hardcoded logic](https://github.com/bazelbuild/bazel-gazelle/blob/master/language/proto/known_proto_imports.go)
  for what googleapis deps look like, additional work is needed to override
  that. With this sample configuration, the following build command succeeds:

```bash
$ bazel build @googleapis//google/api:annotations_cc_library
Target @googleapis//google/api:annotations_cc_library up-to-date:
  bazel-bin/external/googleapis/google/api/libannotations_cc_library.a
  bazel-bin/external/googleapis/google/api/libannotations_cc_library.so
```

Another example using the Bazel repository:

```python
load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository")

proto_repository(
    name = "bazelapis",
    build_directives = [
        "gazelle:exclude third_party",
        "gazelle:proto_language go enable true",
        "gazelle:proto_language closure enabled true",
        "gazelle:prefix github.com/bazelbuild/bazelapis",
    ],
    build_file_expunge = True,
    build_file_proto_mode = "file",
    cfgs = ["//proto:config.yaml"],
    imports = [
        "@googleapis//:imports.csv",
        "@protoapis//:imports.csv",
        "@remoteapis//:imports.csv",
    ],
    strip_prefix = "bazel-02ad3e3bc6970db11fe80f966da5707a6c389fdd",
    type = "zip",
    urls = ["https://codeload.github.com/bazelbuild/bazel/zip/02ad3e3bc6970db11fe80f966da5707a6c389fdd"],
)
```

Takeaways:

- The `build_directives` are setting the `gazelle:prefix` for the `language/go`
  plugin; two `proto_language` configs named in the `proto/config.yaml` are
  being enabled.
- `build_file_expunge` means _remove all existing BUILD files before generating
  new ones_.
- `build_file_proto_mode = "file"` means _make a separate `proto_library` rule
  for every `.proto` file_.
- `cfgs = ["//proto:config.yaml"]` means _use the configuration in this YAML
  file as a base set of gazelle directives_. It is easier/more consistent to
  share the same `config.yaml` file across multiple `proto_repository` rules.

The last one `imports = ["@googleapis//:imports.csv", ...]` requires extra
explanation. When the `proto_repository` gazelle process finishes, it writes a
file named `imports.csv` in the root of the external workspace. This file
records the association between import statements and bazel labels, much like a
`gazelle:resolve` directive:

```csv
# GENERATED FILE, DO NOT EDIT (created by gazelle)
# lang,imp.lang,imp,label
go,go,github.com/googleapis/gapic-showcase/server/genproto,@googleapis//google/example/showcase/v1:compliance_go_proto
go,go,google.golang.org/genproto/firestore/bundle,@googleapis//google/firestore/bundle:bundle_go_proto
go,go,google.golang.org/genproto/googleapis/actions/sdk/v2,@googleapis//google/actions/sdk/v2:account_linking_go_proto
```

Therefore, the `imports` attribute assists gazelle in figuring how to resolve
imports. Therefore, when gazelle is preparing a `go_library` rule and finds a
`main.go` file having an import on
`google.golang.org/genproto/firestore/bundle`, it knows to put
`@googleapis//google/firestore/bundle:bundle_go_proto` in the rule `deps`.

To take advantage of this mechanism in the default workspace, use the
`proto_gazelle` rule.

## proto_gazelle

`proto_gazelle` is not a repository rule: it's just like the typical `gazelle`
rule, but with extra deps resolution superpowers. But, we discuss it here since
it works in conjunction with `proto_repository`:

```python
load("@build_stack_rules_proto//rules:proto_gazelle.bzl", "DEFAULT_LANGUAGES", "proto_gazelle")

proto_gazelle(
    name = "gazelle",
    cfgs = ["//proto:config.yaml"],
    command = "update",
    gazelle = ":gazelle-protobuf",
    imports = [
        "@bazelapis//:imports.csv",
        "@googleapis//:imports.csv",
        "@protoapis//:imports.csv",
        "@remoteapis//:imports.csv",
    ],
)
```

In this example, we are again setting the base gazelle config using the YAML
file (the same one used in for the `proto_repository` rules). We are also now
importing resolve information from four external sources.

With this setup, we can simply place an import statement like
`import "src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream.proto";`
in a `foo.proto` file in the default workspace, and gazelle will automagically
figure out the import dependency tree spanning `@bazelapis`, `@remoteapis`,
`@googleapis`, and the well-known types from `@protoapis`.

This works for any `proto_language`, with any set of custom protoc plugins.

## golden_filegroup

`golden_filegroup` is a utility macro for golden file testing. It works like a
native filegroup, but adds `.update` and `.test` targets. Example:

```py
load("@build_stack_rules_proto//rules:golden_filegroup.bzl", "golden_filegroup")

# golden_filegroup asserts that generated files named in 'srcs' are
# identical to the ones checked into source control.
#
# Usage:
#
# $ bazel build :golden        # not particularly useful, just a regular filegroup
#
# $ bazel test  :golden.test   # checks that generated files are identical to
# ones in git (for CI)
#
# $ bazel run   :golden.update # copies the generated files into source tree
# (then 'git add' to your PR if it looks good)
golden_filegroup(
    name = "golden",
    srcs = [
        ":some_generated_file1.json",
        ":some_generated_file2.json",
    ],
)
```

## Plugin Implementations

The plugin name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/plugin_name convention.

| Plugin                                                                                                                 |
|------------------------------------------------------------------------------------------------------------------------|
| [builtin:cpp](pkg/plugin/builtin/cpp_plugin.go)                                                                        |
| [builtin:csharp](pkg/plugin/builtin/csharp_plugin.go)                                                                  |
| [builtin:java](pkg/plugin/builtin/java_plugin.go)                                                                      |
| [builtin:js:closure](pkg/plugin/builtin/js_closure_plugin.go)                                                          |
| [builtin:js:common](pkg/plugin/builtin/js_common_plugin.go)                                                            |
| [builtin:objc](pkg/plugin/builtin/objc_plugin.go)                                                                      |
| [builtin:php](pkg/plugin/builtin/php_plugin.go)                                                                        |
| [builtin:python](pkg/plugin/builtin/python_plugin.go)                                                                  |
| [builtin:pyi](pkg/plugin/builtin/pyi_plugin.go)                                                                        |
| [builtin:ruby](pkg/plugin/builtin/ruby_plugin.go)                                                                      |
| [grpc:grpc:cpp](pkg/plugin/builtin/grpc_grpc_cpp.go)                                                                   |
| [grpc:grpc:protoc-gen-grpc-python](pkg/plugin/grpc/grpc/protoc-gen-grpc-python.go)                                     |
| [golang:protobuf:protoc-gen-go](pkg/plugin/golang/protobuf/protoc-gen-go.go)                                           |
| [grpc:grpc-go:protoc-gen-go-grpc](pkg/plugin/grpc/grpcgo/protoc-gen-go-grpc.go)                                        |
| [grpc:grpc-java:protoc-gen-grpc-java](pkg/plugin/grpc/grpcjava/protoc-gen-grpc-java.go)                                |
| [grpc:grpc-node:protoc-gen-grpc-node](pkg/plugin/grpc/grpcnode/protoc-gen-grpc-node.go)                                |
| [grpc:grpc-web:protoc-gen-grpc-web](pkg/plugin/grpc/grpcweb/protoc-gen-grpc-web.go)                                    |
| [gogo:protobuf:protoc-gen-combo](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                          |
| [gogo:protobuf:protoc-gen-gogo](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                           |
| [gogo:protobuf:protoc-gen-gogofast](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                       |
| [gogo:protobuf:protoc-gen-gogofaster](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                     |
| [gogo:protobuf:protoc-gen-gogoslick](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                      |
| [gogo:protobuf:protoc-gen-gogotypes](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                      |
| [gogo:protobuf:protoc-gen-gostring](pkg/plugin/gogo/protobuf/protoc-gen-gogo.go)                                       |
| [grpc-ecosystem:grpc-gateway:protoc-gen-grpc-gateway](pkg/plugin/grpcecosystem/grpcgateway/protoc-gen-grpc-gateway.go) |
| [scalapb:scalapb:protoc-gen-scala](pkg/plugin/scalapb/scalapb/protoc_gen_scala.go)                                     |
| [stackb:grpc.js:protoc-gen-grpc-js](pkg/plugin/stackb/grpc_js/protoc-gen-grpc-js.go)                                   |
| [stephenh:ts-proto:protoc-gen-ts-proto](pkg/plugin/stephenh/ts-proto/protoc-gen-ts-proto.go)                           |

## Rule Implementations

The rule name is an opaque string, but by convention they are maven-esqe
artifact identifiers that follow a GitHub org/repo/rule_name convention.

| Plugin                                                                                            |
|---------------------------------------------------------------------------------------------------|
| [stackb:rules_proto:grpc_cc_library](pkg/rule/rules_cc/grpc_cc_library.go)                        |
| [stackb:rules_proto:grpc_closure_js_library](pkg/rule/rules_closure/grpc_closure_js_library.go)   |
| [stackb:rules_proto:grpc_java_library](pkg/rule/rules_java/grpc_java_library.go)                  |
| [stackb:rules_proto:grpc_nodejs_library](pkg/rule/rules_nodejs/grpc_nodejs_library.go)            |
| [stackb:rules_proto:grpc_web_js_library](pkg/rule/rules_nodejs/grpc_web_js_library.go)            |
| [stackb:rules_proto:grpc_py_library](pkg/rule/rules_python/grpc_py_library.go)                    |
| [stackb:rules_proto:proto_cc_library](pkg/rule/rules_cc/proto_cc_library.go)                      |
| [stackb:rules_proto:proto_closure_js_library](pkg/rule/rules_closure/proto_closure_js_library.go) |
| [stackb:rules_proto:proto_compile](pkg/protoc/proto_compile.go)                                   |
| [stackb:rules_proto:proto_compiled_sources](pkg/protoc/proto_compiled_sources.go)                 |
| [stackb:rules_proto:proto_descriptor_set](pkg/protoc/proto_descriptor_set.go)                     |
| [stackb:rules_proto:proto_go_library](pkg/rule/rules_go/go_library.go)                            |
| [stackb:rules_proto:proto_java_library](pkg/rule/rules_java/proto_java_library.go)                |
| [stackb:rules_proto:proto_nodejs_library](pkg/rule/rules_nodejs/proto_nodejs_library.go)          |
| [stackb:rules_proto:proto_py_library](pkg/rule/rules_python/proto_py_library.go)                  |
| [bazelbuild:rules_scala:scala_proto_library](pkg/rule/rules_scala/scala_proto_library.go)         |

Please consult the `example/` directory and unit tests for more additional
detail.


# Writing Custom Plugins and Rules

Custom plugin implementations and rule implementations can be written in golang
or starlark.  Golang implementations are statically compiled into the final
`gazelle_binary` whereas starlark plugins are evaluated at gazelle runtime.

## +/- of golang implementations

- `+` Full power of a statically compiled language, the golang stdlib, and
  external dependencies.
- `+` Easier to test.
- `+` API not experimental.
- `-` Cannot be used in a `proto_repository` rule without forking
  stackb/rules_proto.
- `-` Initial setup harder, often housed within your own custom gazelle
  extension.

Until a dedicated tutorial is available, please consult the source code for
examples.

## +/- of starlark implementations

- `+` More familiar to developer with starlark experience but not golang.
- `+` Easier setup (*.star files in your gazelle repository)
- `+` Possible to use in conjunction with the `proto_repository` rule.
- `-` Limited API: can only reference state that has been already configured via gazelle directives.
- `-` Not possible to implement stateful design.
- `-` No standard library.

Until a dedicated tutorial is available, please consult the reference example in
`example/testdata/starlark_java`.

# History

The original rules_proto was <https://github.com/pubref/rules_proto>. This was
redesigned around the `proto_library` rule and moved to
<https://github.com/stackb/rules_proto>.

Following earlier experiments with aspects, this ruleset was forked to
<https://github.com/rules-proto-grpc/rules_proto_grpc>. Aspect-based compilation
was featured for quite a while there but has recently been deprecated.

Maintaining `stackb/rules_proto` and its polyglot set of languages in its
original v0-v1 form became a full-time (unpaid) job. The main issue is that the
`{LANG}_{PROTO|GRPC}_library` rules are **tightly bound to a specific set of
dependencies**. As such, rules_proto users are tightly bound to the specific
labels named by those rules. This is a problem for the maintainer as one must
keep the dependencies current. It is also a problem for rule consumers: _it
becomes increasingly difficult to upgrade as the dependencies as assumptions and
dependencies drift_.

With `stackb/rules_proto` in its `v2` gazelle-based form, it is largely
**dependency-free**: other than gazelle and the `protoc` toolchain, _there are
no dependencies that you cannot fully control in your own workspace via the
gazelle configuration_.

The gazelle based design also makes things _much_ simpler and powerful, because
the **content of the proto files is the source of truth**. Due to the fact that
Bazel does not permit reading/interpreting a file during the scope of an action,
it is impossible to make a decision about what to do. A prime example of this is
the `go_package` option. If the `go_package` option is present, the location of
the output file for `protoc-gen-go` is completely different. As a result, the
information about the go_package metadata ultimately needs to be duplicated so
that the build system can know about it.

The gazelle-based approach moves all that messy interpretation and evaluation
into a precompiled state; as a result, the work that needs to be done in the
action itself is dramatically simplified.

Furthermore, by parsing the proto files it is easy to support complex custom
plugins that do things like:

- Emit no files (only assert/lint).
- Emit a file only if a specific enum constant is found. With the previous
  design, this was near impossible. With the `v2` design, the `protoc.Plugin`
  implementation can trivially perform that evaluation because it is handed the
  complete proto AST during gazelle evaluation.
