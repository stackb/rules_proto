## Transitivity

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

> TODO: explain this better and provide examples.

## Code Layout

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
