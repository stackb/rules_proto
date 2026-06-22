package rules_rust

import (
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

var rustLibraryKindInfo = rule.KindInfo{
	MergeableAttrs: map[string]bool{
		"srcs":             true,
		"deps":             true,
		"reexports":        true,
		"per_file":         true,
		"per_file_imports": true,
	},
	NonEmptyAttrs: map[string]bool{
		"srcs": true,
	},
	ResolveAttrs: map[string]bool{
		"deps": true,
	},
}

// RustLibrary implements RuleProvider for 'proto_rust_library'-derived rules.
type RustLibrary struct {
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       protoc.DepsResolver
	// PerFile, when true, tells the proto_rust_library macro to expand
	// `srcs` into one rust_library per .proto file plus a per-package
	// façade. Set from the surrounding gazelle:proto file mode — see
	// ProtocConfiguration.IsProtoFileMode.
	PerFile              bool
	id                   label.Label
	protoLibrariesByRule map[label.Label][]protoc.ProtoLibrary
}

// Kind implements part of the RuleProvider interface.
func (s *RustLibrary) Kind() string {
	return s.KindName
}

// Name implements part of the RuleProvider interface.
func (s *RustLibrary) Name() string {
	if pkg := s.Pkg(); pkg != "" {
		return protoc.RustCrateName(pkg)
	}
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Pkg returns the proto package name of the first file in the library, or "".
// Internal helper: drives the rule's crate name (via RustCrateName) and the
// reexport computation. Not propagated to the generated rule as an attribute —
// the consuming Starlark rule no longer accepts a `pkg` attr.
func (s *RustLibrary) Pkg() string {
	files := s.Config.Library.Files()
	if len(files) == 0 {
		return ""
	}
	return files[0].Package().Name
}

// Srcs computes the srcs list for the rule.
func (s *RustLibrary) Srcs() []string {
	srcs := make([]string, 0, len(s.Outputs))
	for _, output := range s.Outputs {
		if strings.HasSuffix(output, ".rs") {
			srcs = append(srcs, protoc.StripRel(s.Config.Rel, output))
		}
	}
	sort.Strings(srcs)
	return srcs
}

// Deps computes the deps list for the rule.
func (s *RustLibrary) Deps() []string {
	deps := s.RuleConfig.GetDeps()
	deps = append(deps, s.fixedDeps()...)
	return protoc.DeduplicateAndSort(deps)
}

// hasServices returns true if any proto file in the library defines services.
func (s *RustLibrary) hasServices() bool {
	for _, f := range s.Config.Library.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// hasWellKnownTypes returns true if any proto file imports a well-known type
// (google/protobuf/*), which requires the prost-types crate at runtime.
func (s *RustLibrary) hasWellKnownTypes() bool {
	for _, f := range s.Config.Library.Files() {
		for _, imp := range f.Imports() {
			if strings.HasPrefix(imp.Filename, "google/protobuf/") {
				return true
			}
		}
	}
	return false
}

// fixedDeps returns the crate dependencies that are always needed.
func (s *RustLibrary) fixedDeps() []string {
	deps := []string{
		"@crates//:prost",
		"@crates//:serde",
		"@crates//:pbjson",
	}
	if s.hasServices() {
		deps = append(deps, "@crates//:tonic")
	}
	if s.hasWellKnownTypes() {
		deps = append(deps, "@crates//:prost-types")
	}
	return deps
}

// Visibility provides visibility labels.
func (s *RustLibrary) Visibility() []string {
	return s.RuleConfig.GetVisibility()
}

// Rule implements part of the RuleProvider interface.
func (s *RustLibrary) Rule(otherGen ...*rule.Rule) *rule.Rule {
	srcs := s.Srcs()
	deps := s.Deps()
	visibility := s.Visibility()
	imports := s.Config.Library.Imports()

	// Check if an existing rule with the same kind and name has already been
	// generated.  If so, merge into it rather than creating a new rule.
	for _, other := range otherGen {
		if other.Kind() == s.Kind() && other.Name() == s.Name() {
			otherLabel := label.New("", s.Config.Rel, other.Name())
			otherSrcs := other.AttrStrings("srcs")
			otherDeps := other.AttrStrings("deps")
			otherVis := other.AttrStrings("visibility")
			otherImports, _ := other.PrivateAttr(config.GazelleImportsKey).([]string)

			other.SetAttr("srcs", protoc.DeduplicateAndSort(append(otherSrcs, srcs...)))
			other.SetAttr("deps", protoc.DeduplicateAndSort(append(otherDeps, deps...)))
			other.SetAttr("visibility", protoc.DeduplicateAndSort(append(otherVis, visibility...)))
			other.SetPrivateAttr(config.GazelleImportsKey, protoc.DeduplicateAndSort(append(otherImports, imports...)))

			s.protoLibrariesByRule[otherLabel] = append(s.protoLibrariesByRule[otherLabel], s.Config.Library)

			return other
		}
	}

	newRule := rule.NewRule(s.Kind(), s.Name())
	newRule.SetAttr("srcs", srcs)
	newRule.SetPrivateAttr(config.GazelleImportsKey, imports)
	s.protoLibrariesByRule[s.id] = []protoc.ProtoLibrary{s.Config.Library}

	if len(deps) > 0 {
		newRule.SetAttr("deps", deps)
	}
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}
	if s.PerFile {
		// Tell the proto_rust_library macro to expand srcs into per-file
		// crates + a per-package façade. See macro docstring in
		// bazel_tools/rust/proto_rust_library.bzl.
		newRule.SetAttr("per_file", true)
	}

	return newRule
}

// Reexports returns "crate_name=proto.package" entries identifying every
// imported package whose proto path is a strict prefix-parent of any of our
// own proto packages. The proto_rust_library Starlark macro uses these to
// generate "pub use ::crate_name::path::*;" re-exports inside the local
// lib.rs at the parent module, which lets prost's relative super::... paths
// for cross-crate references resolve. See the matching prost-side filter in
// extern_paths.ResolveExternPathOptions for context: that filter drops the
// dependency's extern_path entry (which would otherwise make prost skip
// generating the local sub-package), and these re-exports replace what the
// extern_path would have provided for cross-crate type resolution.
func (s *RustLibrary) Reexports() []string {
	ownPkg := s.Pkg()
	if ownPkg == "" {
		return nil
	}

	resolver := protoc.GlobalResolver()
	out := make([]string, 0)
	seen := make(map[string]bool)

	for _, f := range s.Config.Library.Files() {
		for _, imp := range f.Imports() {
			results := resolver.Resolve("proto", "prost_extern", imp.Filename)
			if len(results) == 0 {
				continue
			}
			impPkg := results[0].Label.Pkg
			impCrate := results[0].Label.Name
			if impPkg == "" || impCrate == "" {
				continue
			}
			if !strings.HasPrefix(ownPkg, impPkg+".") {
				// Not a strict prefix-parent of our own package — handled
				// via the regular extern_path mechanism.
				continue
			}
			entry := impCrate + "=" + impPkg
			if seen[entry] {
				continue
			}
			seen[entry] = true
			out = append(out, entry)
		}
	}

	sort.Strings(out)
	return out
}

// Imports implements part of the RuleProvider interface.
func (s *RustLibrary) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	libs, ok := s.protoLibrariesByRule[s.id]
	if !ok {
		return nil
	}
	return protoc.ProtoLibraryImportSpecsForKind(r.Kind(), libs...)
}

// Resolve implements part of the RuleProvider interface.
func (s *RustLibrary) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	s.Resolver(c, ix, r, imports, from)
	if reexports := s.Reexports(); len(reexports) > 0 {
		r.SetAttr("reexports", reexports)
	}
	if s.PerFile {
		if perFileImports := s.PerFileImports(); len(perFileImports) > 0 {
			r.SetAttr("per_file_imports", protoc.MakeStringListDict(perFileImports))
		}
	}
}

// PerFileImports returns a stem→sibling-stems map describing which per-file
// crates each per-file crate in this library must depend on at the bazel
// dep-graph level.
//
// Per-file mode emits one rust_library per .proto file (`<facade>__<stem>`).
// When a file's prost/serde/tonic codegen references a type defined in a
// sibling file in the same proto package, the resulting Rust code uses an
// absolute path through the sibling's per-file crate (e.g.
// `::omnistac_uss_proto__uss_stream::UssStream`) because the extern_path
// registry routes same-package cross-file references that way. Without an
// explicit `deps =` edge for the sibling crate, the consuming per-file
// crate fails to compile with `cannot find crate`.
//
// The mapping is computed once per rule (over all merged proto_libraries
// tracked under `protoLibrariesByRule`) so it survives the rule-merge that
// collapses N per-file proto_libraries into one proto_rust_library.
//
// Cross-package imports are NOT recorded here — they're already handled by
// the cross-package façade entry the gazelle resolver adds to `deps`.
func (s *RustLibrary) PerFileImports() map[string][]string {
	if !s.PerFile {
		return nil
	}

	libs, ok := s.protoLibrariesByRule[s.id]
	if !ok || len(libs) == 0 {
		libs = []protoc.ProtoLibrary{s.Config.Library}
	}

	// Build the set of stems that actually correspond to per-file crates in
	// this package — i.e. files that emit codegen (have messages, enums, or
	// services). Files like a shared `package.proto` are not crates and must
	// not be listed as deps.
	crateStems := make(map[string]bool)
	for _, lib := range libs {
		for _, f := range lib.Files() {
			if !f.HasMessages() && !f.HasEnums() && !f.HasServices() {
				continue
			}
			crateStems[strings.TrimSuffix(f.Basename, ".proto")] = true
		}
	}

	result := make(map[string][]string)
	seen := make(map[string]map[string]bool)

	for _, lib := range libs {
		for _, f := range lib.Files() {
			if !f.HasMessages() && !f.HasEnums() && !f.HasServices() {
				continue
			}
			ownStem := strings.TrimSuffix(f.Basename, ".proto")
			if seen[ownStem] == nil {
				seen[ownStem] = make(map[string]bool)
			}
			for _, imp := range f.Imports() {
				impDir := path.Dir(imp.Filename)
				if impDir == "." {
					impDir = ""
				}
				if impDir != f.Dir {
					// Cross-bazel-package import — handled by the façade
					// dep, not by an inter-per-file edge.
					continue
				}
				impStem := strings.TrimSuffix(path.Base(imp.Filename), ".proto")
				if impStem == ownStem {
					continue
				}
				if !crateStems[impStem] {
					// Sibling has no codegen (e.g. shared `package.proto`)
					// so there's no per-file crate to depend on.
					continue
				}
				if seen[ownStem][impStem] {
					continue
				}
				seen[ownStem][impStem] = true
				result[ownStem] = append(result[ownStem], impStem)
			}
		}
	}

	for stem, sibs := range result {
		sort.Strings(sibs)
		result[stem] = sibs
	}

	return result
}
