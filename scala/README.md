# `scala`

Rules for `scala_grpc_{compile|library}` don't produce code that compiles!  Use `@//io_bazel_rules_scala//scala_proto:scala_proto.bzl` instead

| Rule | Description |
| ---: | :--- |
| [scala_proto_compile](#scala_proto_compile) | Generates *.scala protobuf artifacts |
| [scala_grpc_compile](#scala_grpc_compile) | Generates *.scala protobuf+gRPC artifacts |
| [scala_proto_library](#scala_proto_library) | Generates *.scala protobuf library |
| [scala_grpc_library](#scala_grpc_library) | Generates *.scala protobuf+gRPC library |

---

## `scala_proto_compile`

Generates *.scala protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//scala:deps.bzl", "scala_proto_compile")

scala_proto_compile()

# rules_go used here to compile a wrapper around the protoc-gen-scala plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//scala:scala_proto_compile.bzl", "scala_proto_compile")

scala_proto_compile(
    name = "person_scala_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins). To apply plugin specific options, use the `options` attribute on `proto_plugin`          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |

---

## `scala_grpc_compile`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates *.scala protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//scala:deps.bzl", "scala_grpc_compile")

scala_grpc_compile()

# rules_go used here to compile a wrapper around the protoc-gen-scala plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//scala:scala_grpc_compile.bzl", "scala_grpc_compile")

scala_grpc_compile(
    name = "greeter_scala_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins). To apply plugin specific options, use the `options` attribute on `proto_plugin`          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |

---

## `scala_proto_library`

Generates *.scala protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//scala:deps.bzl", "scala_proto_library")

scala_proto_library()

# rules_go used here to compile a wrapper around the protoc-gen-scala plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//scala:scala_proto_library.bzl", "scala_proto_library")

scala_proto_library(
    name = "person_scala_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins). To apply plugin specific options, use the `options` attribute on `proto_plugin`          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |

---

## `scala_grpc_library`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates *.scala protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//scala:deps.bzl", "scala_grpc_library")

scala_grpc_library()

# rules_go used here to compile a wrapper around the protoc-gen-scala plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//scala:scala_grpc_library.bzl", "scala_grpc_library")

scala_grpc_library(
    name = "greeter_scala_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins). To apply plugin specific options, use the `options` attribute on `proto_plugin`          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
