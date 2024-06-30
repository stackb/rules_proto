"proto_nodejs_library.bzl provides a js_library for proto files."

load("@aspect_rules_js//js:defs.bzl", "js_library")

def proto_nodejs_library(**kwargs):
    js_library(**kwargs)
