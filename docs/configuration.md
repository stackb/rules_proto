---
layout: default
title: Configuration
permalink: guides/configuration
parent: Guides
nav_order: 2
---

If you haven't done so already, [install gazelle](installation).

Configuration of the `protobuf` gazelle language involves adding *gazelle
directives* to your build files.  Directives in a package are inherited by child
packages, so the typical place to put common directives is in the root
`BUILD.bazel` file.  However, if you have a common subdirectory for your protos
like `./proto/BUILD.bazel`, that would be an appropriate place as well.

There are three directives you'll want to become familiar with:

1. `gazelle:proto_rule`
2. `gazelle:proto_plugin`
3. `gazelle:proto_language`

A simple setup could be:

```python
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile
# gazelle:proto_plugin python implementation builtin:python
# gazelle:proto_language python rule proto_compile
# gazelle:proto_language python plugin python
# gazelle:proto_language python enabled true
```

Let's break this down a bit:

1. The `proto_rule` directive reads as "create a new rule configration named
   "proto_compile" whose golang implementation is registered as
   `stackb:rules_proto:proto_compile`".  The `protobuf` gazelle extension
   maintains a global registry of rule implementations; see
   [rules](/rules_proto/rules) for a list of pre-registered rules, or
   [custom rules](/rules_proto/rules/custom) for a guide to implementing your own.  The config
   name is just a string identifier, choose whatever suits your preference.
1. The `proto_plugin` directive reads as "create a new plugin configration named
   'python' whose golang implementation is registered as `builtin:python`".  The
   `protobuf` gazelle extension maintains a global registry of plugin
   implementations; see [plugins](/rules_proto/plugins) for a list of
   pre-registered plugins, or [custom plugins](/plugins_proto/plugins/custom)
   for a guide to implementing your own.  The config name is just a string
   identifier, choose whatever suits your preference.
1. The `proto_language` directive reads as "create a new language configuration
   and add the 'proto_compile' rule config".  The next line says "add the
   'python' plugin to the 'python' language config.  Finally, the last line
   flips the language to be enabled.  

## proto_rule

The `gazelle:proto_rule` directive is a tuple of strings `NAME KEY VALUE`.
`NAME` is just an identifier.  The list of available keys are (see
[language_rule_config.go] for more info):

| proto_rule key   | description                                                           |
| ---------------- | --------------------------------------------------------------------- |
| `implementation` | Identifier for a rule that has been installed into the rule registry. |
| `visibility`     | Default visibility for the generated rule.                            |
| `dep`            | Bazel label to be added to the list of `deps` for the generated rule. |
| `enabled`        | Toggle to enable/disable rule generation on a per-package basis.      |

`proto_rule` creates/updates a configuration for an entity that knows how to
generate a build rule based on a single `proto_library` rule.

## proto_plugin

The `gazelle:proto_plugin` directive is a tuple of strings `NAME KEY VALUE`.
`NAME` is just an identifier.  The list of available keys are (see
[language_plugin_config.go] for more info):

| proto_plugin key | description                                                                         |
| ---------------- | ----------------------------------------------------------------------------------- |
| `implementation` | Identifier for a plugin that has been installed into the plugin registry.           |
| `option`         | plugin option to be passed to the compiler invocation (e.g. `import_style=closure`) |
| `label`          | Bazel label for the `ProtoPluginInfo` provider. ^1                                    |
| `enabled`        | Toggle to enable/disable the plugin on a per-package basis.                         |

^1: Each `proto_plugin` implementation must also have a corresponding
`proto_plugin` rule (see [proto_plugin.bzl]).  The primary purpose of the
`proto_plugin` rule is to associate the binary `tool` for the plugin.  For
example the tool for `protoc-gen-go` tool is
`@com_github_grpc_grpc_go//cmd/protoc-gen-go`.  Use `label` value to override
the default that is provided by the implementation.

`proto_plugin` creates/updates a configuration for an entity that works with a
rule.  The primary job of a plugin is to predict what output files are going to
be produced by a protoc plugin binary.  For example, the gazelle proto_plugin
implementation for `protoc-gen-go` knows that for any file `foo.proto` with
`package a.b.c`, a file `{EXECROOT}/a/b/c/foo.pb.go` will be generated *unless*
that file also has a go_package option like `go_package github.com/bar/baz:baz`;
in this case the output file will be `{EXECROOT}/github.com/bar/baz/foo.pb.go`.

## proto_language

The `gazelle:proto_language` directive is a tuple of strings `NAME KEY VALUE`.
`NAME` is just an identifier.  The list of available keys are (see
[language_config.go] for more info):

`proto_language` creates/updates a configuration for an entity that associates
rules and plugins together, and can be enabled/disabled on a per-package basis.

[language_rule_config.go]: https://github.com/stackb/rules_proto/language_rule_config.go
[language_plugin_config.go]: https://github.com/stackb/rules_proto/language_plugin_config.go
[language_config.go]: https://github.com/stackb/rules_proto/language_config.go
[proto_plugin.bzl]: https://github.com/stackb/rules_proto/rules/proto_plugin.bzl

## intents

For directives `gazelle:proto_*` whose form is `NAME KEY VALUE`, the `KEY` can
take an "intent modifier" `+KEY` or `-KEY` to express positive or negative
intent (by default, intent is positive so `KEY` and `+KEY` are equivalent).
Examples:

"Remove the mypy plugin from the python language in this BUILD file":

```
# gazelle:proto_language python -plugin mypy
```

"Turn off gRPC for the gogofast plugin in this BUILD file":

```
# gazelle:proto_plugin gogofast -option plugins=grpc
```

