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

// TestResolveTransitiveExternPaths_PerFileSiblingTypes covers the headline
// per-file behaviour: when our per-file proto_library imports another file
// in the *same* proto package (a sibling per-file crate), ResolveTransitive-
// ExternPaths emits per-type extern_path entries routing each sibling type
// through its own per-file crate, AND suppresses the per-package extern_path
// that would otherwise hijack our own crate's references via prost's
// longest-prefix matching.
func TestResolveTransitiveExternPaths_PerFileSiblingTypes(t *testing.T) {
	resolver := protoc.GlobalResolver()

	// Two per-file crates in the same proto package. The per-file crate
	// naming convention is `<pkg>__<file_stem>` — see protoc-gen-prost's
	// registerExternPaths().
	const pkg = "perfiletest.pkg"
	const ourCrate = "perfiletest_pkg__order"
	const siblingCrate = "perfiletest_pkg__trade"

	resolver.Provide("proto", "prost_extern",
		"perfiletest/pkg/order.proto",
		label.New("", pkg, ourCrate))
	resolver.Provide("proto", "prost_extern",
		"perfiletest/pkg/trade.proto",
		label.New("", pkg, siblingCrate))

	// Per-type registry, keyed by proto package: each entry is
	// (typeName, crateName).
	resolver.Provide("proto", prost.PerFileTypeProvideKind, pkg,
		label.New("", "Order", ourCrate))
	resolver.Provide("proto", prost.PerFileTypeProvideKind, pkg,
		label.New("", "Trade", siblingCrate))

	// order.proto imports trade.proto.
	resolver.Provide("proto", "depends",
		"perfiletest/pkg/order.proto",
		label.New("", "perfiletest/pkg", "trade.proto"))

	r := makeLibraryRule("order_proto", "perfiletest/pkg", []string{"order.proto"})
	from := label.New("", "perfiletest/pkg", "order_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	sort.Strings(got)

	want := []string{
		// Per-type entry for the sibling's type — routes references in our
		// crate to the sibling's per-file crate. The trailing `::Trade`
		// is needed so prost lands on the type (not the crate itself) when
		// substituting references.
		"extern_path=.perfiletest.pkg.Trade=::perfiletest_pkg__trade::Trade",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveTransitiveExternPaths:\n got: %v\nwant: %v", got, want)
	}

	// Must NOT contain a per-package extern_path for our own package —
	// that would hijack every type reference in our crate.
	for _, opt := range got {
		if opt == "extern_path=.perfiletest.pkg=::perfiletest_pkg__trade" {
			t.Errorf("per-package extern_path for own proto package leaked: %v", got)
		}
	}

	// Must NOT contain a per-type entry for our own type (would prevent
	// prost from generating its definition locally).
	for _, opt := range got {
		if opt == "extern_path=.perfiletest.pkg.Order=::perfiletest_pkg__order" {
			t.Errorf("per-type extern_path for own type leaked: %v", got)
		}
	}
}

// TestResolveTransitiveExternPaths_PerFileCrossPackageStillEmitsPerPackage
// ensures the per-file machinery is additive: cross-*package* references
// still get a per-package extern_path (because the dep isn't in our proto
// package), even when we're a per-file crate ourselves.
func TestResolveTransitiveExternPaths_PerFileCrossPackageStillEmitsPerPackage(t *testing.T) {
	resolver := protoc.GlobalResolver()

	// Our per-file crate.
	resolver.Provide("proto", "prost_extern",
		"perfilecross/own/o.proto",
		label.New("", "perfilecross.own", "perfilecross_own__o"))

	// A dep in a *different* proto package — its crate is the conventional
	// per-package (façade) name.
	resolver.Provide("proto", "prost_extern",
		"perfilecross/dep/d.proto",
		label.New("", "perfilecross.dep", "perfilecross_dep"))

	// o.proto imports d.proto (cross-package).
	resolver.Provide("proto", "depends",
		"perfilecross/own/o.proto",
		label.New("", "perfilecross/dep", "d.proto"))

	r := makeLibraryRule("o_proto", "perfilecross/own", []string{"o.proto"})
	from := label.New("", "perfilecross/own", "o_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	want := []string{"extern_path=.perfilecross.dep=::perfilecross_dep"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveTransitiveExternPaths:\n got: %v\nwant: %v", got, want)
	}
}

// TestResolveTransitiveExternPaths_PerFileNoSiblingsRegistered verifies that
// when no PerFileTypeProvideKind entries are registered for our proto
// packages (i.e. the per-file plugin path wasn't taken — e.g. package-mode
// surroundings), ResolveTransitiveExternPaths emits zero per-type entries.
// Guards against regression for the default per-package convention.
func TestResolveTransitiveExternPaths_PerFileNoSiblingsRegistered(t *testing.T) {
	resolver := protoc.GlobalResolver()

	resolver.Provide("proto", "prost_extern",
		"perfilenone/own/own.proto",
		label.New("", "perfilenone.own", "perfilenone_own"))

	r := makeLibraryRule("own_proto", "perfilenone/own", []string{"own.proto"})
	from := label.New("", "perfilenone/own", "own_proto")

	got := prost.ResolveTransitiveExternPaths(r, from)
	if len(got) != 0 {
		t.Errorf("expected no extern paths, got %v", got)
	}
}

// TestResolveTransitiveExternPaths_PerFileMultiFileOwn covers the merged
// per-file case: our gazelle rule actually owns *multiple* per-file
// libraries (each in its own proto_library) merged into a single
// proto_rust_library via MergedProtoLibrariesKey. Per-type entries for
// every type any of our own files defines must be skipped — only sibling
// types from siblings that aren't part of our merge set should emerge.
func TestResolveTransitiveExternPaths_PerFileMultiFileOwn(t *testing.T) {
	resolver := protoc.GlobalResolver()

	const pkg = "perfilemulti.pkg"
	const aCrate = "perfilemulti_pkg__a"
	const bCrate = "perfilemulti_pkg__b"
	const cCrate = "perfilemulti_pkg__c"

	// Three sibling per-file crates in the same proto package — we own
	// `a` and `b`, `c` is an external sibling we'd reference via its own
	// crate. (Contrived, but it covers the multi-own / single-sibling
	// composition.)
	resolver.Provide("proto", "prost_extern",
		"perfilemulti/pkg/a.proto", label.New("", pkg, aCrate))
	resolver.Provide("proto", "prost_extern",
		"perfilemulti/pkg/b.proto", label.New("", pkg, bCrate))
	resolver.Provide("proto", "prost_extern",
		"perfilemulti/pkg/c.proto", label.New("", pkg, cCrate))

	resolver.Provide("proto", prost.PerFileTypeProvideKind, pkg,
		label.New("", "A", aCrate))
	resolver.Provide("proto", prost.PerFileTypeProvideKind, pkg,
		label.New("", "B", bCrate))
	resolver.Provide("proto", prost.PerFileTypeProvideKind, pkg,
		label.New("", "C", cCrate))

	// a.proto imports c.proto so c is on the dep stack; b.proto imports
	// a.proto so the merge graph touches every file.
	resolver.Provide("proto", "depends",
		"perfilemulti/pkg/a.proto",
		label.New("", "perfilemulti/pkg", "c.proto"))
	resolver.Provide("proto", "depends",
		"perfilemulti/pkg/b.proto",
		label.New("", "perfilemulti/pkg", "a.proto"))

	// Build a merged rule owning a.proto + b.proto. c.proto is not ours.
	r := rule.NewRule("proto_library", "ab_merged")
	libA := protoc.NewOtherProtoLibrary(nil,
		makeLibraryRule("a_proto", "perfilemulti/pkg", []string{"a.proto"}),
		protoc.NewFile("perfilemulti/pkg", "a.proto"))
	libB := protoc.NewOtherProtoLibrary(nil,
		makeLibraryRule("b_proto", "perfilemulti/pkg", []string{"b.proto"}),
		protoc.NewFile("perfilemulti/pkg", "b.proto"))
	r.SetPrivateAttr(protoc.MergedProtoLibrariesKey, []protoc.ProtoLibrary{libA, libB})

	from := label.New("", "perfilemulti/pkg", "ab_merged")

	got := prost.ResolveTransitiveExternPaths(r, from)
	sort.Strings(got)

	// Only C should appear (A and B are ours). The rust_path includes the
	// trailing `::C` so prost emits `::<crate>::C` rather than `::<crate>`
	// when substituting type references.
	want := []string{
		"extern_path=.perfilemulti.pkg.C=::perfilemulti_pkg__c::C",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ResolveTransitiveExternPaths:\n got: %v\nwant: %v", got, want)
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
