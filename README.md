# `rules_proto` [![Build Status](https://travis-ci.org/pubref/rules_proto.svg?branch=master)](https://travis-ci.org/pubref/rules_proto)

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

> NOTE: These rules are in a *pre-release* state.  The routeguide examples have
> been developed thus far with the goal of getting them to compile/build, not to
> ensure the routeguide client/server is actually correct.  Do expect everything
> to compile, but not to work right!

## Usage

### Add rules_proto your `WORKSPACE`

Specify the language(s) you'd like use by loading the language-specific
`{lang}/deps.bzl` file and call the macro of your choice corresponding to the
proto rule you plan to use.  To override any dependency, declare it in your
workspace before calling the macro.  Example:

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

RULES_PROTO_VERSION = "{replace-with-commit-id}"
RULES_PROTO_SHA256 = "{replace-with-sha256}"

http_archive(
    name = "build_stack_rules_proto",
    urls = ["https://github.com/pubref/rules_proto/archive/%s.tar.gz" % RULES_PROTO_VERSION],
    sha256 = RULES_PROTO_SHA256,
    strip_prefix = "rules_proto-" + RULES_PROTO_VERSION,
)
```

Refer to the `{lang}/README.md` for instructions on workspace dependencies for
that language. [TODO: write these files!]

### Use the rules in your `BUILD.bazel` files

To build a c++ gRPC library:

```python
load("@build_stack_rules_protobuf//cpp:library.bzl", "cpp_grpc_library")

proto_library(
  name = "proto_library",
  srcs = ["api.proto"],
)

cpp_grpc_library(
  name = "api",
  deps = [":proto_library"],
)
```

## Code Layout

Each language `{lang}` has a top-level subdirectory that contains ~4 files:

1. `deps.bzl`: contains macro functions that declare repository rule
   dependencies for that language.  There are typically `n` macros that
   correspond to the names of the rules in `compile.bzl` and `library.bzl`. Load
   only what you need.
2. `compile.bzl`: contains the rules `{lang}_proto_compile` and
   `{lang}_grpc_compile` (if available).
3. `library.bzl`: contains the rules `{lang}_proto_library` and
   `{lang}_grpc_library` (if available).
4. `BUILD.bazel`: contains `proto_plugin()` declarations for the available
   plugins for that language.

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

## Developing Custom Plugins

Follow the pattern seen in the multiple examples in this repository.  The basic idea is:

1. Load the plugin rule: `load("@build_stack_rules_proto//:plugin.bzl", "proto_plugin")`.
2. Define the rule, giving it a `name`, `options` (not mandatory), `tool`, and
   `outputs`.  
3. `tool` is a label that refers to the binary executable for the plugin itself.
4. `outputs` is a list of strings that predicts the pattern of files generated
   by the plugin.  Specifying outputs is the only attribute that requires much
   mental effort. [TODO: article here with example writing a custom plugin].

## Contributing

Contributions welcome; please create Issues or GitHub pull requests.