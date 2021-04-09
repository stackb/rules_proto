# rules for gogo
load("@build_stack_rules_proto//rules:combo_proto_compile.bzl", _combo_proto_compile = "combo_proto_compile")
load("@build_stack_rules_proto//rules:gogo_proto_compile.bzl", _gogo_proto_compile = "gogo_proto_compile")
load("@build_stack_rules_proto//rules:gogofast_proto_compile.bzl", _gogofast_proto_compile = "gogofast_proto_compile")
load("@build_stack_rules_proto//rules:gogofaster_proto_compile.bzl", _gogofaster_proto_compile = "gogofaster_proto_compile")
load("@build_stack_rules_proto//rules:gogoslick_proto_compile.bzl", _gogoslick_proto_compile = "gogoslick_proto_compile")
load("@build_stack_rules_proto//rules:gogotypes_proto_compile.bzl", _gogotypes_proto_compile = "gogotypes_proto_compile")
load("@build_stack_rules_proto//rules:gostring_proto_compile.bzl", _gostring_proto_compile = "gostring_proto_compile")

combo_proto_compile = _combo_proto_compile
gogo_proto_compile = _gogo_proto_compile
gogofast_proto_compile = _gogofast_proto_compile
gogofaster_proto_compile = _gogofaster_proto_compile
gogoslick_proto_compile = _gogoslick_proto_compile
gogotypes_proto_compile = _gogotypes_proto_compile
gostring_proto_compile = _gostring_proto_compile
