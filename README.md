# `rules_proto` [![Build Status](https://travis-ci.org/pubref/rules_proto.svg?branch=master)](https://travis-ci.org/stackb/rules_proto)

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

Add rules_proto your `WORKSPACE`:

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    urls = ["https://github.com/stackb/rules_proto/archive/v0.9.tar.gz"],
    sha256 = 4329663fe6c523425ad4d3c989a8ac026b04e1acedeceb56aa4b190fa7f3973d,
    strip_prefix = "rules_proto-v0.9",
)
```

Follow instructions in the language-specific `README.md` for additional
workspace dependencies and build rule usage.  


## Rules

| Lang | Rule | Description |
| ---: | :--- | :--- |
| [android](/android) | [android_proto_compile](/android#android_proto_compile) | Generates android protobuf artifacts |
| [android](/android) | [android_grpc_compile](/android#android_grpc_compile) | Generates android protobuf+gRPC artifacts |
| [android](/android) | [android_proto_library](/android#android_proto_library) | Generates android protobuf library |
| [android](/android) | [android_grpc_library](/android#android_grpc_library) | Generates android protobuf+gRPC library |
| [closure](/closure) | [closure_proto_compile](/closure#closure_proto_compile) | Generates closure *.js protobuf+gRPC files |
| [closure](/closure) | [closure_proto_library](/closure#closure_proto_library) | Generates a closure_library with compiled protobuf *.js files |
| [cpp](/cpp) | [cpp_proto_compile](/cpp#cpp_proto_compile) | Generates *.h,*.cc protobuf artifacts |
| [cpp](/cpp) | [cpp_grpc_compile](/cpp#cpp_grpc_compile) | Generates *.h,*.cc protobuf+gRPC artifacts |
| [cpp](/cpp) | [cpp_proto_library](/cpp#cpp_proto_library) | Generates *.h,*.cc protobuf library |
| [cpp](/cpp) | [cpp_grpc_library](/cpp#cpp_grpc_library) | Generates *.h,*.cc protobuf+gRPC library |
| [csharp](/csharp) | [csharp_proto_compile](/csharp#csharp_proto_compile) | Generates csharp protobuf artifacts |
| [csharp](/csharp) | [csharp_grpc_compile](/csharp#csharp_grpc_compile) | Generates csharp protobuf+gRPC artifacts |
| [csharp](/csharp) | [csharp_proto_library](/csharp#csharp_proto_library) | Generates csharp protobuf library |
| [csharp](/csharp) | [csharp_grpc_library](/csharp#csharp_grpc_library) | Generates csharp protobuf+gRPC library |
| [dart](/dart) | [dart_proto_compile](/dart#dart_proto_compile) | Generates dart protobuf artifacts |
| [dart](/dart) | [dart_grpc_compile](/dart#dart_grpc_compile) | Generates dart protobuf+gRPC artifacts |
| [dart](/dart) | [dart_proto_library](/dart#dart_proto_library) | Generates dart protobuf library |
| [dart](/dart) | [dart_grpc_library](/dart#dart_grpc_library) | Generates dart protobuf+gRPC library |
| [go](/go) | [go_proto_compile](/go#go_proto_compile) | Generates *.go protobuf artifacts |
| [go](/go) | [go_grpc_compile](/go#go_grpc_compile) | Generates *.go protobuf+gRPC artifacts |
| [go](/go) | [go_proto_library](/go#go_proto_library) | Generates *.go protobuf library |
| [go](/go) | [go_grpc_library](/go#go_grpc_library) | Generates *.go protobuf+gRPC library |
| [java](/java) | [java_proto_compile](/java#java_proto_compile) | Generates a srcjar with protobuf *.java files |
| [java](/java) | [java_grpc_compile](/java#java_grpc_compile) | Generates a srcjar with protobuf+gRPC *.java files |
| [java](/java) | [java_proto_library](/java#java_proto_library) | Generates a jar with compiled protobuf *.class files |
| [java](/java) | [java_grpc_library](/java#java_grpc_library) | Generates a jar with compiled protobuf+gRPC *.class files |
| [node](/node) | [node_proto_compile](/node#node_proto_compile) | Generates node *.js protobuf artifacts |
| [node](/node) | [node_grpc_compile](/node#node_grpc_compile) | Generates node *.js protobuf+gRPC artifacts |
| [node](/node) | [node_proto_library](/node#node_proto_library) | Generates node *.js protobuf library |
| [node](/node) | [node_grpc_library](/node#node_grpc_library) | Generates node *.js protobuf+gRPC library |
| [objc](/objc) | [objc_proto_compile](/objc#objc_proto_compile) | Generates objc protobuf artifacts |
| [objc](/objc) | [objc_grpc_compile](/objc#objc_grpc_compile) | Generates objc protobuf+gRPC artifacts |
| [php](/php) | [php_proto_compile](/php#php_proto_compile) | Generates php protobuf artifacts |
| [php](/php) | [php_grpc_compile](/php#php_grpc_compile) | Generates php protobuf+gRPC artifacts |
| [python](/python) | [python_proto_compile](/python#python_proto_compile) | Generates *.py protobuf artifacts |
| [python](/python) | [python_grpc_compile](/python#python_grpc_compile) | Generates *.py protobuf+gRPC artifacts |
| [python](/python) | [python_proto_library](/python#python_proto_library) | Generates *.py protobuf library |
| [python](/python) | [python_grpc_library](/python#python_grpc_library) | Generates *.py protobuf+gRPC library |
| [ruby](/ruby) | [ruby_proto_compile](/ruby#ruby_proto_compile) | Generates *.ruby protobuf artifacts |
| [ruby](/ruby) | [ruby_grpc_compile](/ruby#ruby_grpc_compile) | Generates *.ruby protobuf+gRPC artifacts |
| [ruby](/ruby) | [ruby_proto_library](/ruby#ruby_proto_library) | Generates *.rb protobuf library |
| [ruby](/ruby) | [ruby_grpc_library](/ruby#ruby_grpc_library) | Generates *.rb protobuf+gRPC library |
| [rust](/rust) | [rust_proto_compile](/rust#rust_proto_compile) | Generates rust protobuf artifacts |
| [rust](/rust) | [rust_grpc_compile](/rust#rust_grpc_compile) | Generates rust protobuf+gRPC artifacts |
| [rust](/rust) | [rust_proto_library](/rust#rust_proto_library) | Generates rust protobuf library |
| [rust](/rust) | [rust_grpc_library](/rust#rust_grpc_library) | Generates rust protobuf+gRPC library |
| [scala](/scala) | [scala_proto_compile](/scala#scala_proto_compile) | Generates *.scala protobuf artifacts |
| [scala](/scala) | [scala_grpc_compile](/scala#scala_grpc_compile) | Generates *.scala protobuf+gRPC artifacts |
| [scala](/scala) | [scala_proto_library](/scala#scala_proto_library) | Generates *.py protobuf library |
| [scala](/scala) | [scala_grpc_library](/scala#scala_grpc_library) | Generates *.py protobuf+gRPC library |
| [swift](/swift) | [swift_proto_compile](/swift#swift_proto_compile) | Generates swift protobuf artifacts |
| [swift](/swift) | [swift_grpc_compile](/swift#swift_grpc_compile) | Generates swift protobuf+gRPC artifacts |
| [swift](/swift) | [swift_proto_library](/swift#swift_proto_library) | Generates swift protobuf library |
| [swift](/swift) | [swift_grpc_library](/swift#swift_grpc_library) | Generates swift protobuf+gRPC library |
| [gogo](/github.com/gogo/protobuf) | [gogo_proto_compile](/github.com/gogo/protobuf#gogo_proto_compile) | Generates gogo protobuf artifacts |
| [gogo](/github.com/gogo/protobuf) | [gogo_grpc_compile](/github.com/gogo/protobuf#gogo_grpc_compile) | Generates gogo protobuf+gRPC artifacts |
| [gogo](/github.com/gogo/protobuf) | [gogo_proto_library](/github.com/gogo/protobuf#gogo_proto_library) | Generates gogo protobuf library |
| [gogo](/github.com/gogo/protobuf) | [gogo_grpc_library](/github.com/gogo/protobuf#gogo_grpc_library) | Generates gogo protobuf+gRPC library |
| [gogo](/github.com/gogo/protobuf) | [gogofast_proto_compile](/github.com/gogo/protobuf#gogofast_proto_compile) | Generates gogofast protobuf artifacts |
| [gogo](/github.com/gogo/protobuf) | [gogofast_grpc_compile](/github.com/gogo/protobuf#gogofast_grpc_compile) | Generates gogofast protobuf+gRPC artifacts |
| [gogo](/github.com/gogo/protobuf) | [gogofast_proto_library](/github.com/gogo/protobuf#gogofast_proto_library) | Generates gogofast protobuf library |
| [gogo](/github.com/gogo/protobuf) | [gogofast_grpc_library](/github.com/gogo/protobuf#gogofast_grpc_library) | Generates gogofast protobuf+gRPC library |
| [gogo](/github.com/gogo/protobuf) | [gogofaster_proto_compile](/github.com/gogo/protobuf#gogofaster_proto_compile) | Generates gogofaster protobuf artifacts |
| [gogo](/github.com/gogo/protobuf) | [gogofaster_grpc_compile](/github.com/gogo/protobuf#gogofaster_grpc_compile) | Generates gogofaster protobuf+gRPC artifacts |
| [gogo](/github.com/gogo/protobuf) | [gogofaster_proto_library](/github.com/gogo/protobuf#gogofaster_proto_library) | Generates gogofaster protobuf library |
| [gogo](/github.com/gogo/protobuf) | [gogofaster_grpc_library](/github.com/gogo/protobuf#gogofaster_grpc_library) | Generates gogofaster protobuf+gRPC library |
| [grpc_web](/github.com/grpc/grpc-web) | [closure_grpc_compile](/github.com/grpc/grpc-web#closure_grpc_compile) | Generates a closure *.js protobuf+gRPC files |
| [grpc_web](/github.com/grpc/grpc-web) | [commonjs_grpc_compile](/github.com/grpc/grpc-web#commonjs_grpc_compile) | Generates a commonjs *.js protobuf+gRPC files |
| [grpc_web](/github.com/grpc/grpc-web) | [commonjs_dts_grpc_compile](/github.com/grpc/grpc-web#commonjs_dts_grpc_compile) | Generates a commonjs_dts *.js protobuf+gRPC files |
| [grpc_web](/github.com/grpc/grpc-web) | [ts_grpc_compile](/github.com/grpc/grpc-web#ts_grpc_compile) | Generates a commonjs *.ts protobuf+gRPC files |
| [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_grpc_compile](/github.com/grpc-ecosystem/grpc-gateway#gateway_grpc_compile) | Generates grpc-gateway *.go files |
| [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_swagger_compile](/github.com/grpc-ecosystem/grpc-gateway#gateway_swagger_compile) | Generates grpc-gateway swagger *.json files |
| [grpc-gateway](/github.com/grpc-ecosystem/grpc-gateway) | [gateway_grpc_library](/github.com/grpc-ecosystem/grpc-gateway#gateway_grpc_library) | Generates grpc-gateway library files |

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

