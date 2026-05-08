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
		"srcs": true,
		"deps": true,
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
		return strings.ReplaceAll(pkg, ".", "_")
	}
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Pkg returns the proto package name from the first file in the library.
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

	if pkg := s.Pkg(); pkg != "" {
		newRule.SetAttr("pkg", pkg)
	}
	if len(deps) > 0 {
		newRule.SetAttr("deps", deps)
	}
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
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
}
