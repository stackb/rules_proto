load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")

load(
    "//:deps.bzl",
    "com_github_scalapb_scalapb",
    "io_bazel_rules_go",
    "io_bazel_rules_scala",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def com_thesamet_scalapb_scalapb_json4s(**kwargs):
    if "com_thesamet_scalapb_scalapb_json4s" not in native.existing_rules():
        jvm_maven_import_external(
            name = "com_thesamet_scalapb_scalapb_json4s",
            artifact = "com.thesamet.scalapb:scalapb-json4s_2.12:0.7.1",
            server_urls = ["https://central.maven.org/maven2"],
            artifact_sha256 = "6c8771714329464e03104b6851bfdc3e2e4967276e1a9bd2c87c3b5a6d9c53c7",
        )

def org_json4s_json4s_jackson_2_12(**kwargs):
    if "org_json4s_json4s_jackson_2_12" not in native.existing_rules():
        jvm_maven_import_external(
            name = "org_json4s_json4s_jackson_2_12",
            artifact = "org.json4s:json4s-jackson_2.12:3.6.1",
            server_urls = ["https://central.maven.org/maven2"],
            artifact_sha256 = "83b854a39e69f022ad3d7dd3da664623252dc822ed4ed1117304f39115c88043",
        )

def org_json4s_json4s_core_2_12(**kwargs):
    if "org_json4s_json4s_core_2_12" not in native.existing_rules():
        jvm_maven_import_external(
            name = "org_json4s_json4s_core_2_12",
            artifact = "org.json4s:json4s-core_2.12:3.6.1",
            server_urls = ["https://central.maven.org/maven2"],
            artifact_sha256 = "e0f481509429a24e295b30ba64f567bad95e8d978d0882ec74e6dab291fcdac0",
        )

def org_json4s_json4s_ast_2_12(**kwargs):
    if "org_json4s_json4s_ast_2_12" not in native.existing_rules():
        jvm_maven_import_external(
            name = "org_json4s_json4s_ast_2_12",
            artifact = "org.json4s:json4s-ast_2.12:3.6.1",
            server_urls = ["https://central.maven.org/maven2"],
            artifact_sha256 = "39c7de601df28e32eb0c4e3d684ec65bbf2e59af83c6088cda12688d796f7746",
        )

def scala_proto_compile(**kwargs):
    protobuf(**kwargs)
    io_bazel_rules_go(**kwargs)
    io_bazel_rules_scala(**kwargs)
    com_github_scalapb_scalapb(**kwargs)

def scala_grpc_compile(**kwargs):
    scala_proto_compile(**kwargs)

def scala_proto_library(**kwargs):
    scala_proto_compile(**kwargs)

def scala_grpc_library(**kwargs):
    scala_grpc_compile(**kwargs)
    scala_proto_library(**kwargs)

    # This one actually only needed for routeguide example
    com_thesamet_scalapb_scalapb_json4s(**kwargs)
    org_json4s_json4s_core_2_12(**kwargs)
    org_json4s_json4s_jackson_2_12(**kwargs)
    org_json4s_json4s_ast_2_12(**kwargs)
