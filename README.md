# rules_proto

`rules_proto` provides build rules for compiling protocol buffer files.

## Overview

Each language `${LANG}` has a top-level subdirectory that contains ~4 files:

1. `deps.bzl`: contains macro functions that declare repository rule
   dependencies for that language.  There are typically 2 functions
   `${LANG}_proto_deps()` and `${LANG}_grpc_deps()`.  If you only need proto
   support (any not grpc support), load the `_proto_` macro, otherwise load the
   `_grpc_` macro function and call it in your `WORKSPACE`.
2. `compile.bzl`: contains the rules `${LANG}_proto_compile` and
   `${LANG}_grpc_compile` (if available).
3. `library.bzl`: contains the rules `${LANG}_proto_library` and
   `${LANG}_grpc_library` (if available).
4. `BUILD.bazel`: contains `proto_plugin()` declarations for the available
   plugins. These are consumed by the `${LANG}_proto_compile` and
   `${LANG}_grpc_compile` rules.

The root directory contains the base rule defintions:

* `plugin.bzl`: A build rule that defines the name, tool binary, and options for
  a particular proto plugin.

* `compile.bzl`: A build rule that contains the `proto_compile` rule.  This rule
  calls protoc with a given list of plugins and generates output files.

## Developing Custom Plugins

Follow the pattern seen in the multiple examples in this repository.  The basic idea is:

1. Load the plugin rule: `load("@org_pubref_rules_proto//:plugin.bzl", "proto_plugin")`.
2. Define the rule, giving it a `name`, `options` (not mandatory), `tool`, and
   `outputs`.  
3. `tool` is a label that refers to the binary executable for the plugin itself.
4. `outputs` is a list of strings that predicts the pattern of files generated
   by the plugin.  Specifying outputs is the only attribute that requires much
   mental effort.

