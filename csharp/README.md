# `csharp`

**NOTE 1**: the csharp_* rules currently don't play nicely with sandboxing.  You may see errors like:

~~~python
The user's home directory could not be determined. Set the 'DOTNET_CLI_HOME' environment variable to specify the directory to use.
~~~

or

~~~python
System.ArgumentNullException: Value cannot be null.
Parameter name: path1
   at System.IO.Path.Combine(String path1, String path2)
   at Microsoft.DotNet.Configurer.CliFallbackFolderPathCalculator.get_DotnetUserProfileFolderPath()
   at Microsoft.DotNet.Configurer.FirstTimeUseNoticeSentinel..ctor(CliFallbackFolderPathCalculator cliFallbackFolderPathCalculator)
   at Microsoft.DotNet.Cli.Program.ProcessArgs(String[] args, ITelemetry telemetryClient)
   at Microsoft.DotNet.Cli.Program.Main(String[] args)
~~~

To remedy this, use --strategy=CoreCompile=standalone for the csharp rules (put it in your .bazelrc file).

**NOTE 2**: the csharp nuget dependency sha256 values do not appear stable.

| Rule | Description |
| ---: | :--- |
| [csharp_proto_compile](#csharp_proto_compile) | Generates csharp protobuf artifacts |
| [csharp_grpc_compile](#csharp_grpc_compile) | Generates csharp protobuf+gRPC artifacts |
| [csharp_proto_library](#csharp_proto_library) | Generates csharp protobuf library |
| [csharp_grpc_library](#csharp_grpc_library) | Generates csharp protobuf+gRPC library |

---

## `csharp_proto_compile`

Generates csharp protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//csharp:deps.bzl", "csharp_proto_compile")

csharp_proto_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//csharp:csharp_proto_compile.bzl", "csharp_proto_compile")

csharp_proto_compile(
    name = "person_csharp_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")

def csharp_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//csharp:csharp")),
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

## `csharp_grpc_compile`

Generates csharp protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//csharp:deps.bzl", "csharp_grpc_compile")

csharp_grpc_compile()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//csharp:csharp_grpc_compile.bzl", "csharp_grpc_compile")

csharp_grpc_compile(
    name = "greeter_csharp_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")

def csharp_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//csharp:csharp")),
            str(Label("//csharp:grpc_csharp")),
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

## `csharp_proto_library`

Generates csharp protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//csharp:deps.bzl", "csharp_proto_library")

csharp_proto_library()

load("@io_bazel_rules_dotnet//dotnet:defs.bzl", 
    "dotnet_register_toolchains", 
    "core_register_sdk",
    "dotnet_repositories",
)

core_version = "v2.1.503"

dotnet_register_toolchains(
    core_version = core_version,
)

dotnet_register_toolchains(
    core_version = core_version,
)

core_register_sdk(
    name = "core_sdk",
    core_version = core_version
)

dotnet_repositories()

load("@build_stack_rules_proto//csharp/nuget:packages.bzl", nuget_packages = "packages")

nuget_packages()

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_protobuf_packages")

nuget_protobuf_packages()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//csharp:csharp_proto_library.bzl", "csharp_proto_library")

csharp_proto_library(
    name = "person_csharp_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//csharp:csharp_proto_compile.bzl", "csharp_proto_compile")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def csharp_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")
    transitive = kwargs.get("transitive")

    name_pb = name + "_pb"
    csharp_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        transitive = transitive,
        verbose = verbose,
    )

    core_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "@google.protobuf//:netstandard1.0_core",
            "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
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

## `csharp_grpc_library`

Generates csharp protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//csharp:deps.bzl", "csharp_grpc_library")

csharp_grpc_library()

load("@io_bazel_rules_dotnet//dotnet:defs.bzl", 
    "dotnet_register_toolchains", 
    "core_register_sdk",
    "dotnet_repositories",
)

core_version = "v2.1.503"

dotnet_register_toolchains(
    core_version = core_version,
)

dotnet_register_toolchains(
    core_version = core_version,
)

core_register_sdk(
    name = "core_sdk",
    core_version = core_version
)

dotnet_repositories()

load("@build_stack_rules_proto//csharp/nuget:packages.bzl", nuget_packages = "packages")

nuget_packages()

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_protobuf_packages")

nuget_protobuf_packages()

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_grpc_packages")

nuget_grpc_packages()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//csharp:csharp_grpc_library.bzl", "csharp_grpc_library")

csharp_grpc_library(
    name = "greeter_csharp_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//csharp:csharp_grpc_compile.bzl", "csharp_grpc_compile")
load("@io_bazel_rules_dotnet//dotnet:defs.bzl", "core_library")

def csharp_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")
    transitive = kwargs.get("transitive")

    name_pb = name + "_pb"
    csharp_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        transitive = transitive,
        verbose = verbose,
    )

    core_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "@google.protobuf//:netstandard1.0_core",
            "@io_bazel_rules_dotnet//dotnet/stdlib.core:system.io.dll",
            "@grpc.core//:netstandard1.5_core",
            "@system.interactive.async//:netstandard2.0_core",
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

