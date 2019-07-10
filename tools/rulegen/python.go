package main

var pythonGrpcLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Lang.Name }}_deps")

{{ .Lang.Name }}_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@com_apt_itude_rules_pip//rules:dependencies.bzl", "pip_rules_dependencies")

pip_rules_dependencies()

load("@com_apt_itude_rules_pip//rules:repository.bzl", "pip_repository")

pip_repository(
    name = "grpc_py2_deps",
    python_interpreter = "python2",
    requirements = "@build_stack_rules_proto//python/requirements:grpc.txt",
)

pip_repository(
    name = "grpc_py3_deps",
    python_interpreter = "python3",
    requirements = "@build_stack_rules_proto//python/requirements:grpc.txt",
)`)

var pythonProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")

def python_proto_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create {{ .Lang.Name }} library
    native.py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
    )

PROTO_DEPS = [
    "@com_google_protobuf//:protobuf_python",
]

# Alias
py_proto_library = python_proto_library`)

var pythonGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")

def python_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Pick deps based on python version
    if "python_version" not in kwargs or kwargs["python_version"] == "PY3":
        grpc_deps = GRPC_PYTHON3_DEPS
    elif kwargs["python_version"] == "PY2":
        grpc_deps = GRPC_PYTHON2_DEPS
    else:
        fail("The 'python_version' attribute to python_grpc_library must be one of ['PY2', 'PY3']")


    # Create {{ .Lang.Name }} library
    native.py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_google_protobuf//:protobuf_python",
        ] + grpc_deps,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
    )

GRPC_PYTHON2_DEPS = [
    "@grpc_py2_deps//grpcio"
]

GRPC_PYTHON3_DEPS = [
    "@grpc_py3_deps//grpcio"
]

# Alias
py_grpc_library = python_grpc_library`)

func makePython() *Language {
	return &Language{
		Dir:   "python",
		Name:  "python",
		DisplayName: "Python",
		Notes: mustTemplate("Rules for generating Python protobuf and gRPC `.py` files and libraries using standard Protocol Buffers and gRPC. Libraries are created with the Bazel native `py_library`"),
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "python_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//python:python"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates Python protobuf `.py` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//python:python", "//python:grpc_python"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Python protobuf+gRPC `.py` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_proto_library",
				Kind:             "proto",
				Implementation:   pythonProtoLibraryRuleTemplate,
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates a Python protobuf library using `py_library`",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_grpc_library",
				Kind:             "grpc",
				Implementation:   pythonGrpcLibraryRuleTemplate,
				WorkspaceExample: pythonGrpcLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates a Python protobuf+gRPC library using `py_library`",
				Attrs:            append(aspectProtoCompileAttrs, []*Attr{
					&Attr{
						Name:      "python_version",
						Type:      "string",
						Default:   "PY3",
						Doc:       "Specify the Python version to use for the bundled dependencies. Valid values are \"PY3\" (the default) and \"PY2\"",
						Mandatory: false,
					},
				}...),
				SkipTestPlatforms: []string{"windows"},
			},
		},
	}
}
