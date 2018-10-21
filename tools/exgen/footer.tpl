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
