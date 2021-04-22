package protoc

import (
	"fmt"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/rule"
)

// OtherProtoLibrary implements the ProtoLibrary interface from an existing ProtoLibrary rule.
type OtherProtoLibrary struct {
	rule  *rule.Rule
	files []*ProtoFile
}

// Name implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Name() string {
	return s.rule.Name()
}

// BaseName implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) BaseName() string {
	name := s.rule.Name()
	if !strings.HasSuffix(name, "_proto") {
		panic(fmt.Sprintf("Unexpected proto_library name %q (it should always end in '_proto')", name))
	}
	return name[0 : len(name)-len("_proto")]
}

// Rule implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Rule() *rule.Rule {
	return s.rule
}

// Files implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Files() []*ProtoFile {
	return s.files
}

// Deps implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Deps() []string {
	return s.rule.AttrStrings("deps")
}
