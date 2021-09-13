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
	// StripImportPrefix returns the strip_import_prefix or the empty string.
	StripImportPrefix() string
	// Files returns the list of proto files in the rule.
	Files() []*File
}
