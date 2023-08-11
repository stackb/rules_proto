load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")
load("@build_stack_rules_proto//rules/ts:proto_ts_library.bzl", "proto_ts_library")

proto_library(
    name = "timestamppb_proto",
    srcs = ["timestamp.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "timestamp_ts_proto_compile",
    options = {"@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto": [
        "emitImportedFiles=false",
    ]},
    outputs = ["timestamp.ts"],
    plugins = ["@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto"],
    proto = "timestamppb_proto",
)

proto_ts_library(
    name = "timestamp_ts_proto",
    srcs = ["timestamp.ts"],
    declaration = True,
    tags = ["manual"],
    tsconfig = {
        "compilerOptions": {
            "lib": ["es2015"],
            "declaration": True,
        },
    },
    visibility = ["//visibility:public"],
    deps = [
        "//:node_modules/long",
        "//:node_modules/protobufjs",
    ],
)
