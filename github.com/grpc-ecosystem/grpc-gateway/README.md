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
load("@build_stack_rules_proto//:deps.bzl", "io_bazel_rules_go", "bazel_gazelle")
io_bazel_rules_go()
bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_grpc_compile")
gateway_grpc_compile()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

load("@com_github_grpc_ecosystem_grpc_gateway//:repositories.bzl", grpc_gateway_repositories = "repositories")
grpc_gateway_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:gateway_grpc_compile.bzl", "gateway_grpc_compile")

gateway_grpc_compile(
    name = "api_gateway_grpc",
    deps = ["@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway/example/api:api_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")

def gateway_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//github.com/grpc-ecosystem/grpc-gateway:grpc-gateway")),
        ],
        **kwargs
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

---

## `gateway_swagger_compile`

Generates grpc-gateway swagger *.json files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "io_bazel_rules_go", "bazel_gazelle")
io_bazel_rules_go()
bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_swagger_compile")
gateway_swagger_compile()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

load("@com_github_grpc_ecosystem_grpc_gateway//:repositories.bzl", grpc_gateway_repositories = "repositories")
grpc_gateway_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:gateway_swagger_compile.bzl", "gateway_swagger_compile")

gateway_swagger_compile(
    name = "api_gateway_grpc",
    deps = ["@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway/example/api:api_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")

def gateway_swagger_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//github.com/grpc-ecosystem/grpc-gateway:swagger")),
        ],
        **kwargs
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

---

## `gateway_grpc_library`

Generates grpc-gateway library files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "io_bazel_rules_go", "bazel_gazelle")
io_bazel_rules_go()
bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_grpc_library")
gateway_grpc_library()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

load("@com_github_grpc_ecosystem_grpc_gateway//:repositories.bzl", grpc_gateway_repositories = "repositories")
grpc_gateway_repositories()
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

### `IMPLEMENTATION`

```python
load("//github.com/grpc-ecosystem/grpc-gateway:gateway_grpc_compile.bzl", "gateway_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

def gateway_grpc_library(**kwargs):
    name = kwargs.get("name")
    importpath = kwargs.get("importpath")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    compilers = kwargs.get("compilers")
    if not compilers:
        compilers = [
            "@io_bazel_rules_go//proto:go_grpc",
            "@com_github_grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
        ]

    go_proto_library(
        name = name,
        compilers = compilers,
        importpath = importpath,
		proto = deps[0],
		deps = ["@go_googleapis//google/api:annotations_go_proto"],
        visibility = visibility,
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
