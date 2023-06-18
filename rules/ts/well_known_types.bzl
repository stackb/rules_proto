load("@build_stack_rules_proto//rules/ts:proto_ts_library.bzl", "proto_ts_library")
load("@build_stack_rules_proto//rules:proto_compile.bzl", "proto_compile")

WELL_KNOWN_PROTO_MAP = {
    "any": ("any.ts", []),
    "api": (
        "api.ts",
        [
            "source_context",
            "type",
        ],
    ),
    "compiler_plugin": (
        "compiler/plugin.ts",
        ["descriptor"],
    ),
    "descriptor": ("descriptor.ts", []),
    "duration": ("duration.ts", []),
    "empty": ("empty.ts", []),
    "field_mask": ("field_mask.ts", []),
    "source_context": ("source_context.ts", []),
    "struct": ("struct.ts", []),
    "timestamp": ("timestamp.ts", []),
    "type": (
        "type.ts",
        [
            "any",
            "source_context",
        ],
    ),
    "wrappers": ("wrappers.ts", []),
}

def well_known_ts_proto_compile(**kwargs):
    for proto in WELL_KNOWN_PROTO_MAP.items():
        proto_compile(
            name = "wkt_" + proto[0] + "_proto_compile",
            options = {"@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto": [
                "emitImportedFiles=false",
                "esModuleInterop=true",
            ]},
            outputs = [
                proto[1][0],
            ],
            plugins = ["@build_stack_rules_proto//plugin/stephenh/ts-proto:protoc-gen-ts-proto"],
            proto = "@com_google_protobuf//:" + proto[0] + "_proto",
            **kwargs
        )

def well_known_ts_proto_library(deps, **kwargs):
    for proto in WELL_KNOWN_PROTO_MAP.items():
        proto_ts_library(
            name = "wkt_" + proto[0] + "_ts_proto",
            srcs = [
                proto[1][0],
            ],
            deps = deps + ["wkt_" + dep + "_ts_proto" for dep in proto[1][1]],
            **kwargs
        )
