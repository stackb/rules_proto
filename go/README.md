# `go`

| Rule | Description |
| ---: | :--- |
| [go_proto_compile](#go_proto_compile) | Generates *.go protobuf artifacts |
| [go_grpc_compile](#go_grpc_compile) | Generates *.go protobuf+gRPC artifacts |
| [go_proto_library](#go_proto_library) | Generates *.go protobuf library |
| [go_grpc_library](#go_grpc_library) | Generates *.go protobuf+gRPC library |

---

## `go_proto_compile`

Generates *.go protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//go:deps.bzl", "go_proto_compile")

go_proto_compile()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//go:go_proto_compile.bzl", "go_proto_compile")

go_proto_compile(
    name = "person_go_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
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
| has_services   | `bool` | `False`    | If the proto files(s) have a service rpc, generate grpc outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `go_grpc_compile`

Generates *.go protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//go:deps.bzl", "go_grpc_compile")

go_grpc_compile()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//go:go_grpc_compile.bzl", "go_grpc_compile")

go_grpc_compile(
    name = "greeter_go_grpc",
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
| has_services   | `bool` | `False`    | If the proto files(s) have a service rpc, generate grpc outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `go_proto_library`

Generates *.go protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//go:deps.bzl", "go_proto_library")

go_proto_library()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//go:go_proto_library.bzl", "go_proto_library")

go_proto_library(
    name = "person_go_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/go/example/go_proto_library/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
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
| has_services   | `bool` | `False`    | If the proto files(s) have a service rpc, generate grpc outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `go_grpc_library`

Generates *.go protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//go:deps.bzl", "go_grpc_library")

go_grpc_library()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//go:go_grpc_library.bzl", "go_grpc_library")

go_grpc_library(
    name = "greeter_go_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/go/example/go_grpc_library/greeter",
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
| has_services   | `bool` | `False`    | If the proto files(s) have a service rpc, generate grpc outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

