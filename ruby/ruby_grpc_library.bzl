load("//ruby:ruby_grpc_compile.bzl", "ruby_grpc_compile")
load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_library")

def ruby_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    ruby_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create ruby library
    ruby_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        includes = ["{package}/%s" % name_pb],
        visibility = kwargs.get("visibility"),
    )
