load("@rules_cc//cc:defs.bzl", "cc_library")

def proto_cc_library(**kwargs):
    cc_library(**kwargs)
