package rules_rust

import (
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
		"srcs":      true,
		"deps":      true,
		"reexports": true,
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
	KindName             string
	RuleNameSuffix       string
	Outputs              []string
	Config               *protoc.ProtocConfiguration
	RuleConfig           *protoc.LanguageRuleConfig
	Resolver             protoc.DepsResolver
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
}
