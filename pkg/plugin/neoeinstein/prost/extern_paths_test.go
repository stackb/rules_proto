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
		"extern_path=.extern.dep_a=::depA_rs",
		"extern_path=.extern.dep_b=::depB_rs",
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

// TestResolveTransitiveExternPaths_SubpackageOfImport verifies that when the
// current library's proto package is a sub-package of an imported library's
// proto package, ResolveTransitiveExternPaths emits the imported package's
// extern_path entry (this is the prost variant — no self-extern override is
// added; that's the job of ResolveExternPathOptionsForReferences).
func TestResolveTransitiveExternPaths_SubpackageOfImport(t *testing.T) {
	resolver := protoc.GlobalResolver()

	resolver.Provide("proto", "prost_extern",
		"subpkg/parent/p.proto",
		label.New("", "subpkg.parent", "parent_rs"))

	resolver.Provide("proto", "prost_extern",
		"subpkg/parent/child/c.proto",
		label.New("", "subpkg.parent.child", "child_rs"))

	resolver.Provide("proto", "depends",
		"subpkg/parent/child/c.proto",
		label.New("", "subpkg/parent", "p.proto"))

	r := makeLibraryRule("child_proto", "subpkg/parent/child", []string{"c.proto"})
	from := label.New("", "subpkg/parent/child", "child_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	want := []string{
		"extern_path=.subpkg.parent=::parent_rs",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveTransitiveExternPaths:\n got: %v\nwant: %v", got, want)
	}
}

// TestResolveExternPathOptionsForReferences_SubpackageOfImport verifies the
// reference-emitting variant (used by prost-serde and tonic) DOES add a self
// extern_path override for the current sub-package, so prost's longest-
// prefix-wins matching routes own-package references to crate::... rather
// than the parent extern crate.
func TestResolveExternPathOptionsForReferences_SubpackageOfImport(t *testing.T) {
	resolver := protoc.GlobalResolver()

	resolver.Provide("proto", "prost_extern",
		"refs/parent/p.proto",
		label.New("", "refs.parent", "parent_rs"))

	resolver.Provide("proto", "prost_extern",
		"refs/parent/child/c.proto",
		label.New("", "refs.parent.child", "child_rs"))

	resolver.Provide("proto", "depends",
		"refs/parent/child/c.proto",
		label.New("", "refs/parent", "p.proto"))

	r := makeLibraryRule("child_proto", "refs/parent/child", []string{"c.proto"})
	from := label.New("", "refs/parent/child", "child_proto")

	cfg := &protoc.PluginConfiguration{Options: nil}
	got := prost.ResolveExternPathOptionsForReferences(cfg, r, from)
	want := []string{
		"extern_path=.refs.parent.child=crate",
		"extern_path=.refs.parent=::parent_rs",
	}
	sort.Strings(want)
	sort.Strings(got)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveExternPathOptionsForReferences:\n got: %v\nwant: %v", got, want)
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
	want := []string{"extern_path=.sibling.a.x=::x_rs"}
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
			"extern_path=.stale.pkg=::stale_rs",
		},
	}

	got := prost.ResolveExternPathOptions(cfg, r, from)
	for _, opt := range got {
		if opt == "extern_path=.stale.pkg=::stale_rs" {
			t.Errorf("stale extern_path option was not filtered: %v", got)
		}
	}

	want := []string{"compile_well_known_types=true"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveExternPathOptions:\n got: %v\nwant: %v", got, want)
	}
}

// TestResolveTransitiveExternPaths_GoogleProtobufNotSkipped guards against the
// historical hard-coded skip of google/protobuf/* in the dep walk. When a
// downstream repo registers a rust target for google.protobuf (via
// proto_language rust enable true on that package), the extern_path entry
// must flow through so consumers reference ::google_protobuf instead of the
// prost-build default ::prost_types — which carries no serde impls.
func TestResolveTransitiveExternPaths_GoogleProtobufNotSkipped(t *testing.T) {
	resolver := protoc.GlobalResolver()

	resolver.Provide("proto", "prost_extern",
		"google/protobuf/duration.proto",
		label.New("", "google.protobuf", "google_protobuf"))
	resolver.Provide("proto", "prost_extern",
		"wkttest/consumer/c.proto",
		label.New("", "wkttest.consumer", "wkttest_consumer"))

	resolver.Provide("proto", "depends",
		"wkttest/consumer/c.proto",
		label.New("", "google/protobuf", "duration.proto"))

	r := makeLibraryRule("c_proto", "wkttest/consumer", []string{"c.proto"})
	from := label.New("", "wkttest/consumer", "c_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	want := []string{"extern_path=.google.protobuf=::google_protobuf"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveTransitiveExternPaths:\n got: %v\nwant: %v", got, want)
	}
}

// TestResolveProstOptions_AddsCompileWellKnownTypes_OwnPackage verifies that
// when the library compiles google.protobuf protos itself, the prost plugin
// prepends compile_well_known_types=true so prost-build emits well-known
// type definitions locally (default behaviour is to skip them, which leaves
// the matching prost-serde impls referencing undefined structs).
func TestResolveProstOptions_AddsCompileWellKnownTypes_OwnPackage(t *testing.T) {
	resolver := protoc.GlobalResolver()

	resolver.Provide("proto", "prost_extern",
		"wktown/google/protobuf/any.proto",
		label.New("", "google.protobuf", "google_protobuf"))

	r := makeLibraryRule("google_protobuf_proto", "wktown/google/protobuf", []string{"any.proto"})
	from := label.New("", "wktown/google/protobuf", "google_protobuf_proto")

	plugin := &prost.ProtocGenProstPlugin{}
	got := plugin.ResolvePluginOptions(&protoc.PluginConfiguration{}, r, from)

	if len(got) == 0 || got[0] != "compile_well_known_types=true" {
		t.Errorf("expected compile_well_known_types=true at head of options for own-google.protobuf library, got: %v", got)
	}
}

// TestResolveProstOptions_AddsCompileWellKnownTypes_ExternPath verifies that
// when a library *consumes* google.protobuf via a foreign crate, the prost
// plugin emits compile_well_known_types=true alongside the extern_path.
// Without the flag, prost-build registers its default
// .google.protobuf=::prost_types extern path and ExternPaths::insert errors
// out with "duplicate extern Protobuf path: .google.protobuf".
func TestResolveProstOptions_AddsCompileWellKnownTypes_ExternPath(t *testing.T) {
	resolver := protoc.GlobalResolver()

	resolver.Provide("proto", "prost_extern",
		"wktdep/google/protobuf/duration.proto",
		label.New("", "google.protobuf", "google_protobuf"))
	resolver.Provide("proto", "prost_extern",
		"wktdep/consumer/c.proto",
		label.New("", "wktdep.consumer", "wktdep_consumer"))
	resolver.Provide("proto", "depends",
		"wktdep/consumer/c.proto",
		label.New("", "wktdep/google/protobuf", "duration.proto"))

	r := makeLibraryRule("c_proto", "wktdep/consumer", []string{"c.proto"})
	from := label.New("", "wktdep/consumer", "c_proto")

	plugin := &prost.ProtocGenProstPlugin{}
	got := plugin.ResolvePluginOptions(&protoc.PluginConfiguration{}, r, from)

	wantHead := "compile_well_known_types=true"
	wantExtern := "extern_path=.google.protobuf=::google_protobuf"
	if len(got) == 0 || got[0] != wantHead {
		t.Errorf("expected %q as first option, got: %v", wantHead, got)
	}
	found := false
	for _, opt := range got {
		if opt == wantExtern {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected %q in options, got: %v", wantExtern, got)
	}
}

// TestResolveProstOptions_NoCompileWellKnownTypes_WhenIrrelevant verifies
// that the flag is NOT emitted when the library neither compiles nor
// consumes google.protobuf. Guards against accidental over-emission.
func TestResolveProstOptions_NoCompileWellKnownTypes_WhenIrrelevant(t *testing.T) {
	resolver := protoc.GlobalResolver()

	resolver.Provide("proto", "prost_extern",
		"wktneg/leaf/l.proto",
		label.New("", "wktneg.leaf", "wktneg_leaf"))

	r := makeLibraryRule("l_proto", "wktneg/leaf", []string{"l.proto"})
	from := label.New("", "wktneg/leaf", "l_proto")

	plugin := &prost.ProtocGenProstPlugin{}
	got := plugin.ResolvePluginOptions(&protoc.PluginConfiguration{}, r, from)

	for _, opt := range got {
		if opt == "compile_well_known_types=true" {
			t.Errorf("did not expect compile_well_known_types=true for irrelevant library, got: %v", got)
		}
	}
}

// TestResolveTransitiveExternPaths_MergedLibrariesOwn covers the merged-library
// case: proto_compile/proto_compiled_sources collapses two proto_libraries into
// one rule via the protos= attribute. Before the fix, only the first library's
// srcs went into ownFiles, so files from the second library were treated as
// external in the dep walk and produced a self-referential extern_path entry
// like .google.api=::google_api. The fix is for the rule to carry the full
// []ProtoLibrary set under MergedProtoLibrariesKey, and ResolveTransitive-
// ExternPaths to iterate all of them.
func TestResolveTransitiveExternPaths_MergedLibrariesOwn(t *testing.T) {
	resolver := protoc.GlobalResolver()

	// Two libraries sharing one proto package (mergetest.merged) — emulates
	// e.g. google/api's annotations_proto + http_proto pair, both declaring
	// `package google.api;`.
	resolver.Provide("proto", "prost_extern",
		"mergetest/merged/a.proto",
		label.New("", "mergetest.merged", "mergetest_merged"))
	resolver.Provide("proto", "prost_extern",
		"mergetest/merged/b.proto",
		label.New("", "mergetest.merged", "mergetest_merged"))

	// a.proto imports b.proto — the dep walk would reach b.proto from a.proto.
	resolver.Provide("proto", "depends",
		"mergetest/merged/a.proto",
		label.New("", "mergetest/merged", "b.proto"))

	// Build a single merged rule directly: instead of ProtoLibraryKey only,
	// set MergedProtoLibrariesKey with both backing libraries.
	r := rule.NewRule("proto_library", "merged_proto")
	libA := protoc.NewOtherProtoLibrary(nil,
		makeLibraryRule("a_proto", "mergetest/merged", []string{"a.proto"}),
		protoc.NewFile("mergetest/merged", "a.proto"))
	libB := protoc.NewOtherProtoLibrary(nil,
		makeLibraryRule("b_proto", "mergetest/merged", []string{"b.proto"}),
		protoc.NewFile("mergetest/merged", "b.proto"))
	r.SetPrivateAttr(protoc.MergedProtoLibrariesKey, []protoc.ProtoLibrary{libA, libB})

	from := label.New("", "mergetest/merged", "merged_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	for _, opt := range got {
		if opt == "extern_path=.mergetest.merged=::mergetest_merged" {
			t.Errorf("got self-referential extern_path for merged-in own library, all options: %v", got)
		}
	}
	if len(got) != 0 {
		t.Errorf("expected no extern_path entries (all files are own), got: %v", got)
	}
}
