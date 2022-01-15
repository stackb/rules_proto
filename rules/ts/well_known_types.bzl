load("@build_stack_rules_proto//rules/ts:proto_ts_library.bzl", "proto_ts_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

WELL_KNOWN_PROTO_MAP = {
    "any": ("google/protobuf/any.ts", []),
    "api": (
        "google/protobuf/api.ts",
        [
            "source_context",
            "type",
        ],
    ),
    "compiler_plugin": (
        "google/protobuf/compiler/plugin.ts",
        ["descriptor"],
    ),
    "descriptor": ("google/protobuf/descriptor.ts", []),
    "duration": ("google/protobuf/duration.ts", []),
    "empty": ("google/protobuf/empty.ts", []),
    "field_mask": ("google/protobuf/field_mask.ts", []),
    "source_context": ("google/protobuf/source_context.ts", []),
    "struct": ("google/protobuf/struct.ts", []),
    "timestamp": ("google/protobuf/timestamp.ts", []),
    "type": (
        "google/protobuf/type.ts",
        [
            "any",
            "source_context",
        ],
    ),
    "wrappers": ("google/protobuf/wrappers.ts", []),
}

def wkt_ts_proto_compile(**kwargs):
    [proto_compile(
        name = "wkt_" + proto[0] + "_ts_proto_compile",
        options = {"@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto": [
            "emitImportedFiles=false",
            "esModuleInterop=true",
        ]},
        outputs = [
            proto[1][0],
        ],
        plugins = ["@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto"],
        proto = "@com_google_protobuf//:" + proto[0] + "_proto",
        **kwargs,
    ) for proto in WELL_KNOWN_PROTO_MAP.items()]

def wkt_ts_proto(**kwargs):
    _deps = [
        "@npm_tsc//long",
        "@npm_tsc//protobufjs",
    ]

    [proto_ts_library(
        name = "wkt_" + proto[0] + "_ts_proto",
        srcs = [
            proto[1][0],
        ],
        args = [
            "--allowSyntheticDefaultImports",
            "--downlevelIteration",
            "--lib ES2015",
        ],
        tsc = "@npm_tsc//typescript/bin:tsc",
        deps = _deps + ["wkt_" + dep + "_ts_proto" for dep in proto[1][1]],
        **kwargs,
    ) for proto in WELL_KNOWN_PROTO_MAP.items()]
