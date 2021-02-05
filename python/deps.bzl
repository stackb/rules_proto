load("//:deps.bzl", "com_google_protobuf")

def py_proto_library(**kwargs):
    com_google_protobuf()
