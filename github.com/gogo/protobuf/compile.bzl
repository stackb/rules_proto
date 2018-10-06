load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")


def gogo_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:gogo"))],
        **kwargs
    )

def gogo_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_gogo"))],
        **kwargs
    )


def gogoslick_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:gogoslick"))],
        **kwargs
    )

def gogoslick_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_gogoslick"))],
        **kwargs
    )


def gogofast_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:gogofast"))],
        **kwargs
    )

def gogofast_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_gogofast"))],
        **kwargs
    )


def gogofaster_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:gogofaster"))],
        **kwargs
    )

def gogofaster_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_gogofaster"))],
        **kwargs
    )


def gogotypes_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:gogotypes"))],
        **kwargs
    )

def gogotypes_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_gogotypes"))],
        **kwargs
    )
