package main

var pythonGrpcLibraryWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()`)

var pythonGrpclibLibraryWorkspaceTemplate = mustTemplate(`load("@rules_proto_grpc//{{ .Lang.Dir }}:repositories.bzl", rules_proto_grpc_{{ .Lang.Name }}_repos="{{ .Lang.Name }}_repos")

rules_proto_grpc_{{ .Lang.Name }}_repos()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@rules_python//python:repositories.bzl", "py_repositories")
py_repositories()

load("@rules_python//python:pip.bzl", "pip_install")
pip_install(
    name = "rules_proto_grpc_py3_deps",
    python_interpreter = "python3",
    requirements = "@rules_proto_grpc//python:requirements.txt",
)`)

var pythonProtoLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@rules_python//python:defs.bzl", "py_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create {{ .Lang.Name }} library
    py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

PROTO_DEPS = [
    "@com_google_protobuf//:protobuf_python",
]`)

var pythonGrpcLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_{{ .Rule.Kind }}_compile.bzl", "{{ .Lang.Name }}_{{ .Rule.Kind }}_compile")
load("@rules_python//python:defs.bzl", "py_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create {{ .Lang.Name }} library
    py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    "@com_google_protobuf//:protobuf_python",
    "@com_github_grpc_grpc//src/python/grpcio/grpc:grpcio",
]`)

var pythonGrpclibLibraryRuleTemplate = mustTemplate(`load("//{{ .Lang.Dir }}:{{ .Lang.Name }}_grpclib_compile.bzl", "{{ .Lang.Name }}_grpclib_compile")
load("@rules_python//python:defs.bzl", "py_library")

def {{ .Rule.Name }}(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_grpclib_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create {{ .Lang.Name }} library
    py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_google_protobuf//:protobuf_python",
        ] + GRPC_DEPS,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags"),
    )

GRPC_DEPS = [
    # Don't use requirement(), since rules_proto_grpc_py3_deps doesn't necessarily exist when
    # imported by defs.bzl
    "@rules_proto_grpc_py3_deps//pypi__grpclib",
]`)

func makePython() *Language {
	return &Language{
		Dir:   "python",
		Name:  "python",
		DisplayName: "Python",
		Notes: mustTemplate("Rules for generating Python protobuf and gRPC `.py` files and libraries using standard Protocol Buffers and gRPC or [grpclib](https://github.com/vmagamedov/grpclib). Libraries are created with `py_library` from `rules_python`. To use the fast C++ Protobuf implementation, you can add `--define=use_fast_cpp_protos=true` to your build, but this requires you setup the path to your Python headers.\n\nNote: On Windows, the path to Python for `pip_install` may need updating to `Python.exe`, depending on your install."),
		Flags: commonLangFlags,
		Aliases: map[string]string{
			"py_proto_compile": "python_proto_compile",
			"py_grpc_compile": "python_grpc_compile",
			"py_grpclib_compile": "python_grpclib_compile",
			"py_proto_library": "python_proto_library",
			"py_grpc_library": "python_grpc_library",
			"py_grpclib_library": "python_grpclib_library",
		},
		Rules: []*Rule{
			&Rule{
				Name:             "python_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//python:python_plugin"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates Python protobuf `.py` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//python:python_plugin", "//python:grpc_python_plugin"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Python protobuf+gRPC `.py` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_grpclib_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//python:python_plugin", "//python:grpclib_python_plugin"},
				WorkspaceExample: pythonGrpclibLibraryWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates Python protobuf+grpclib `.py` artifacts (supports Python 3 only)",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"},
			},
			&Rule{
				Name:             "python_proto_library",
				Kind:             "proto",
				Implementation:   pythonProtoLibraryRuleTemplate,
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoLibraryExampleTemplate,
				Doc:              "Generates a Python protobuf library using `py_library` from `rules_python`",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "python_grpc_library",
				Kind:             "grpc",
				Implementation:   pythonGrpcLibraryRuleTemplate,
				WorkspaceExample: pythonGrpcLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates a Python protobuf+gRPC library using `py_library` from `rules_python`",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"},
			},
			&Rule{
				Name:             "python_grpclib_library",
				Kind:             "grpc",
				Implementation:   pythonGrpclibLibraryRuleTemplate,
				WorkspaceExample: pythonGrpclibLibraryWorkspaceTemplate,
				BuildExample:     grpcLibraryExampleTemplate,
				Doc:              "Generates a Python protobuf+grpclib library using `py_library` from `rules_python` (supports Python 3 only)",
				Attrs:            aspectProtoCompileAttrs,
				SkipTestPlatforms: []string{"windows"},
			},
		},
	}
}
