load("//ruby:ruby_proto_compile.bzl", "ruby_proto_compile")
load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_library")

def ruby_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    
    ruby_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    ruby_library(
        name = name,
        srcs = [name_pb],
        includes = ["{package}/%s" % name_pb],   
        visibility = visibility,
    )

