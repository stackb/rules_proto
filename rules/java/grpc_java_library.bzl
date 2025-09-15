"grpc_java_library.bzl provides a java_library for grpc files."

load("@rules_java//java:java_library.bzl", "java_library")

def grpc_java_library(**kwargs):
    java_library(**kwargs)
