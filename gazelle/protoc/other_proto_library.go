package protoc

import (
	"fmt"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
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
	matching := make([]*ProtoFile, 0)
	return matching
}

// Deps implements part of the ProtoLibrary interface
func (s *OtherProtoLibrary) Deps() []string {
	return s.rule.AttrStrings("deps")
}

// Srcs returns the srcs attribute
func (s *OtherProtoLibrary) Srcs() []string {
	return s.rule.AttrStrings("deps")
}

func matchingFiles(files map[string]*ProtoFile, srcs []label.Label) []*ProtoFile {
	matching := make([]*ProtoFile, 0)
	for _, src := range srcs {
		if file, ok := files[src.Name]; ok {
			matching = append(matching, file)
		}
	}
	return matching
}
