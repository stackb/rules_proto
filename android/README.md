# `android`

| Rule | Description |
| ---: | :--- |
| [android_proto_compile](#android_proto_compile) | Generates android protobuf artifacts |
| [android_grpc_compile](#android_grpc_compile) | Generates android protobuf+gRPC artifacts |
| [android_proto_library](#android_proto_library) | Generates android protobuf library |
| [android_grpc_library](#android_grpc_library) | Generates android protobuf+gRPC library |

---

## `android_proto_compile`

Generates android protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//android:deps.bzl", "android_proto_compile")

android_proto_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//android:android_proto_compile.bzl", "android_proto_compile")

android_proto_compile(
    name = "person_android_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")

def android_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//android:javalite")),
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

## `android_grpc_compile`

Generates android protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")

io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
)

load("@build_stack_rules_proto//android:deps.bzl", "android_grpc_compile")

android_grpc_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//android:android_grpc_compile.bzl", "android_grpc_compile")

android_grpc_compile(
    name = "greeter_android_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")

def android_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//android:javalite")),
            str(Label("//android:grpc_javalite")),
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

## `android_proto_library`

Generates android protobuf library

### `WORKSPACE`

```python

# The set of dependencies loaded here is excessive for android proto alone
# (but simplifies our setup) 
load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")
io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
)

load("@build_stack_rules_proto//android:deps.bzl", "android_proto_library")
android_proto_library()

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")
android_sdk_repository(name = "androidsdk")

load("@gmaven_rules//:gmaven.bzl", "gmaven_rules")
gmaven_rules()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//android:android_proto_library.bzl", "android_proto_library")

android_proto_library(
    name = "person_android_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//android:android_proto_compile.bzl", "android_proto_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def android_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    android_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    android_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("//android:proto_deps")),
        ],
        exports = [
            str(Label("//android:proto_deps")),
        ],
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

---

## `android_grpc_library`

Generates android protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")
io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
)

load("@build_stack_rules_proto//android:deps.bzl", "android_grpc_library")
android_grpc_library()

#load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")
#grpc_deps()

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")
android_sdk_repository(name = "androidsdk")

load("@gmaven_rules//:gmaven.bzl", "gmaven_rules")
gmaven_rules()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//android:android_grpc_library.bzl", "android_grpc_library")

android_grpc_library(
    name = "greeter_android_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//android:android_grpc_compile.bzl", "android_grpc_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def android_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    android_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    android_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("//android:grpc_deps")),
        ],
        exports = [
            str(Label("//android:grpc_deps")),
        ],
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

