# `dart`

| Rule | Description |
| ---: | :--- |
| [dart_proto_compile](#dart_proto_compile) | Generates dart protobuf artifacts |
| [dart_grpc_compile](#dart_grpc_compile) | Generates dart protobuf+gRPC artifacts |
| [dart_proto_library](#dart_proto_library) | Generates dart protobuf library |
| [dart_grpc_library](#dart_grpc_library) | Generates dart protobuf+gRPC library |

---

## `dart_proto_compile`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates dart protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//dart:deps.bzl", "dart_proto_compile")

dart_proto_compile()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
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

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile_attrs", "proto_compile_impl")
load("//:aspect.bzl", "ProtoLibraryAspectNodeInfo", "proto_compile_aspect_attrs", "proto_compile_aspect_impl")
load("//:plugin.bzl", "ProtoPluginInfo")

# "Aspects should be top-level values in extension files that define them."

dart_proto_compile_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    provides = ["proto_compile", ProtoLibraryAspectNodeInfo],
    attr_aspects = ["deps"],
    attrs = proto_compile_aspect_attrs + {
        "_plugins": attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                str(Label("//dart:dart")),
            ],
        ),
    },
)

_rule = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs + {
        "deps": attr.label_list(
            mandatory = True,
            providers = ["proto", "proto_compile", ProtoLibraryAspectNodeInfo],
            aspects = [dart_proto_compile_aspect],
        ),    
    },
)

def dart_proto_compile(**kwargs):
    _rule(
        verbose_string = "%s" % kwargs.get("verbose", 0),
        plugin_options_string = ";".join(kwargs.get("plugin_options", [])),
        **kwargs)

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

## `dart_grpc_compile`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates dart protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//dart:deps.bzl", "dart_grpc_compile")

dart_grpc_compile()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
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

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile_attrs", "proto_compile_impl")
load("//:aspect.bzl", "ProtoLibraryAspectNodeInfo", "proto_compile_aspect_attrs", "proto_compile_aspect_impl")
load("//:plugin.bzl", "ProtoPluginInfo")

# "Aspects should be top-level values in extension files that define them."

dart_grpc_compile_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    provides = ["proto_compile", ProtoLibraryAspectNodeInfo],
    attr_aspects = ["deps"],
    attrs = proto_compile_aspect_attrs + {
        "_plugins": attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                str(Label("//dart:grpc_dart")),
            ],
        ),
    },
)

_rule = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs + {
        "deps": attr.label_list(
            mandatory = True,
            providers = ["proto", "proto_compile", ProtoLibraryAspectNodeInfo],
            aspects = [dart_grpc_compile_aspect],
        ),    
    },
)

def dart_grpc_compile(**kwargs):
    _rule(
        verbose_string = "%s" % kwargs.get("verbose", 0),
        plugin_options_string = ";".join(kwargs.get("plugin_options", [])),
        **kwargs)

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

## `dart_proto_library`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates dart protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//dart:deps.bzl", "dart_proto_library")

dart_proto_library()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
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

### `IMPLEMENTATION`

```python
load("//dart:dart_proto_compile.bzl", "dart_proto_compile")
load("@io_bazel_rules_dart//dart/build_rules:core.bzl", "dart_library")

def dart_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    dart_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )
    dart_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("@vendor_protobuf//:protobuf")),
        ],
        pub_pkg_name = name,
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

## `dart_grpc_library`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates dart protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//dart:deps.bzl", "dart_grpc_library")

dart_grpc_library()

# rules_go used here to compile a wrapper around the protoc-gen-grpc plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
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

### `IMPLEMENTATION`

```python
load("//dart:dart_grpc_compile.bzl", "dart_grpc_compile")
load("@io_bazel_rules_dart//dart/build_rules:core.bzl", "dart_library")

def dart_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    dart_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )
    dart_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("@vendor_protobuf//:protobuf")),
            str(Label("@vendor_grpc//:grpc")),
        ],
        pub_pkg_name = name,
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

