"proto_java_library.bzl provides a java_library for proto files."

load("@rules_java//java:java_library.bzl", "java_library")

def proto_java_library(**kwargs):
    java_library(**kwargs)
