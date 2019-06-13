# `grpc-gateway`

| Rule | Description |
| ---: | :--- |
| [gateway_grpc_compile](#gateway_grpc_compile) | Generates grpc-gateway *.go files |
| [gateway_swagger_compile](#gateway_swagger_compile) | Generates grpc-gateway swagger *.json files |
| [gateway_grpc_library](#gateway_grpc_library) | Generates grpc-gateway library files |

---

## `gateway_grpc_compile`

Generates grpc-gateway *.go files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_grpc_compile")

gateway_grpc_compile()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:gateway_grpc_compile.bzl", "gateway_grpc_compile")

gateway_grpc_compile(
    name = "api_gateway_grpc",
    deps = ["@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway/example/api:api_proto"],
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
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |

---

## `gateway_swagger_compile`

Generates grpc-gateway swagger *.json files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_swagger_compile")

gateway_swagger_compile()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:gateway_swagger_compile.bzl", "gateway_swagger_compile")

gateway_swagger_compile(
    name = "api_gateway_grpc",
    deps = ["@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway/example/api:api_proto"],
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
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |

---

## `gateway_grpc_library`

Generates grpc-gateway library files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_grpc_library")

gateway_grpc_library()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:gateway_grpc_library.bzl", "gateway_grpc_library")

gateway_grpc_library(
    name = "api_gateway_library",
    importpath = "github.com/stackb/rules_proto/github.com/grpc-ecosystem/grpc-gateway/examples/api",
    deps = ["@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway/example/api:api_proto"],
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
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |
