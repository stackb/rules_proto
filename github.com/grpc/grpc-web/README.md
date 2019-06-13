# `grpc-web`

| Rule | Description |
| ---: | :--- |
| [closure_grpc_compile](#closure_grpc_compile) | Generates a closure *.js protobuf+gRPC files |
| [commonjs_grpc_compile](#commonjs_grpc_compile) | Generates a commonjs *.js protobuf+gRPC files |
| [commonjs_dts_grpc_compile](#commonjs_dts_grpc_compile) | Generates a commonjs_dts *.js protobuf+gRPC files |
| [ts_grpc_compile](#ts_grpc_compile) | Generates a commonjs *.ts protobuf+gRPC files |
| [closure_grpc_library](#closure_grpc_library) | Generates protobuf closure library *.js files |

---

## `closure_grpc_compile`

Generates a closure *.js protobuf+gRPC files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:deps.bzl", "closure_grpc_compile")

closure_grpc_compile()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:closure_grpc_compile.bzl", "closure_grpc_compile")

closure_grpc_compile(
    name = "greeter_grpc-web_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (`native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins)          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `True`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |

---

## `commonjs_grpc_compile`

Generates a commonjs *.js protobuf+gRPC files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:deps.bzl", "commonjs_grpc_compile")

commonjs_grpc_compile()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:commonjs_grpc_compile.bzl", "commonjs_grpc_compile")

commonjs_grpc_compile(
    name = "greeter_grpc-web_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (`native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins)          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `True`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |

---

## `commonjs_dts_grpc_compile`

Generates a commonjs_dts *.js protobuf+gRPC files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:deps.bzl", "commonjs_dts_grpc_compile")

commonjs_dts_grpc_compile()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:commonjs_dts_grpc_compile.bzl", "commonjs_dts_grpc_compile")

commonjs_dts_grpc_compile(
    name = "greeter_grpc-web_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (`native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins)          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `True`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |

---

## `ts_grpc_compile`

Generates a commonjs *.ts protobuf+gRPC files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:deps.bzl", "ts_grpc_compile")

ts_grpc_compile()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:ts_grpc_compile.bzl", "ts_grpc_compile")

ts_grpc_compile(
    name = "greeter_grpc-web_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (`native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins)          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `True`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |

---

## `closure_grpc_library`

Generates protobuf closure library *.js files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:deps.bzl", "closure_grpc_library")

closure_grpc_library()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc/grpc-web:closure_grpc_library.bzl", "closure_grpc_library")

closure_grpc_library(
    name = "greeter_grpc-web_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `Flags`

| Category | Flag | Value | Description |
| --- | --- | --- | --- |
| build | incompatible_disallow_struct_provider_syntax | false |  |
| build | incompatible_use_toolchain_resolution_for_java_rules | false |  |

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (`native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins)          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `True`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |
