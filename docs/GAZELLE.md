# `protobuf` gazelle extension

From a high-level, the `protobuf` extension works as follows:

- **Flags Phase**: Registers command-line flags for the extension including
  `proto_configs` (YAML config files), `proto_imports_in/out` (import index
  files), `proto_rule` (custom Starlark rules), and `proto_plugin` (custom
  Starlark plugins).  Most of these flags are abstracted away by the
  `proto_gazelle` rule.
- **Config Phase**: Processes configuration directives from BUILD files and
  applies settings. Handles `gazelle:proto_language`, `gazelle:proto_plugin`,
  and `gazelle:proto_rule` directives. Loads YAML config files and import index
  files specified in flags.
- **Generate Phase**: Reads `proto_library` rules, parses `.proto` files named
  in `proto_library.srcs` and language-specific rules based on configured
  languages and plugins. Generates `proto_compile` and language specific rules
  as per configuration.  Records proto file imports and dependencies for later
  resolution.
- **Resolve Phase**: Translates imports into Bazel dependencies.

## Setup

### Step 1: create a `rules_proto_config.yaml` file

Configuration of the extension is expressed in YAML and/or gazelle directives of
three core concepts:

- `gazelle:proto_plugin`: metadata about a proto plugin.
- `gazelle:proto_rule`: metadata about a rule that generates Bazel targets
- `gazelle:proto_language`: composes plugins and rules together.

For example, the following are equivalent:

```yaml
plugins:
  - name: protoc-gen-go
    label: "@//proto/plugin:protoc-gen-go"
    implementation: golang:protobuf:protoc-gen-go
    deps:
      - "@org_golang_google_protobuf//reflect/protoreflect"
      - "@org_golang_google_protobuf//runtime/protoimpl"
  - name: protoc-gen-go-grpc
    label: "@//proto/plugin:protoc-gen-go-grpc"
    implementation: grpc:grpc-go:protoc-gen-go-grpc
    deps:
      - "@org_golang_google_grpc//:go_default_library"
      - "@org_golang_google_grpc//codes"
      - "@org_golang_google_grpc//status"
rules:
  - name: proto_compile
    implementation: stackb:rules_proto:proto_compile
  - name: proto_compiled_sources
    implementation: stackb:rules_proto:proto_compiles_sources
  - name: proto_go_library
    implementation: stackb:rules_proto:proto_go_library
languages:
  - name: go
    plugins:
      - protoc-gen-go
      - protoc-gen-go-grpc
    rules:
      - proto_compile
      - proto_go_library
    enabled: false
```

```py
# gazelle:proto_plugin protoc-gen-go implementation golang:protobuf:protoc-gen-go
# gazelle:proto_plugin protoc-gen-go label @//proto/plugin:protoc-gen-go
# gazelle:proto_plugin protoc-gen-go deps @org_golang_google_protobuf//reflect/protoreflect
# gazelle:proto_plugin protoc-gen-go deps @org_golang_google_protobuf//runtime/protoimpl
# gazelle:proto_plugin protoc-gen-go-grpc implementation golang:protobuf:protoc-gen-go-grpc
# gazelle:proto_plugin protoc-gen-go label @//proto/plugin:protoc-gen-go-grpc
# gazelle:proto_plugin protoc-gen-go-grpc deps @org_golang_google_grpc//:go_default_library
# gazelle:proto_plugin protoc-gen-go-grpc deps @org_golang_google_grpc//codes
# gazelle:proto_plugin protoc-gen-go-grpc deps @org_golang_google_grpc//status

# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_rule proto_compile visibility //visibility:public
# gazelle:proto_rule proto_compiled_source implementation stackb:rules_proto:proto_compiled_source
# gazelle:proto_rule proto_compiled_source visibility //visibility:public
# gazelle:proto_rule proto_go_library implementation stackb:rules_proto:proto_go_library
# gazelle:proto_rule proto_go_library visibility //visibility:public

# gazelle:proto_language go plugin protoc-gen-go
# gazelle:proto_language go plugin protoc-gen-go-grpc
# gazelle:proto_language go rule proto_compile
# gazelle:proto_language go rule proto_go_library
```

The YAML form is typically used to assemble/compose a base set of
proto_languages that can then be "turned on" (enabled) and tweaked using build
directives.

So while a YAML file is not strictly required, it is more ergonomic and
generally recommended.

```py
exports_files(["rules_proto_config.yaml"])
```

> The file should be exported as it will be read by a `repository_ctx`:

### Step 2 (optional): custom `proto_plugin` rules

It is not strictly necessary to configure your own `proto_plugin` rules as the
individual plugin implementations include a default.

For example, the `golang:protobuf:protoc-gen-go` implementation defaults to
`@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go`:

```go
func (p *ProtocGenGoPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
    ...
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/golang/protobuf", "protoc-gen-go"),
		Outputs: p.outputs(ctx.ProtoLibrary, mappings),
		Options: ctx.PluginConfig.GetOptions(),
	}
}
```

```py
# ../external/build_stack_rules_proto+/plugin/golang/protobuf/BUILD.bazel:3:13
proto_plugin(
  name = "protoc-gen-go",
  tool = "@org_golang_google_protobuf//cmd/protoc-gen-go:protoc-gen-go",
  visibility = ["//visibility:public"],
)
```

In the current example, `@//proto/plugin:protoc-gen-go` (in the default
workspace) is the target `proto_plugin` rule.  An additional option has been
included for demonstrative purposes.

```py
proto_plugin(
  name = "protoc-gen-go",
  options = [
    "experimental_strip_nonfunctional_codegen=true",
  ],
  tool = "@org_golang_google_protobuf//cmd/protoc-gen-go:protoc-gen-go",
  visibility = ["//visibility:public"],
)
```

### Step 3: configure `gazelle_binary` to include `@build_stack_rules_proto//language/protobuf`

```python
load("@bazel_gazelle//:def.bzl", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle-protobuf",
    languages = [
        "@bazel_gazelle//language/go",
        "@bazel_gazelle//language/proto",
        "@build_stack_rules_proto//language/protobuf",
    ],
)

gazelle(
    name = "gazelle",
    gazelle = ":gazelle-protobuf",
)
```

> Order matters.  Since the `protobuf` inspects `proto_library` rules generated
> by `proto` extension, that one should precede it. 

### Step 2: configure the `gazelle` or `proto_gazelle` rule with the base YAML configuration

> NOTE: the `proto_gazelle` rule is very similar to the `gazelle` rule, with a
> few additional attributes.

```py
load("@build_stack_rules_proto//rules:proto_gazelle.bzl", "proto_gazelle")

proto_gazelle(
    name = "gazelle",
    cfgs = ["//:rules_proto_config.yaml"],
    gazelle = ":gazelle-protobuf",
    imports = [
        "@bazelapis//:imports.csv",
        "@googleapis//:imports.csv",
        "@protoapis//:imports.csv",
        "@remoteapis//:imports.csv",
    ],
)
```

- `gazelle`: an `attr.label` that should point to a `gazelle_binary` rule.
- `cfgs`: an optional `attr.label_list` that references base YAML configuration files.
- `imports`: an optional `attr.label_list` that references CSV files that are
  produced by the `proto_repository` rule.
  - The format is CSV, but you can think of this as a list of `gazelle:resolve
    <lang> <import> <label>` directives.
  - For example, after parsing `@protoapis//:imports.csv`, the in-memory data
    store for gazelle `resolve` extension will contain entries like
  `gazelle:resolve proto google/protobuf/any.proto
  @protoapis//google/proto:any_proto`.

### Step 3: tweak base configuration

In the current example, your root `BUILD.bazel` file might contain:

```py
# -- Gazelle language "walk" ---
# gazelle:exclude vendor

# -- Gazelle language "go" ---
# gazelle:prefix github.com/stackb/bazel-aquery-differ
# gazelle:go_generate_proto false

# -- Gazelle language "protobuf" ---
# gazelle:proto_language go enable true
```

### Step 4: `bazel run gazelle`

Following a gazelle invocation, the following rules are generated in
`example/person:BUILD.bazel` having a file `example/person:person.proto`:

```py
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@build_stack_rules_proto//rules/go:proto_go_library.bzl", "proto_go_library")
load("@rules_proto//proto:defs.bzl", "proto_library")

proto_library(
    name = "person_proto",
    srcs = ["person.proto"],
    visibility = ["//visibility:public"],
    deps = ["//example/place:place_proto"],
)

proto_compile(
    name = "person_go_compile",
    output_mappings = ["person.pb.go=github.com/stackb/rules_proto/example/person/person.pb.go"],
    outputs = ["person.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go"],
    proto = "person_proto",
    visibility = ["//visibility:public"],
)

proto_go_library(
    name = "person_go_proto",
    srcs = ["person.pb.go"],
    importpath = "github.com/stackb/rules_proto/example/person",
    visibility = ["//visibility:public"],
    deps = [
        "//example/place:place_go_proto",
        "@org_golang_google_protobuf//proto",
        "@org_golang_google_protobuf//reflect/protoreflect",
        "@org_golang_google_protobuf//runtime/protoimpl",
    ],
)
```

### Default Workspace Vendoring Strategy: `proto_compiled_sources`

For IDE completion and UX, you may decide that it is preferable to vendor the
generated source files into your repo.  This is the use case for
`proto_compiled_sources`, which resembles `proto_compile` but has additional
implicit targets `.update` (copies the `.pb.go` files into the source tree) and
`_test` (asserts that the vendored copy is up-to-date).

You **could** change the `language.rules` section of the YAML file to modify the
base configuration, but in this example we will use "intents" to change the base
configuration instead.  The following removes the default rules and adds
`proto_compiled_sources` instead.

```py
# gazelle:proto_language go +rule proto_compiled_sources
# gazelle:proto_language go -rule proto_compile
# gazelle:proto_language go -rule proto_go_library
```

With this change, we get:

```py
proto_compiled_sources(
    name = "person_go_compiled_sources",
    srcs = ["person.pb.go"],
    output_mappings = ["person.pb.go=github.com/stackb/rules_proto/example/person/person.pb.go"],
    plugins = ["@build_stack_rules_proto//plugin/golang/protobuf:protoc-gen-go"],
    proto = "person_proto",
    visibility = ["//visibility:public"],
)
```

To generate `person.pb.go` and copy it into place, invoke:

```sh
$ bazel run //example/person:person_go_compiled_sources.update
Target //example/person:person_go_compiled_sources.update up-to-date:
  bazel-bin/example/person/person_go_compiled_sources.update.update.json
INFO: Running command line: bazel-bin/example/person/person_go_compiled_sources.update.bash
Target @@//example/person:person_go_compiled_sources: output files copied to source tree:
  example/person/person.pb.go
```

Running gazelle again generates the `go_library` rule using the typical `go`
gazelle extension:

```py
go_library(
    name = "person",
    srcs = ["person.pb.go"],
    importpath = "github.com/stackb/rules_proto/example/person",
    visibility = ["//visibility:public"],
    deps = [
        "//example/place:place_go_proto",
        "@org_golang_google_protobuf//reflect/protoreflect:go_default_library",
        "@org_golang_google_protobuf//runtime/protoimpl:go_default_library",
    ],
)
```

If you have many proto rules, it would be tedious to have to update them all
with individual `bazel run <target>.update` invocations.  To simplify this, use
the `proto_compile_assets` rule:

```py
load("@build_stack_rules_proto//rules:proto_compile_assets.bzl", "proto_compile_assets")

proto_compile_assets(
    name = "assets",
    deps = [
        "//example/person:person_go_compiled_sources",
        "//example/place:place_go_compiled_sources",
        "//example/thing:thing_go_compiled_sources",
    ],
)
```

> NOTE: the `proto_compile_assets` is not generated by gazelle currently.  It
> will require manual rule edits or your own hand-rolled gazelle extension for
> this (not difficult).

### Advanced Topic: external proto deps with `proto_repository`

A common strategy for using protos from other sources is to simply copy them
into your own repo.  This works but can be tedious and hard to maintain.  The
`proto_repository` rule / module extension is like `go_repository` in that it
works like `http_archive` to fetch an external tarball, and then runs `gazelle`
in the external workspace after downloading and extracting the archive.

For example, consider the use case where you want to integrate with the build
event service or one of the remote execution services.  These have complex proto
dependencies that are tedious to get right.

To use this, add the following to your `MODULE.bazel` file:

```py

proto_repository = use_extension("@build_stack_rules_proto//extensions:proto_repository.bzl", "proto_repository", dev_dependency = True)

proto_repository.archive(
    name = "protobufapis",
    build_directives = [
        "gazelle:exclude testdata",
        "gazelle:exclude google/protobuf/compiler/ruby",
        "gazelle:exclude google/protobuf/bridge",
        "gazelle:exclude google/protobuf/util/internal/testdata",
        "gazelle:proto_language go enable true",
    ],
    build_file_proto_mode = "file",
    build_file_generation = "clean",
    cfgs = ["@//:rules_proto_config.yaml"],
    deleted_files = [
        "google/protobuf/*test*.proto",
        "google/protobuf/*unittest*.proto",
        "google/protobuf/late*.proto",
        "google/protobuf/sample*.proto",
        "google/protobuf/cpp_features.proto",
        "google/protobuf/internal_options.proto",
        "google/protobuf/compiler/cpp/*test*.proto",
        "google/protobuf/util/*test*.proto",
        "google/protobuf/util/*unittest*.proto",
        "google/protobuf/util/json_format*.proto",
    ],
    sha256 = "bb1fd58473c47c747a3f00fc45ced1d562bba4bf645db07cc889fe86dee279ca",
    strip_prefix = "protobuf-4fbd1111a292d04746c732573025e3251de0bb9c/src",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/4fbd1111a292d04746c732573025e3251de0bb9c.tar.gz"],
)

# Commit: 60e1300d4a0b60b85b3df167ddc4062ac7cc4f44
# Date: 2025-09-08 18:27:45 +0000 UTC
# URL: https://github.com/googleapis/googleapis/commit/60e1300d4a0b60b85b3df167ddc4062ac7cc4f44
#
# feat: release initial client libraries for Cloud Location Finder (https://cloud.google.com/location-finder)
#
# Clients can now use this v1 client library via the following methods: ListCloudLocations, GetCloudLocation, SearchCloudLocations
#
# PiperOrigin-RevId: 804516757
# Size: 11515047 (12 MB)
proto_repository.archive(
    name = "googleapis",
    build_file_generation = "clean",
    build_directives = [
        "gazelle:exclude google/ads/googleads/v19/services",
        "gazelle:exclude google/ads/googleads/v20/services",
        "gazelle:exclude google/ads/googleads/v21/services",
        "gazelle:exclude google/ads/googleads/v7/services",
        "gazelle:exclude google/ads/googleads/v8/services",
        "gazelle:exclude google/cloud/recommendationengine/v1beta1",
        "gazelle:exclude google/devtools/containeranalysis/v1beta1",
        "gazelle:exclude google/example",
        "gazelle:exclude google/maps/weather/v1",
        "gazelle:proto_language go enable true",
    ],
    build_file_proto_mode = "file",
    cfgs = ["//:rules_proto_config.yaml"],
    imports = ["@protobufapis//:imports.csv"],
    reresolve_known_proto_imports = True,
    sha256 = "b1f729e116312e1bed9a6c0b812e8d6071755dcf93ff4f665c07bbf517dd61a6",
    strip_prefix = "googleapis-60e1300d4a0b60b85b3df167ddc4062ac7cc4f44",
    urls = ["https://github.com/googleapis/googleapis/archive/60e1300d4a0b60b85b3df167ddc4062ac7cc4f44.tar.gz"],
)

# Commit: 6777112ef7defa6705b1ebd2831d6c7efeb12ba2
# Date: 2024-09-26 07:13:55 +0000 UTC
# URL: https://github.com/bazelbuild/remote-apis/commit/6777112ef7defa6705b1ebd2831d6c7efeb12ba2
#
# Improve phrasing of digest_function added in #311 (#312)
#
# - Unlike the other messages where a digest_function field was added,
#   ExecuteActionMetadata is returned by the server to the client. This
#   means that the words "client" and "server" need to be swapped around
#   in some but not all places.
#
# - To ensure backward compatibility, we permit digest functions that were
#   defined at the time this field was added to remain unset. When the
#   digest_function fields were added to the other messages, the last one
#   to be added was MURMUR3. In the meantime we've added BLAKE3 and
#   SHA256TREE, so for this specific field we must list those as well.
# Size: 136595 (137 kB)
proto_repository.archive(
    name = "remoteapis",
    build_directives = [
        "gazelle:exclude third_party",
        "gazelle:proto_language go enable true",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    cfgs = ["//:rules_proto_config.yaml"],
    imports = [
        "@googleapis//:imports.csv",
        "@protobufapis//:imports.csv",
    ],
    sha256 = "7b6847779f18fe0a586c8629b9347cf5e54edb0c9fb7cd7b56c489c0209409c2",
    strip_prefix = "remote-apis-6777112ef7defa6705b1ebd2831d6c7efeb12ba2",
    urls = ["https://github.com/bazelbuild/remote-apis/archive/6777112ef7defa6705b1ebd2831d6c7efeb12ba2.tar.gz"],
)

# Commit: 526225e836561307065f2389c9e2163064fda084
# Date: 2025-09-04 17:49:52 +0000 UTC
# URL: https://github.com/bazelbuild/bazel/commit/526225e836561307065f2389c9e2163064fda084
#
# Release 8.4.0 (2025-09-04)
#
# Release Notes:
# Size: 23065443 (23 MB)
proto_repository.archive(
    name = "bazelapis",
    build_directives = [
        "gazelle:exclude src/java_tools/import_deps_checker/javatests/com/google/devtools/build/importdeps/testdata",
        "gazelle:exclude src/tools/android/java/com/google/devtools/build/android",
        "gazelle:exclude third_party",
        "gazelle:proto_language go enable true",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/stardoc_output.proto=github.com/bazelbuild/bazel/src/main/protobuf/stardoc_output",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/java/com/google/devtools/build/lib/packages/metrics/package_load_metrics.proto=github.com/bazelbuild/bazel/src/main/java/com/google/devtools/build/lib/packages/metrics",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/strategy_policy.proto=github.com/bazelbuild/bazel/src/main/protobuf/strategy_policy",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/execution_graph.proto=github.com/bazelbuild/bazel/src/main/protobuf/execution_graph",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/remote_scrubbing.proto=github.com/bazelbuild/bazel/src/main/protobuf/remote_scrubbing",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/memory_pressure.proto=github.com/bazelbuild/bazel/src/main/protobuf/memory_pressure",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/cache_salt.proto=github.com/bazelbuild/bazel/src/main/protobuf/cache_salt",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/crash_debugging.proto=github.com/bazelbuild/bazel/src/main/protobuf/crash_debugging",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/file_invalidation_data.proto=github.com/bazelbuild/bazel/src/main/protobuf/file_invalidation_data",
        "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/bazel_output_service_rev2.proto=github.com/bazelbuild/bazel/src/main/protobuf/bazel_output_service_rev2",
    ] + [
      "gazelle:proto_plugin protoc-gen-go option Msrc/main/protobuf/{0}.proto=github.com/bazelbuild/bazel/src/main/protobuf/{0}" % name 
      for name in [
        "stardoc_output", "strategy_policy", "execution_graph", "remote_scrubbing", "memory_pressure", "cache_salt", "crash_debugging", "file_invalidation_data", "bazel_output_service_rev2",
      ],
    ],
    deleted_files = [
        "src/main/protobuf/bazel_output_service.proto",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    cfgs = ["//:rules_proto_config.yaml"],
    imports = [
        "@googleapis//:imports.csv",
        "@protobufapis//:imports.csv",
        "@remoteapis//:imports.csv",
    ],
    sha256 = "6037e5df4d97c1402298e19f0c511ae5f44757c1df2e92faad0b2a25c24ae2f8",
    strip_prefix = "bazel-526225e836561307065f2389c9e2163064fda084",
    urls = ["https://github.com/bazelbuild/bazel/archive/526225e836561307065f2389c9e2163064fda084.tar.gz"],
)

use_repo(
    proto_repository,
    "bazelapis",
    "googleapis",
    "protobufapis",
    "remoteapis",
)
```

There is a lot to unpack here:

- `sha256`, `strip_prefix`, `urls` are what you'd expect from an `http_archive`
- `build_directives`, `build_file_generation = "clean"`, and
  `build_file_proto_mode = "file"` are similar to `go_repository`.  Clean will
  delete all pre-existing BUILD files in the archive, and `file` mode is best
  for ensuring dependencies resolve correctly (and fine-grained).
- `cfgs` and `imports` are similar to that from the `proto_gazelle` rule (the
  base YAML config).  The `imports` setup the correct resolve directives such
  that everything stays well-defined.

With this configuration in place, you can include external proto imports in
`.proto` files, or go `imports` into your `.go` files and gazelle will
automatically find them.

```
grep 'build_event_stream' < /private/var/tmp/_bazel_pcj/01715f374a279cb6785f19e874fc1201/external/build_stack_rules_proto++proto_repository+bazelapis/imports.csv
go,go,github.com/bazelbuild/bazelapis/src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream/build_event_stream,@bazelapis//src/main/java/com/google/devtools/build/lib/buildeventstream/proto:build_event_stream_go_proto
proto,proto,src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream.proto,@bazelapis//src/main/java/com/google/devtools/build/lib/buildeventstream/proto:build_event_stream_proto
protobuf,proto_go_library,src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream.proto,@bazelapis//src/main/java/com/google/devtools/build/lib/buildeventstream/proto:build_event_stream_go_proto
```

For example, consider the minimal example in `cmd/bescli/main.go`:

```go
package main

import (
	"log"

	"github.com/bazelbuild/bazelapis/src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream/build_event_stream"
)

func main() {
	log.Printf("hello, %+v", build_event_stream.BuildEventId{})
}
```

Running gazelle will produce:

```py
load("@rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "bescli_lib",
    srcs = ["main.go"],
    importpath = "github.com/stackb/bazel-aquery-differ/cmd/bescli",
    visibility = ["//visibility:private"],
    deps = ["@bazelapis//src/main/java/com/google/devtools/build/lib/buildeventstream/proto:build_event_stream_go_proto"],
)

go_binary(
    name = "bescli",
    embed = [":bescli_lib"],
    visibility = ["//visibility:public"],
)
```

Which is runnable as follows:

```sh
$ bazel run //cmd/bescli
Target //cmd/bescli:bescli up-to-date:
  bazel-bin/cmd/bescli/bescli_/bescli
INFO: Running command line: bazel-bin/cmd/bescli/bescli_/bescli
2025/09/10 20:08:16 hello, {state:{NoUnkeyedLiterals:{} DoNotCompare:[] DoNotCopy:[] atomicMessageInfo:<nil>} Id:<nil> unknownFields:[] sizeCache:0}
```

### advanced topic: vendoring external go proto assets with `proto_go_modules`

> NOTE: this section is specific to golang only

Given the above, we still have the issue that while `bazel` and `gazelle`
understood how to work with `build_event_service`, the IDE and traditional go
tooling like `go mod tidy` do not.  To address this problem, consider
`proto_go_modules`.  This rule copies over generated go proto assets to a
`local/` folder (for reasons beyond the scope of this document, `vendor/` does
not work, it has to be a different directory).

`proto_go_modules` is a separate gazelle extension and a runnable rule.  To use
it, first add the extension to your `gazelle_binary`:

```diff
gazelle_binary(
    name = "gazelle-protobuf",
    languages = [
        "@gazelle//language/proto:go_default_library",
        "@gazelle//language/go:go_default_library",
        "@build_stack_rules_proto//language/protobuf",
+       "@build_stack_rules_proto//language/proto_go_modules",
    ],
)
```

Enable the extension in each `proto_repository` workspaces by naming at least
one rule kind to collect (typically `proto_go_library`):

```diff
proto_repository.archive(
    name = "protobufapis",
    build_directives = [
        ...
        "gazelle:proto_language go enable true",
+       "gazelle:proto_go_modules_index_kind proto_go_library",
    ],
    build_file_proto_mode = "file",
    build_file_generation = "clean",
    cfgs = ["@//:rules_proto_config.yaml"],
    ...
)
```

A new rule will be generated:

```py
# bazel query @protobufapis//:proto_go_modules --output build
proto_go_modules(
    name = "proto_go_modules",
    visibility = ["//visibility:public"],
    deps = [
        "@protobufapis//google/protobuf:any_go_proto",
        "@protobufapis//google/protobuf:api_go_proto",
        "@protobufapis//google/protobuf:descriptor_go_proto",
        "@protobufapis//google/protobuf:duration_go_proto",
        "@protobufapis//google/protobuf:empty_go_proto",
        "@protobufapis//google/protobuf:field_mask_go_proto",
        "@protobufapis//google/protobuf:source_context_go_proto",
        "@protobufapis//google/protobuf:struct_go_proto",
        "@protobufapis//google/protobuf:timestamp_go_proto",
        "@protobufapis//google/protobuf:type_go_proto",
        "@protobufapis//google/protobuf:wrappers_go_proto",
        "@protobufapis//google/protobuf/compiler:plugin_go_proto",
    ],
)
```

This rule `deps` includes all the known `proto_go_library` macros in the
external repository.  These are `go_library` rules (having provider `GoArchive`)
 are used to determine direct and transitive deps.

Next, create a `proto_go_modules` rule in the default workspace manually (no
need for gazelle to do this).  It looks similar but uses the `modules` attribute
to configure a "universe" of available archives transitively.

```py
proto_go_modules(
    name = "proto_go_modules",
    imports = [
        "github.com/bazelbuild/bazelapis/src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream/build_event_stream",
    ],
    srcroot = "genproto",
    modules = [
        "@protobufapis//:proto_go_modules",
        "@bazelapis//:proto_go_modules",
        "@remoteapis//:proto_go_modules",
        "@googleapis//:proto_go_modules",
    ],
)
```

Use the `imports` attribute to name a "top-level" importpath that should be vendored into the default workspace.  All transitive proto modules will also be vendored.

```sh
$ bazel run @//:proto_go_modules
```

```sh
$ tree genproto
genproto
├── github.com
│   └── bazelbuild
│       ├── bazel
│       │   └── src
│       │       └── main
│       │           ├── java
│       │           │   └── com
│       │           │       └── google
│       │           │           └── devtools
│       │           │               └── build
│       │           │                   └── lib
│       │           │                       └── packages
│       │           │                           └── metrics
│       │           │                               ├── go.mod
│       │           │                               └── package_load_metrics.pb.go
│       │           └── protobuf
│       │               └── strategy_policy
│       │                   ├── go.mod
│       │                   └── strategy_policy.pb.go
│       └── bazelapis
│           └── src
│               └── main
│                   ├── java
│                   │   └── com
│                   │       └── google
│                   │           └── devtools
│                   │               └── build
│                   │                   └── lib
│                   │                       └── buildeventstream
│                   │                           └── proto
│                   │                               └── build_event_stream
│                   │                                   └── build_event_stream
│                   │                                       ├── build_event_stream.pb.go
│                   │                                       └── go.mod
│                   └── protobuf
│                       ├── action_cache
│                       │   ├── action_cache.pb.go
│                       │   └── go.mod
│                       ├── command_line
│                       │   ├── command_line.pb.go
│                       │   └── go.mod
│                       ├── failure_details
│                       │   ├── failure_details.pb.go
│                       │   └── go.mod
│                       ├── invocation_policy
│                       │   ├── go.mod
│                       │   └── invocation_policy.pb.go
│                       └── option_filters
│                           ├── go.mod
│                           └── option_filters.pb.go
└── google.golang.org
    └── protobuf
        └── types
            ├── descriptorpb
            │   ├── descriptor.pb.go
            │   └── go.mod
            └── known
                ├── anypb
                │   ├── any.pb.go
                │   └── go.mod
                ├── durationpb
                │   ├── duration.pb.go
                │   └── go.mod
                └── timestamppb
                    ├── go.mod
                    └── timestamp.pb.go
```

```sh
$ cat genproto/github.com/bazelbuild/bazelapis/src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream/build_event_stream/go.mod

module github.com/bazelbuild/bazelapis/src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream/build_event_stream
go 1.23.1
```

The default `go.mod` file has been modified with additional `replace` directives that point to the `genproto/` modules.

```sh
replace github.com/bazelbuild/bazelapis/src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream/build_event_stream => ./genproto/github.com/bazelbuild/bazelapis/src/main/java/com/google/devtools/build/lib/buildeventstream/proto/build_event_stream/build_event_stream
```

## Appendix 1: Registered Flags

- `--proto_configs`: optional config.yaml file(s) that provide preconfiguration
- `--proto_imports_in`: index files to parse and load symbols from
- `--proto_imports_out`: filename where index should be written
- `--proto_repo_name`: external name of this repository
- `--reresolve_known_proto_imports`: if true, re-resolve hardcoded proto_library
  deps on go_googleapis and com_google_protobuf from language/proto from the
  index
- `--proto_rule`: register custom starlark rule of the form
  `<file_name>%<rule_name>`
- `--proto_plugin`: register custom starlark plugin of the form
  `<file_name>%<plugin_name>`

## Appendix 2: Known Directives

### `gazelle:proto_plugin` - Configure Plugin Metadata

The `proto_plugin` directive configures metadata about a protoc plugin. Format:

```
# gazelle:proto_plugin <plugin_name> <parameter> <value>
```

Parameters:
- `implementation`: Specifies the plugin implementation (e.g., `builtin:java`,
  `golang:protobuf:protoc-gen-go`).
  - Plugin implementations are written in go (statically compiled) or starlark
    (dynamically evaluated). Implementations are responsible for building a
    `PluginConfiguration` object that is used to render a `proto_compile` rule.
    An implementation is tied to the implementation of the protoc plugin itself
    and must be able to accurately predict the files that will be produced as a
    result of the compilation inputs and plugin options.  For example, the
    implementation must know that given a proto file `example.proto` having
    `go_package = 'github.com/example/api/v1`, the output produced by `protoc`
    will be `{output_dir}/github.com/example/api/v1/example.pb.go`. 
- `label`: The bazel label that points to a `proto_plugin` rule.  A set of
  `proto_plugin` rules are available out of the box from
  `@stack_build_rules_proto`, or you can bring your own.
- `option`: Plugin-specific options.
- `dep`: Dependencies required by generated code
- `enabled`: Whether the plugin is enabled (`true`/`false`)

###

Examples:
```bash
# Configure the built-in Go plugin
# gazelle:proto_plugin go implementation builtin:go

# Configure a custom plugin with its label
# gazelle:proto_plugin my_plugin label //tools:protoc-gen-my-plugin

# Add options to a plugin (use + prefix to add, - to remove)
# gazelle:proto_plugin grpc_java option grpc
# gazelle:proto_plugin grpc_java +option java_multiple_files=true
# gazelle:proto_plugin grpc_java -option grpc

# Add dependencies to generated code
# gazelle:proto_plugin go_grpc dep @org_golang_google_grpc//:grpc
```

### `gazelle:proto_rule` - Configure Rule Generation

The `proto_rule` directive configures how Bazel rules are generated. Format:

```
# gazelle:proto_rule <rule_name> <parameter> <value>
```

Parameters:
- `implementation`: Specifies the rule implementation (e.g.,
  `stackb:rules_proto:proto_compile`)
- `enabled`: Whether the rule is enabled (`true`/`false`)
- `dep`: Dependencies to add to generated rules
- `attr`: Set arbitrary attributes on generated rules
- `resolve`: Custom import resolution patterns (regex supported)

Examples:
```bash
# Configure a proto_compile rule
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# Enable/disable a rule
# gazelle:proto_rule go_proto_library enabled true

# Add dependencies to generated rules
# gazelle:proto_rule go_proto_library dep @com_github_golang_protobuf//proto:go_default_library
# gazelle:proto_rule go_proto_library +dep //my/lib:go_default_library

# Set custom attributes
# gazelle:proto_rule my_proto_library attr visibility //visibility:public

# Configure custom import resolution
# gazelle:proto_rule go_proto_library resolve google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/$1pb
```

### `gazelle:proto_language` - Compose Languages from Plugins and Rules

The `proto_language` directive binds plugins and rules together to form a
complete language configuration. Format:

```
# gazelle:proto_language <language_name> <parameter> <value>
```

Parameters:
- `enabled`: Whether the language is enabled (`true`/`false`)
- `plugin`: Which plugin(s) to use for this language
- `rule`: Which rule to use for generating targets

Examples:
```bash
# Enable Go proto generation
# gazelle:proto_language go enabled true
# gazelle:proto_language go plugin go
# gazelle:proto_language go rule go_proto_library

# Configure gRPC for Go
# gazelle:proto_language go_grpc enabled true
# gazelle:proto_language go_grpc plugin grpc-go
# gazelle:proto_language go_grpc rule go_grpc_library

# Add multiple plugins to a language
# gazelle:proto_language java plugin java
# gazelle:proto_language java +plugin grpc_java
```

### Intent System (+ and - Prefixes)

Many directive parameters support an "intent" system:
- No prefix or `+` prefix: Add/enable the value
- `-` prefix: Remove/disable the value

This allows incremental configuration across multiple BUILD files:

```bash
# In //proto/BUILD.bazel - base configuration
# gazelle:proto_language go plugin go
# gazelle:proto_language go dep @com_github_golang_protobuf//proto

# In //proto/services/BUILD.bazel - add gRPC support
# gazelle:proto_language go +plugin grpc-go
# gazelle:proto_language go +dep @org_golang_google_grpc//:grpc

# In //proto/legacy/BUILD.bazel - remove gRPC
# gazelle:proto_language go -plugin grpc-go
```
