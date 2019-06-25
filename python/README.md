# `python`

| Rule | Description |
| ---: | :--- |
| [python_proto_compile](#python_proto_compile) | Generates *.py protobuf artifacts |
| [python_grpc_compile](#python_grpc_compile) | Generates *.py protobuf+gRPC artifacts |
| [python_proto_library](#python_proto_library) | Generates *.py protobuf library |
| [python_grpc_library](#python_grpc_library) | Generates *.py protobuf+gRPC library |

---

## `python_proto_compile`

Generates *.py protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//python:deps.bzl", "python_proto_compile")

python_proto_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:python_proto_compile.bzl", "python_proto_compile")

python_proto_compile(
    name = "person_python_proto",
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
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `python_grpc_compile`

Generates *.py protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//python:deps.bzl", "python_grpc_compile")

python_grpc_compile()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:python_grpc_compile.bzl", "python_grpc_compile")

python_grpc_compile(
    name = "greeter_python_grpc",
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
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `python_proto_library`

Generates *.py protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//python:deps.bzl", "python_proto_library")

python_proto_library()

load("@io_bazel_rules_python//python:pip.bzl", "pip_import", "pip_repositories")

pip_repositories()

pip_import(
    name = "protobuf_py_deps",
    requirements = "@build_stack_rules_proto//python/requirements:protobuf.txt",
)

load("@protobuf_py_deps//:requirements.bzl", protobuf_pip_install = "pip_install")

protobuf_pip_install()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:python_proto_library.bzl", "python_proto_library")

python_proto_library(
    name = "person_python_library",
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
| transitive   | `bool` | `True`    | Generate outputs for both *.proto directly named in `deps` AND all their transitive proto_library dependencies          |
| transitivity   | `string_dict` | `{}`    | Transitive filters to apply when the 'transitive' property is enabled. This string_dict can be used to exclude or explicitly include protos from the compilation list by using `exclude` or `include` respectively as the dict value          |

---

## `python_grpc_library`

Generates *.py protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//python:deps.bzl", "python_grpc_library")

python_grpc_library()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@io_bazel_rules_python//python:pip.bzl", "pip_import", "pip_repositories")

pip_repositories()

pip_import(
    name = "protobuf_py_deps",
    requirements = "@build_stack_rules_proto//python/requirements:protobuf.txt",
)

load("@protobuf_py_deps//:requirements.bzl", protobuf_pip_install = "pip_install")

protobuf_pip_install()

pip_import(
    name = "grpc_py_deps",
    requirements = "@build_stack_rules_proto//python/requirements:grpc.txt",
)

load("@grpc_py_deps//:requirements.bzl", grpc_pip_install = "pip_install")

grpc_pip_install()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:python_grpc_library.bzl", "python_grpc_library")

python_grpc_library(
    name = "greeter_python_library",
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
| transitive   | `bool` | `True`    | Generate outputs for both *.proto directly named in `deps` AND all their transitive proto_library dependencies          |
| transitivity   | `string_dict` | `{}`    | Transitive filters to apply when the 'transitive' property is enabled. This string_dict can be used to exclude or explicitly include protos from the compilation list by using `exclude` or `include` respectively as the dict value          |
