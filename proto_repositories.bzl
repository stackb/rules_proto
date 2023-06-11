load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository")

def proto_repositories():
    proto_repository(
        name = "protoapis",
        build_directives = [
            "gazelle:exclude testdata",
            "gazelle:exclude google/protobuf/compiler/ruby",
            "gazelle:exclude google/protobuf/util",
            "gazelle:proto_language go enable false",
            "gazelle:proto_language descriptor enable true",
        ],
        build_file_expunge = True,
        build_file_proto_mode = "file",
        cfgs = ["//example:config.yaml"],
        deleted_files = [
            # some of these are unparseable with curret parser version
            "google/protobuf/unittest_custom_options.proto",
            "google/protobuf/unittest.proto",
            "google/protobuf/unittest_lite.proto",
            "google/protobuf/unittest_mset_wire_format.proto",
            "google/protobuf/unittest_retention.proto",
            "google/protobuf/map_unittest.proto",
            "google/protobuf/map_lite_unittest.proto",
            "google/protobuf/map_proto2_unittest.proto",
            "google/protobuf/test_messages_proto2.proto",
            "google/protobuf/test_messages_proto3.proto",
        ],
        strip_prefix = "protobuf-a74f54b724bdc2fe0bfc271f4dc0ceb159805625/src",
        type = "zip",
        # https://github.com/protocolbuffers/protobuf/releases/tag/v23.2
        urls = ["https://codeload.github.com/protocolbuffers/protobuf/zip/a74f54b724bdc2fe0bfc271f4dc0ceb159805625"],
    )
