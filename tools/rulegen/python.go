package main

var pythonGrpcLibraryWorkspaceTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@io_bazel_rules_python//python:pip.bzl", "pip_import", "pip_repositories")

pip_repositories()

pip_import(
    name = "grpc_py_deps",
    requirements = "@build_stack_rules_proto//python/requirements:grpc.txt",
)

load("@grpc_py_deps//:requirements.bzl", grpc_pip_install = "pip_install")

grpc_pip_install()`)

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
        deps = [
            "@com_google_protobuf//:protobuf_python",
        ],
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
    )

# Alias
py_proto_library = python_proto_library`)

var pythonGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@grpc_py_deps//:requirements.bzl", grpc_requirements = "all_requirements")

def python_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create {{ .Lang.Name }} library
    native.py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_google_protobuf//:protobuf_python",
        ] + grpc_requirements,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
    )

# Alias
py_grpc_library = python_grpc_library`)

func makePython() *Language {
	return &Language{
		Dir:   "python",
		Name:  "python",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "python_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//python:python"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates *.py protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//python:python", "//python:grpc_python"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates *.py protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_proto_library",
				Kind:             "proto",
				Implementation:   pythonProtoLibraryRuleTemplate,
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates *.py protobuf library",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_grpc_library",
				Kind:             "grpc",
				Implementation:   pythonGrpcLibraryRuleTemplate,
				WorkspaceExample: pythonGrpcLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates *.py protobuf+gRPC library",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
