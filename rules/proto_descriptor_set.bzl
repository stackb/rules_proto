"proto_descriptor_set.bzl wraps the proto_descriptor_set rule from @rules_proto."

load("@rules_proto//proto:defs.bzl", _proto_descriptor_set = "proto_descriptor_set")

def proto_descriptor_set(**kwargs):
    _proto_descriptor_set(**kwargs)
