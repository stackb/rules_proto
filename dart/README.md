# `dart`

| Rule | Description |
| ---: | :--- |
| [dart_proto_compile](#dart_proto_compile) | Generates dart protobuf artifacts |
| [dart_grpc_compile](#dart_grpc_compile) | Generates dart protobuf+gRPC artifacts |
| [dart_proto_library](#dart_proto_library) | Generates dart protobuf library |
| [dart_grpc_library](#dart_grpc_library) | Generates dart protobuf+gRPC library |

---

## `dart_proto_compile`

Generates dart protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//dart:deps.bzl", "dart_proto_compile")

dart_proto_compile()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

dart_repositories()

load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

dart_protoc_plugin_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//dart:dart_proto_compile.bzl", "dart_proto_compile")

dart_proto_compile(
    name = "person_dart_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `Flags`

| Category | Flag | Value | Description |
| --- | --- | --- | --- |
| build | incompatible_no_transitive_loads | false |  |
| build | incompatible_disable_deprecated_attr_params | false |  |
| build | incompatible_enable_cc_toolchain_resolution | false |  |
| build | incompatible_require_ctx_in_configure_features | false |  |
| build | incompatible_depset_is_not_iterable | false |  |
| build | incompatible_depset_union | false |  |
| build | incompatible_disallow_struct_provider_syntax | false |  |

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

---

## `dart_grpc_compile`

Generates dart protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//dart:deps.bzl", "dart_grpc_compile")

dart_grpc_compile()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

dart_repositories()

load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

dart_protoc_plugin_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//dart:dart_grpc_compile.bzl", "dart_grpc_compile")

dart_grpc_compile(
    name = "greeter_dart_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `Flags`

| Category | Flag | Value | Description |
| --- | --- | --- | --- |
| build | incompatible_no_transitive_loads | false |  |
| build | incompatible_disable_deprecated_attr_params | false |  |
| build | incompatible_enable_cc_toolchain_resolution | false |  |
| build | incompatible_require_ctx_in_configure_features | false |  |
| build | incompatible_depset_is_not_iterable | false |  |
| build | incompatible_depset_union | false |  |
| build | incompatible_disallow_struct_provider_syntax | false |  |

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

---

## `dart_proto_library`

Generates dart protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//dart:deps.bzl", "dart_proto_library")

dart_proto_library()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

dart_repositories()

load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

dart_protoc_plugin_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//dart:dart_proto_library.bzl", "dart_proto_library")

dart_proto_library(
    name = "person_dart_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `Flags`

| Category | Flag | Value | Description |
| --- | --- | --- | --- |
| build | incompatible_no_transitive_loads | false |  |
| build | incompatible_disable_deprecated_attr_params | false |  |
| build | incompatible_enable_cc_toolchain_resolution | false |  |
| build | incompatible_require_ctx_in_configure_features | false |  |
| build | incompatible_depset_is_not_iterable | false |  |
| build | incompatible_depset_union | false |  |
| build | incompatible_disallow_struct_provider_syntax | false |  |

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

---

## `dart_grpc_library`

Generates dart protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//dart:deps.bzl", "dart_grpc_library")

dart_grpc_library()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

dart_repositories()

load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

dart_protoc_plugin_deps()

load("@dart_pub_deps_grpc//:deps.bzl", dart_grpc_deps = "pub_deps")

dart_grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//dart:dart_grpc_library.bzl", "dart_grpc_library")

dart_grpc_library(
    name = "greeter_dart_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `Flags`

| Category | Flag | Value | Description |
| --- | --- | --- | --- |
| build | incompatible_no_transitive_loads | false |  |
| build | incompatible_disable_deprecated_attr_params | false |  |
| build | incompatible_enable_cc_toolchain_resolution | false |  |
| build | incompatible_require_ctx_in_configure_features | false |  |
| build | incompatible_depset_is_not_iterable | false |  |
| build | incompatible_depset_union | false |  |
| build | incompatible_disallow_struct_provider_syntax | false |  |

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

