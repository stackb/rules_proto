package grpc_test

import (
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

// TestMain is failing in CI:
// error: mkdir /home/runner/.cache/bazel/_bazel_runner/43914acf8bee29fb1d82f3e3dee22a49/sandbox/linux-sandbox/789/bazel_testing: read-only file system
// FIXME(pcj)
func SkipTestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: txtar,
	})
}

func SkipTestBuild(t *testing.T) {
	if err := bazel_testing.RunBazel("build", ":all"); err != nil {
		t.Fatal(err)
	}
}

var txtar = `
-- WORKSPACE --
local_repository(
    name = "build_stack_rules_proto",
    path = "../../build_stack_rules_proto",
)

-- BUILD.bazel --

# gazelle:proto_language builtins rule proto_compile
# gazelle:proto_language builtins plugin grpc:grpc:protoc-gen-grpc-python

-- foo.proto --
syntax = "proto3";

message FooRequest {
    string name = 1;
}

message Foo {
    string name = 1;
}

service Fooer {
	rpc GetFoo(FooRequest) returns (Foo) {}
}
`
