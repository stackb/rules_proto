# `swift`


**NOTE**: The swift rules are essentially non-functional.  The protoc-plugin "core dumps" despite all efforts thus far on linux.

| Rule | Description |
| ---: | :--- |
| [swift_proto_compile](#swift_proto_compile) | Generates swift protobuf artifacts |
| [swift_grpc_compile](#swift_grpc_compile) | Generates swift protobuf+gRPC artifacts |
| [swift_proto_library](#swift_proto_library) | Generates swift protobuf library |
| [swift_grpc_library](#swift_grpc_library) | Generates swift protobuf+gRPC library |

---

## `swift_proto_compile`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates swift protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//swift:deps.bzl", "swift_proto_compile")
swift_proto_compile()

# rules_go used here to compile a wrapper around the protoc-gen-swift plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@build_stack_rules_proto//swift:repositories.bzl", "swift_toolchain")
# Default values work with linux, x86_64, /usr/local/bin/clang. 
swift_toolchain(
	#root = "/home/pcj/.local/share/umake/swift/swift-lang/usr",
)

# Uncomment for ocal development with swift installed on your machine
# load("@build_bazel_rules_swift//swift:repositories.bzl", "swift_rules_dependencies")
# swift_rules_dependencies()


```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//swift:swift_proto_compile.bzl", "swift_proto_compile")

swift_proto_compile(
    name = "person_swift_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile_attrs", "proto_compile_impl")
load("//:aspect.bzl", "ProtoLibraryAspectNodeInfo", "proto_compile_aspect_attrs", "proto_compile_aspect_impl")
load("//:plugin.bzl", "ProtoPluginInfo")

# "Aspects should be top-level values in extension files that define them."

_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    attr_aspects = ["deps"],
    attrs = proto_compile_aspect_attrs + {
        "_plugins": attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                str(Label("//swift:swift")),
            ],
        ),
    },
)

_rule = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs + {
        "deps": attr.label_list(
            mandatory = True,
            providers = ["proto", ProtoLibraryAspectNodeInfo],
            aspects = [_aspect],
        ),    
    },
)

def swift_proto_compile(**kwargs):
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

## `swift_grpc_compile`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates swift protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//swift:deps.bzl", "swift_grpc_compile")
swift_grpc_compile()

# rules_go used here to compile a wrapper around the protoc-gen-swift plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@build_stack_rules_proto//swift:repositories.bzl", "swift_toolchain")
# Default values work with linux, x86_64, /usr/local/bin/clang. 
swift_toolchain(
	#root = "/home/pcj/.local/share/umake/swift/swift-lang/usr",
)

# Uncomment for ocal development with swift installed on your machine
# load("@build_bazel_rules_swift//swift:repositories.bzl", "swift_rules_dependencies")
# swift_rules_dependencies()


```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//swift:swift_grpc_compile.bzl", "swift_grpc_compile")

swift_grpc_compile(
    name = "greeter_swift_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile_attrs", "proto_compile_impl")
load("//:aspect.bzl", "ProtoLibraryAspectNodeInfo", "proto_compile_aspect_attrs", "proto_compile_aspect_impl")
load("//:plugin.bzl", "ProtoPluginInfo")

# "Aspects should be top-level values in extension files that define them."

_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    attr_aspects = ["deps"],
    attrs = proto_compile_aspect_attrs + {
        "_plugins": attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                str(Label("//swift:grpc_swift")),
            ],
        ),
    },
)

_rule = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs + {
        "deps": attr.label_list(
            mandatory = True,
            providers = ["proto", ProtoLibraryAspectNodeInfo],
            aspects = [_aspect],
        ),    
    },
)

def swift_grpc_compile(**kwargs):
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

## `swift_proto_library`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates swift protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//swift:deps.bzl", "swift_proto_library")
swift_proto_library()

# rules_go used here to compile a wrapper around the protoc-gen-swift plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@build_stack_rules_proto//swift:repositories.bzl", "swift_toolchain")
# Default values work with linux, x86_64, /usr/local/bin/clang. 
swift_toolchain(
	#root = "/home/pcj/.local/share/umake/swift/swift-lang/usr",
)

# Uncomment for ocal development with swift installed on your machine
# load("@build_bazel_rules_swift//swift:repositories.bzl", "swift_rules_dependencies")
# swift_rules_dependencies()


```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//swift:swift_proto_library.bzl", "swift_proto_library")

swift_proto_library(
    name = "person_swift_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//swift:swift_proto_compile.bzl", "swift_proto_compile")
load("@build_bazel_rules_swift//swift:swift.bzl", _swift_proto_library = "swift_proto_library")

def swift_proto_library(**kwargs):
    _swift_proto_library(**kwargs)

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

## `swift_grpc_library`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates swift protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//swift:deps.bzl", "swift_grpc_library")
swift_grpc_library()

# rules_go used here to compile a wrapper around the protoc-gen-swift plugin
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@build_stack_rules_proto//swift:repositories.bzl", "swift_toolchain")
# Default values work with linux, x86_64, /usr/local/bin/clang. 
swift_toolchain(
	#root = "/home/pcj/.local/share/umake/swift/swift-lang/usr",
)

# Uncomment for ocal development with swift installed on your machine
# load("@build_bazel_rules_swift//swift:repositories.bzl", "swift_rules_dependencies")
# swift_rules_dependencies()


```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//swift:swift_grpc_library.bzl", "swift_grpc_library")

swift_grpc_library(
    name = "greeter_swift_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//swift:swift_grpc_compile.bzl", "swift_grpc_compile")
load("@build_bazel_rules_swift//swift:swift.bzl", _swift_proto_library = "swift_proto_library")

def swift_grpc_library(**kwargs):
    _swift_proto_library(**kwargs)

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

