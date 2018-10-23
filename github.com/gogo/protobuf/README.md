# `gogo`

| Rule | Description |
| ---: | :--- |
| [gogo_proto_compile](#gogo_proto_compile) | Generates gogo protobuf artifacts |
| [gogo_grpc_compile](#gogo_grpc_compile) | Generates gogo protobuf+gRPC artifacts |
| [gogo_proto_library](#gogo_proto_library) | Generates gogo protobuf library |
| [gogo_grpc_library](#gogo_grpc_library) | Generates gogo protobuf+gRPC library |
| [gogofast_proto_compile](#gogofast_proto_compile) | Generates gogofast protobuf artifacts |
| [gogofast_grpc_compile](#gogofast_grpc_compile) | Generates gogofast protobuf+gRPC artifacts |
| [gogofast_proto_library](#gogofast_proto_library) | Generates gogofast protobuf library |
| [gogofast_grpc_library](#gogofast_grpc_library) | Generates gogofast protobuf+gRPC library |
| [gogofaster_proto_compile](#gogofaster_proto_compile) | Generates gogofaster protobuf artifacts |
| [gogofaster_grpc_compile](#gogofaster_grpc_compile) | Generates gogofaster protobuf+gRPC artifacts |
| [gogofaster_proto_library](#gogofaster_proto_library) | Generates gogofaster protobuf library |
| [gogofaster_grpc_library](#gogofaster_grpc_library) | Generates gogofaster protobuf+gRPC library |

---

## `gogo_proto_compile`

Generates gogo protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_proto_compile")

gogo_proto_compile()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogo_proto_compile.bzl", "gogo_proto_compile")

gogo_proto_compile(
    name = "person_gogo_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def gogo_proto_compile(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_gogo_protobuf//protoc-gen-gogo",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//github.com/gogo/protobuf:gogo"))]

    proto_compile(
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogo_grpc_compile`

Generates gogo protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_grpc_compile")

gogo_grpc_compile()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogo_grpc_compile.bzl", "gogo_grpc_compile")

gogo_grpc_compile(
    name = "greeter_gogo_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def gogo_grpc_compile(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_gogo_protobuf//protoc-gen-gogo",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//github.com/gogo/protobuf:gogo"))]

    proto_compile(
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogo_proto_library`

Generates gogo protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_proto_library")

gogo_proto_library()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogo_proto_library.bzl", "gogo_proto_library")

gogo_proto_library(
    name = "person_gogo_library",
    importpath = "github.com/stackb/rules_proto/gogo/example/gogo_proto_library/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
    go_deps = [
		"@com_github_gogo_protobuf//types:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)
```

### `IMPLEMENTATION`

```python
load("//github.com/gogo/protobuf:gogo_proto_compile.bzl", "gogo_proto_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
	"google/protobuf/any.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def gogo_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    gogo_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})) + wkt_mappings,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
        ],
        importpath = importpath,
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogo_grpc_library`

Generates gogo protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_grpc_library")

gogo_grpc_library()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogo_grpc_library.bzl", "gogo_grpc_library")

gogo_grpc_library(
    name = "greeter_gogo_library",
    importpath = "github.com/stackb/rules_proto/gogo/example/gogo_grpc_library/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)
```

### `IMPLEMENTATION`

```python
load("//github.com/gogo/protobuf:gogo_grpc_compile.bzl", "gogo_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
	"google/protobuf/any.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def gogo_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    gogo_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})) + wkt_mappings,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = importpath,
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogofast_proto_compile`

Generates gogofast protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogofast_proto_compile")

gogofast_proto_compile()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogofast_proto_compile.bzl", "gogofast_proto_compile")

gogofast_proto_compile(
    name = "person_gogo_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def gogofast_proto_compile(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_gogo_protobuf//protoc-gen-gogofast",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//github.com/gogo/protobuf:gogofast"))]

    proto_compile(
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogofast_grpc_compile`

Generates gogofast protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogofast_grpc_compile")

gogofast_grpc_compile()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogofast_grpc_compile.bzl", "gogofast_grpc_compile")

gogofast_grpc_compile(
    name = "greeter_gogo_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def gogofast_grpc_compile(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_gogo_protobuf//protoc-gen-gogofast",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//github.com/gogo/protobuf:gogofast"))]

    proto_compile(
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogofast_proto_library`

Generates gogofast protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogofast_proto_library")

gogofast_proto_library()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogofast_proto_library.bzl", "gogofast_proto_library")

gogofast_proto_library(
    name = "person_gogo_library",
    importpath = "github.com/stackb/rules_proto/gogo/example/gogofast_proto_library/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
    go_deps = [
		"@com_github_gogo_protobuf//types:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)
```

### `IMPLEMENTATION`

```python
load("//github.com/gogo/protobuf:gogofast_proto_compile.bzl", "gogofast_proto_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
	"google/protobuf/any.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def gogofast_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    gogofast_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})) + wkt_mappings,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
        ],
        importpath = importpath,
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogofast_grpc_library`

Generates gogofast protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogofast_grpc_library")

gogofast_grpc_library()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogofast_grpc_library.bzl", "gogofast_grpc_library")

gogofast_grpc_library(
    name = "greeter_gogo_library",
    importpath = "github.com/stackb/rules_proto/gogo/example/gogofast_grpc_library/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)
```

### `IMPLEMENTATION`

```python
load("//github.com/gogo/protobuf:gogofast_grpc_compile.bzl", "gogofast_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
	"google/protobuf/any.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def gogofast_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    gogofast_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})) + wkt_mappings,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = importpath,
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogofaster_proto_compile`

Generates gogofaster protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogofaster_proto_compile")

gogofaster_proto_compile()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogofaster_proto_compile.bzl", "gogofaster_proto_compile")

gogofaster_proto_compile(
    name = "person_gogo_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def gogofaster_proto_compile(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_gogo_protobuf//protoc-gen-gogofaster",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//github.com/gogo/protobuf:gogofaster"))]

    proto_compile(
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogofaster_grpc_compile`

Generates gogofaster protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogofaster_grpc_compile")

gogofaster_grpc_compile()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogofaster_grpc_compile.bzl", "gogofaster_grpc_compile")

gogofaster_grpc_compile(
    name = "greeter_gogo_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def gogofaster_grpc_compile(**kwargs):
    # If importpath specified, declare a custom plugin that should correctly
    # predict the output location.
    importpath = kwargs.get("importpath")
    if importpath and not kwargs.get("plugins"):
        name = kwargs.get("name")
        name_plugin = name + "_plugin"
        proto_plugin(
            name = name_plugin,
            outputs = ["{package}/%s/{basename}.pb.go" % importpath],
            tool = "@com_github_gogo_protobuf//protoc-gen-gogofaster",
        )
        kwargs["plugins"] = [name_plugin]
        kwargs.pop("importpath")
    # Define the default plugin if still not defined
    if not kwargs.get("plugins"):
        kwargs["plugins"] = [str(Label("//github.com/gogo/protobuf:gogofaster"))]

    proto_compile(
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogofaster_proto_library`

Generates gogofaster protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogofaster_proto_library")

gogofaster_proto_library()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogofaster_proto_library.bzl", "gogofaster_proto_library")

gogofaster_proto_library(
    name = "person_gogo_library",
    importpath = "github.com/stackb/rules_proto/gogo/example/gogofaster_proto_library/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
    go_deps = [
		"@com_github_gogo_protobuf//types:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)
```

### `IMPLEMENTATION`

```python
load("//github.com/gogo/protobuf:gogofaster_proto_compile.bzl", "gogofaster_proto_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
	"google/protobuf/any.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def gogofaster_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    gogofaster_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})) + wkt_mappings,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
        ],
        importpath = importpath,
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

---

## `gogofaster_grpc_library`

Generates gogofaster protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogofaster_grpc_library")

gogofaster_grpc_library()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:gogofaster_grpc_library.bzl", "gogofaster_grpc_library")

gogofaster_grpc_library(
    name = "greeter_gogo_library",
    importpath = "github.com/stackb/rules_proto/gogo/example/gogofaster_grpc_library/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
)
```

### `IMPLEMENTATION`

```python
load("//github.com/gogo/protobuf:gogofaster_grpc_compile.bzl", "gogofaster_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("//go:utils.bzl", "get_importmappings")

wkt_mappings = get_importmappings({
	"google/protobuf/any.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/duration.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/struct.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/timestamp.proto": "github.com/gogo/protobuf/types",
	"google/protobuf/wrappers.proto": "github.com/gogo/protobuf/types",
})

def gogofaster_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")
    go_deps = kwargs.get("go_deps", [])

    name_pb = name + "_pb"

    gogofaster_grpc_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        plugin_options = get_importmappings(kwargs.pop("importmap", {})) + wkt_mappings,
        visibility = visibility,
    )

    go_library(
        name = name,
        srcs = [name_pb],
        deps = go_deps + [
            "@com_github_gogo_protobuf//proto:go_default_library",
            "@org_golang_google_grpc//:go_default_library",
            "@org_golang_x_net//context:go_default_library",
        ],
        importpath = importpath,
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
| importpath   | `string` | `None`    | Importpath for the generated artifacts          |
| importmap   | `string_dict` | `None`    | A dictionary of the form `{ K: V}` that dictates the importpath `V` for a matching imported proto file `K`          |

