## Overview

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

* If your proto file has services, use the `{lang}_grpc_{compile|library}`
  rule instead.

* Load any external dependencies needed for the rule via the
  `load("@build_stack_rules_proto//{lang}:deps.bzl",
  "{lang}_{proto|grpc}_{compile|library}")`.

## Developers

### Code Layout

Each language `{lang}` has a top-level subdirectory that contains:

1. `{lang}/README.md`: generated documentation for the rule(s).

1. `{lang}/deps.bzl`: contains macro functions that declare repository rule
   dependencies for that language.  The name of the macro corresponds to the
   name of the build rule you'd like to use.

2. `{lang}/{rule}.bzl`: rule implementation(s) of the form
`{lang}_{kind}_{type}`, where `kind` is one of `proto|grpc` and `type` is one of
`compile|library`.

3. `{lang}/BUILD.bazel`: contains `proto_plugin()` declarations for the available
   plugins for that language.

4. `{lang}/example/{rule}/`: contains a generated `WORKSPACE` and `BUILD.bazel`
demonstrating standalone usage.

5. `{lang}/example/routeguide/`: routeguide example implementation, if possible.

The root directory contains the base rule defintions:

* `plugin.bzl`: A build rule that defines the name, tool binary, and options for
  a particular proto plugin.

* `compile.bzl`: A build rule that contains the `proto_compile` rule.  This rule
  calls `protoc` with a given list of plugins and generates output files.

Additional protoc plugins and their rules are scoped to the github repository
name where the plugin resides.  For example, there are 3 grpc-web
implementations in this repo:
`[github.com/improbable-eng/grpc-web](./github.com/improbable-eng/grpc-web),
[github.com/grpc/grpc-web](./github.com/grpc/grpc-web), and
[github.com/stackb/grpc.js](./github.com/stackb/grpc.js).

### Rule Generation

To help maintain consistency of the rule implementations and documentation, all
of the rule implementations are generated by the tool `//tools/rulegen`. Changes
in the main `README.md` should be placed in `tools/rulegen/README.header.md` or
`tools/rulegen/README.footer.md`.  Changes to generated rules should be put in
the source files (example: `tools/rulegen/java.go`).

### Transitivity

Briefly, here's how the rules work:

1. Using the `proto_library` graph, collect all the `*.proto` files directly and
transitively required for the protobuf compilation.

2. Copy the `*.proto` files into a "staging" area in the bazel sandbox such that a
single `-I.` will satisfy all imports.

3. Call `protoc OPTIONS FILELIST` and generate outputs.

The concept of *transitivity* (as defined here) affects which files in the set
`*.proto` files are named in the `FILELIST`.  If we only list direct
dependencies then we say `transitive = False`.  If all files are named, then
`transitive = True`.  The set of files that can be included or excluded from the
`FILELIST` are called *transitivity rules*, which can be defined on a per-rule
or per-plugin basis.  Please grep the codebase for examples of their usage.

## Developing Custom Plugins

Follow the pattern seen in the multiple examples in this repository.  The basic idea is:

1. Load the plugin rule: `load("@build_stack_rules_proto//:plugin.bzl", "proto_plugin")`.
2. Define the rule, giving it a `name`, `options` (not mandatory), `tool`, and
   `outputs`.
3. `tool` is a label that refers to the binary executable for the plugin itself.
4. `outputs` is a list of strings that predicts the pattern of files generated
   by the plugin.  Specifying outputs is the only attribute that requires much
   mental effort.

## Contributing

Contributions welcome; please create Issues or GitHub pull requests.

## ProtoInfo provider

A few notes about the 'proto' provider.  Bazel initially implemented
[providers](https://docs.bazel.build/versions/master/skylark/lib/Provider.html)
using simple strings.  For example, one could write a rule that depends on
`proto_library` rules as follows:

```python
my_custom_rule = rule(
    implementation = my_custom_rule_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "proto_library dependencies",
            mandatory = True,
            providers = ["proto"],
        ),
    }
) 
```

We can then collect a list of the provided "info" objects as follows:

```python
def my_custom_rule_impl(ctx):
    # list<ProtoInfo> A list of ProtoInfo values
    deps = [ dep.proto for dep in ctx.attr.deps ]
```

This style of provider usage is now considered "legacy".  As of bazel 0.22.0
(Jan 2019)
[4817df](https://github.com/bazelbuild/bazel/commit/4817df1e862fe2de3b17905dbbc0c7badff6bb4b),
the
[ProtoInfo](https://docs.bazel.build/versions/master/skylark/lib/ProtoInfo.html)
provider was introduced.  Rules using this provider should be changed to:

```python
my_custom_rule = rule(
    implementation = my_custom_rule_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "proto_library dependencies",
            mandatory = True,
            providers = [ProtoInfo],
        ),
 
```

```python
def my_custom_rule_impl(ctx):
    # <list<ProtoInfo>> A list of ProtoInfo
    deps = [ dep[ProtoInfo] for dep in ctx.attr.deps ]
```

A flag to toggle availability of the "legacy" provider was introduced in
[e3f4ce](https://github.com/bazelbuild/bazel/commit/e3f4ce739f0dbe0fea2b5580655fc96012d162cd)
as `--incompatible_disable_legacy_proto_provider`.  Therefore, if you still want
the legacy behavior specify
`--incompatible_disable_legacy_proto_provider=false`.

Additional reading:

- https://github.com/bazelbuild/bazel/issues/3701
- https://github.com/bazelbuild/bazel/issues/6901
