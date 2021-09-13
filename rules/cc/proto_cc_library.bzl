"proto_cc_library.bzl provides a cc_library for proto files."

load("@rules_cc//cc:defs.bzl", "cc_library")

def proto_cc_library(**kwargs):
    cc_library(**kwargs)
