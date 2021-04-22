package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// ProtoLibrary represents a proto_library targets and its associated parse
// .proto files.
type ProtoLibrary interface {
	// Name returns the name of the rule (e.g. foo_proto)
	Name() string

	// BaseName returns the name of the rule (e.g. foo).  This is typically
	// derived from the proto file package or name.
	BaseName() string

	// Rule returns the underlying rule
	Rule() *rule.Rule

	// Deps lists all direct library dependencies.
	Deps() []string

	// Lists all files references in the rule.
	Files() []*ProtoFile
}
