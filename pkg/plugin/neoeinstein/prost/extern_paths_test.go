package prost_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/prost"
	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

// makeLibraryRule constructs a proto_library rule with the given srcs and a
// ProtoLibraryKey private attr backed by a stub ProtoLibrary so that
// ResolveTransitiveExternPaths can read it.
func makeLibraryRule(name, pkg string, srcs []string) *rule.Rule {
	r := rule.NewRule("proto_library", name)
	r.SetAttr("srcs", srcs)
	files := make([]*protoc.File, len(srcs))
	for i, s := range srcs {
		files[i] = protoc.NewFile(pkg, s)
	}
	lib := protoc.NewOtherProtoLibrary(nil, r, files...)
	r.SetPrivateAttr(protoc.ProtoLibraryKey, lib)
	return r
}

func TestResolveTransitiveExternPaths(t *testing.T) {
	resolver := protoc.GlobalResolver()

	// Register prost_extern entries for two upstream libraries.
	resolver.Provide("proto", "prost_extern",
		"externtest/depA/a.proto",
		label.New("", "extern.dep_a", "depA_rs"))
	resolver.Provide("proto", "prost_extern",
		"externtest/depB/b.proto",
		label.New("", "extern.dep_b", "depB_rs"))

	// Set up the depends graph: own.proto -> depA -> depB, plus a WKT skip.
	resolver.Provide("proto", "depends",
		"externtest/own/own.proto",
		label.New("", "externtest/depA", "a.proto"))
	resolver.Provide("proto", "depends",
		"externtest/own/own.proto",
		label.New("", "google/protobuf", "duration.proto"))
	resolver.Provide("proto", "depends",
		"externtest/depA/a.proto",
		label.New("", "externtest/depB", "b.proto"))

	r := makeLibraryRule("own_proto", "externtest/own", []string{"own.proto"})

	from := label.New("", "externtest/own", "own_proto")
	got := prost.ResolveTransitiveExternPaths(r, from)
	sort.Strings(got)

	want := []string{
		"extern_path=.extern.dep_a=::depA_rs::extern::dep_a",
		"extern_path=.extern.dep_b=::depB_rs::extern::dep_b",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveTransitiveExternPaths:\n got: %v\nwant: %v", got, want)
	}

	// Second call should hit the cache and return the same slice.
	got2 := prost.ResolveTransitiveExternPaths(r, from)
	if !reflect.DeepEqual(got2, got) {
		t.Errorf("cached call differs:\n got: %v\nwant: %v", got2, got)
	}
}

func TestResolveTransitiveExternPaths_OwnFilesSkipped(t *testing.T) {
	resolver := protoc.GlobalResolver()

	// Register the library's own proto file as if it had been registered.
	// The function must NOT include own files in the result.
	resolver.Provide("proto", "prost_extern",
		"selftest/me/m.proto",
		label.New("", "selftest.me", "me_rs"))

	r := makeLibraryRule("me_proto", "selftest/me", []string{"m.proto"})
	from := label.New("", "selftest/me", "me_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	if len(got) != 0 {
		t.Errorf("expected empty extern paths for own files, got %v", got)
	}
}

// TestResolveTransitiveExternPaths_SubpackageOfImportSkipped is a regression
// test for the case where the current library's proto package is a
// sub-package of an imported library's proto package. Because prost's
// extern_path matches by package prefix, emitting an extern_path for the
// parent package would cause prost to rewrite the current library's own type
// references into the imported crate. The function must filter such entries
// out.
func TestResolveTransitiveExternPaths_SubpackageOfImportSkipped(t *testing.T) {
	resolver := protoc.GlobalResolver()

	// Imported library: parent proto package "subpkg.parent".
	resolver.Provide("proto", "prost_extern",
		"subpkg/parent/p.proto",
		label.New("", "subpkg.parent", "parent_rs"))

	// Current library: own proto package "subpkg.parent.child" (a sub-package
	// of the imported one) registered against its own file.
	resolver.Provide("proto", "prost_extern",
		"subpkg/parent/child/c.proto",
		label.New("", "subpkg.parent.child", "child_rs"))

	// Dependency edge: child.proto -> parent.proto.
	resolver.Provide("proto", "depends",
		"subpkg/parent/child/c.proto",
		label.New("", "subpkg/parent", "p.proto"))

	r := makeLibraryRule("child_proto", "subpkg/parent/child", []string{"c.proto"})
	from := label.New("", "subpkg/parent/child", "child_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	for _, opt := range got {
		if opt == "extern_path=.subpkg.parent=::parent_rs::subpkg::parent" {
			t.Errorf("extern_path for parent of own package was not filtered: %v", got)
		}
	}
	if len(got) != 0 {
		t.Errorf("expected no extern_paths (only entry would shadow own package), got %v", got)
	}
}

// TestResolveTransitiveExternPaths_SiblingNotFiltered ensures the filter is
// not over-aggressive: a sibling package (one that shares a common prefix but
// is neither equal to nor an ancestor of the current package) must still
// produce an extern_path entry.
func TestResolveTransitiveExternPaths_SiblingNotFiltered(t *testing.T) {
	resolver := protoc.GlobalResolver()

	// Sibling package "sibling.a.x" — shares prefix "sibling.a" with our own
	// "sibling.a.y" but neither is a parent of the other.
	resolver.Provide("proto", "prost_extern",
		"sibling/a/x/x.proto",
		label.New("", "sibling.a.x", "x_rs"))

	// Own package "sibling.a.y".
	resolver.Provide("proto", "prost_extern",
		"sibling/a/y/y.proto",
		label.New("", "sibling.a.y", "y_rs"))

	resolver.Provide("proto", "depends",
		"sibling/a/y/y.proto",
		label.New("", "sibling/a/x", "x.proto"))

	r := makeLibraryRule("y_proto", "sibling/a/y", []string{"y.proto"})
	from := label.New("", "sibling/a/y", "y_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	want := []string{"extern_path=.sibling.a.x=::x_rs::sibling::a::x"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveTransitiveExternPaths:\n got: %v\nwant: %v", got, want)
	}
}

func TestResolveExternPathOptions_FiltersExisting(t *testing.T) {
	// Library with no transitive deps — extern paths come only from cfg.Options
	// after filtering out any pre-existing extern_path= entries.
	r := makeLibraryRule("noop_proto", "exfilter/noop", []string{"n.proto"})
	from := label.New("", "exfilter/noop", "noop_proto")

	cfg := &protoc.PluginConfiguration{
		Options: []string{
			"compile_well_known_types=true",
			"extern_path=.stale.pkg=::stale_rs::stale::pkg",
		},
	}

	got := prost.ResolveExternPathOptions(cfg, r, from)
	for _, opt := range got {
		if opt == "extern_path=.stale.pkg=::stale_rs::stale::pkg" {
			t.Errorf("stale extern_path option was not filtered: %v", got)
		}
	}

	want := []string{"compile_well_known_types=true"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveExternPathOptions:\n got: %v\nwant: %v", got, want)
	}
}
