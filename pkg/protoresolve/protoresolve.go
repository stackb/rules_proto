package protoresolve

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
)

const (
	// ProtoResolve should be used by private attrs
	// private attr for later deps resolution.
	ProtoResolveKey = "_proto_resolve"
)

var GlobalResolver *Index

// Index knows how to read a protoresolve file and can perform proto import cross
// resolution.
type Index struct {
	// knownProtoImports stores a mapping like
	// 'github.com/protocolbuffers/protobuf/google/protobuf/any' ->
	// @com_google_protobuf//:any_go_proto
	knownGoImports map[string]label.Label
	// knownProtoImports stores a mapping like 'google/protobuf/any.proto' ->
	// @com_google_protobuf//:any_proto
	knownProtoImports map[string]label.Label
	// known imports is a mapping such that the top-level key is a ruleKind
	// name, and the second key is a proto import statement.
	knownImports map[string]map[string]label.Label
	entries      []*IndexEntry
}

func NewIndex() *Index {
	return &Index{
		entries:           make([]*IndexEntry, 0),
		knownImports:      make(map[string]map[string]label.Label),
		knownGoImports:    make(map[string]label.Label),
		knownProtoImports: make(map[string]label.Label),
	}
}

// ParseCSVFile reads a protoresolve csv file.
func (idx *Index) ParseCSVFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return idx.ParseCSVReader(f)
}

// ParseCSVFile reads input and returns a list of items.  Comment lines beginning
// with '#' are ignored.
func (idx *Index) ParseCSVReader(in io.Reader) error {
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ",", 4)
		if len(parts) == 4 {
			l, err := label.Parse(parts[0])
			if err != nil {
				return fmt.Errorf("malformed label %s: %v", parts[1], err)
			}
			idx.AddEntry(&IndexEntry{l, parts[1], parts[2], parts[3]})
		}
	}

	return nil
}

func (idx *Index) GetEntriesByKind(kind string) []*IndexEntry {
	entries := make([]*IndexEntry, 0)
	for _, e := range idx.entries {
		if e.Kind == kind {
			entries = append(entries, e)
		}
	}
	return entries
}

func (idx *Index) AddEntry(e *IndexEntry) {
	idx.entries = append(idx.entries, e)
	if e.Kind == "proto_library" && e.Attr == "srcs" {
		idx.knownProtoImports[e.Value] = e.Label
	}
	if e.Attr == "importpath" {
		idx.knownGoImports[e.Value] = e.Label
	}
	if e.Attr == "kind" {
		idx.knownGoImports[e.Value] = e.Label
	}
}

func (idx *Index) AddKind(kind, imp string, l label.Label) {
	known, ok := idx.knownImports[kind]
	if !ok {
		known = make(map[string]label.Label)
		idx.knownImports[kind] = known
	}
	known[imp] = l
}

// CrossResolve provides dependency resolution logic for the proto language extension.
func (idx *Index) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	switch lang {
	case "go":
		return idx.resolveGo(c, ix, imp)
	case "proto":
		return idx.resolveProto(c, ix, imp)
	default:
		return idx.resolveKind(c, ix, imp, lang)
	}
	return nil
}

func (idx *Index) resolveKind(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, kind string) []resolve.FindResult {
	known := idx.knownImports[kind]
	if known == nil {
		return nil
	}
	if l, ok := known[imp.Imp]; ok {
		// log.Println("Cross-resolved", imp.Imp, l.String())
		return []resolve.FindResult{{Label: l}}
	}
	return nil
}

func (idx *Index) resolveGo(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec) []resolve.FindResult {
	if l, ok := idx.knownGoImports[imp.Imp]; ok {
		// log.Println("Cross-resolved", imp.Imp, l.String())
		return []resolve.FindResult{{Label: l}}
	}
	return nil
}

func (idx *Index) resolveProto(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec) []resolve.FindResult {
	if l, ok := idx.knownProtoImports[imp.Imp]; ok {
		// log.Println("Cross-resolved", imp.Imp, l.String())
		return []resolve.FindResult{{Label: l}}
	}
	return nil
}

// IndexEntry represents a line in an protoresolve csv file.
type IndexEntry struct {
	// Label is the rule Label
	Label label.Label
	// Kind is the name of a rule (e.g. proto_library)
	Kind string
	// Attr is the name of an attr (e.g. "srcs")
	Attr string
	// Value is the name of the attr
	Value string
}

// String implements the Stringer interface
func (e *IndexEntry) String() string {
	return fmt.Sprintf("%s,%s,%s,%s", e.Label.String(), e.Kind, e.Attr, e.Value)
}
