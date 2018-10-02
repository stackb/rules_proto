load("//:compile.bzl", "proto_compile")
load("//:plugin.bzl", "proto_plugin")

def gogoslick_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:slick"))],
        **kwargs
    )

def gogoslick_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_slick"))],
        **kwargs
    )

def gogofast_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:fast"))],
        **kwargs
    )

def gogofast_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_fast"))],
        **kwargs
    )

def gogofaster_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:faster"))],
        **kwargs
    )

def gogofaster_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_faster"))],
        **kwargs
    )

def gogotypes_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:types"))],
        **kwargs
    )

def gogotypes_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/gogo/protobuf:grpc_types"))],
        **kwargs
    )
