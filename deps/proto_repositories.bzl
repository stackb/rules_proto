"""
Third-party proto dependencies
"""

load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository")

def proto_repositories():
    """third party proto repositories
    """
    proto_repository(
        name = "protobufapis",
        build_directives = [
            "gazelle:proto_language go enable true",
            "gazelle:exclude google/protobuf/compiler/ruby/**",
            "gazelle:exclude google/protobuf/compiler/cpp/**",
            "gazelle:exclude google/protobuf/util/**",
            "gazelle:exclude google/protobuf/unittest/**",
        ],
        deleted_files = [
            "google/protobuf/any_test.proto",
            "google/protobuf/map_lite_unittest.proto",
            "google/protobuf/map_proto2_unittest.proto",
            "google/protobuf/map_unittest.proto",
            "google/protobuf/test_messages_proto2.proto",
            "google/protobuf/test_messages_proto3.proto",
            "google/protobuf/unittest_arena.proto",
            "google/protobuf/unittest_custom_options.proto",
            "google/protobuf/unittest_drop_unknown_fields.proto",
            "google/protobuf/unittest_embed_optimize_for.proto",
            "google/protobuf/unittest_empty.proto",
            "google/protobuf/unittest_enormous_descriptor.proto",
            "google/protobuf/unittest_import_lite.proto",
            "google/protobuf/unittest_import_public_lite.proto",
            "google/protobuf/unittest_import_public.proto",
            "google/protobuf/unittest_import.proto",
            "google/protobuf/unittest_lazy_dependencies_custom_option.proto",
            "google/protobuf/unittest_lazy_dependencies_enum.proto",
            "google/protobuf/unittest_lazy_dependencies.proto",
            "google/protobuf/unittest_lite_imports_nonlite.proto",
            "google/protobuf/unittest_lite.proto",
            "google/protobuf/unittest_mset_wire_format.proto",
            "google/protobuf/unittest_mset.proto",
            "google/protobuf/unittest_no_field_presence.proto",
            "google/protobuf/unittest_no_generic_services.proto",
            "google/protobuf/unittest_optimize_for.proto",
            "google/protobuf/unittest_preserve_unknown_enum.proto",
            "google/protobuf/unittest_preserve_unknown_enum2.proto",
            "google/protobuf/unittest_proto3_arena_lite.proto",
            "google/protobuf/unittest_proto3_arena.proto",
            "google/protobuf/unittest_proto3_lite.proto",
            "google/protobuf/unittest_proto3_optional.proto",
            "google/protobuf/unittest_proto3.proto",
            "google/protobuf/unittest_well_known_types.proto",
            "google/protobuf/unittest.proto",
        ],
        build_file_expunge = True,
        build_file_proto_mode = "file",
        cfgs = ["//example:config.yaml"],
        sha256 = "bf0e5070b4b99240183b29df78155eee335885e53a8af8683964579c214ad301",
        strip_prefix = "protobuf-3.14.0/src",
        urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.14.0.zip"],
    )

    proto_repository(
        name = "googleapis",
        build_directives = [
            "gazelle:proto_language go enable true",
        ],
        build_file_generation = "on",
        build_file_proto_mode = "file",
        cfgs = ["//example:config.yaml"],
        override_go_googleapis = True,
        sha256 = "b9dbc65ebc738a486265ef7b708e9449bf361541890091983e946557ee0a4bfc",
        strip_prefix = "googleapis-66759bdf6a5ebb898c2a51c8649aefd1ee0b7ffe",
        type = "zip",
        urls = ["https://codeload.github.com/googleapis/googleapis/zip/66759bdf6a5ebb898c2a51c8649aefd1ee0b7ffe"],
    )
