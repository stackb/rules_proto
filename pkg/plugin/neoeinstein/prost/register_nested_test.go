package prost

import (
	"sort"
	"strings"
	"testing"

	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

func TestRegisterNestedTypes_RegistersNestedPaths(t *testing.T) {
	resolver := protoc.NewImportResolver(&protoc.ImportResolverOptions{})

	const src = `syntax = "proto2";
package nested.test;
message Outer {
    message Inner {
        message Leaf {}
    }
    enum NestedEnum { A = 0; }
}
message TopLevel {}
`
	f := protoc.NewFile("nested/test", "x.proto")
	if err := f.ParseReader(strings.NewReader(src)); err != nil {
		t.Fatalf("parse: %v", err)
	}

	for _, msg := range f.Messages() {
		registerNestedTypes(resolver, "nested.test", msg, msg.Name, "perfile_crate")
	}

	var keys []string
	for _, ent := range resolver.Resolve("proto", PerFileTypeProvideKind, "nested.test") {
		keys = append(keys, ent.Label.Pkg)
	}
	sort.Strings(keys)

	for _, want := range []string{
		"Outer",
		"Outer.Inner",
		"Outer.Inner.Leaf",
		"Outer.NestedEnum",
		"TopLevel",
	} {
		found := false
		for _, k := range keys {
			if k == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing nested registration %q; got: %v", want, keys)
		}
	}
}
