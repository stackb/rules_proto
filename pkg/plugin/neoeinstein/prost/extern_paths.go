package prost

import (
	"container/list"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

const (
	// TransitiveExternPathsKey caches the dependency-only extern_path option
	// strings on the library rule's private attrs.
	TransitiveExternPathsKey = "_transitive_extern_paths"
	// OwnProtoPackagesKey caches the set of proto packages the library
	// itself contributes, used to compute self-extern overrides for
	// reference-emitting plugins (serde, tonic).
	OwnProtoPackagesKey = "_own_proto_packages"
)

// ResolveExternPathOptions filters existing extern_path= options from
// cfg.Options, resolves transitive dependency extern paths, and returns the
// combined options list.
//
// This variant is used by protoc-gen-prost. It does NOT add self-extern
// overrides for the library's own packages because prost interprets such an
// entry as "this package is external — skip generating types for it" and
// emits an empty stub.
//
// It also drops any dependency extern_path whose proto package is a strict
// prefix-parent of one of the library's own packages, for the same reason:
// prost's prefix-matching extern_path semantics treat a sub-package as
// matched and skip generation. Cross-crate references that would otherwise
// have used those filtered extern_paths emerge from prost as relative
// super::... paths; the proto_rust_library macro's generated lib.rs adds
// re-export shims to satisfy them.
func ResolveExternPathOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	parents := ResolveTransitiveExternPaths(r, from)
	owns := ownProtoPackages(r, from)
	if len(owns) > 0 {
		filtered := make([]string, 0, len(parents))
		for _, ep := range parents {
			pkg := externPathPackage(ep)
			if pkg != "" && isParentOfAnyOwn(pkg, owns) {
				continue
			}
			filtered = append(filtered, ep)
		}
		parents = filtered
	}
	return mergeExternPathOptions(cfg, parents)
}

// externPathPackage extracts the proto package from an "extern_path=.{pkg}=..."
// option string, or returns "" if the input doesn't match the expected
// format.
func externPathPackage(opt string) string {
	const prefix = "extern_path=."
	if !strings.HasPrefix(opt, prefix) {
		return ""
	}
	rest := opt[len(prefix):]
	eq := strings.IndexByte(rest, '=')
	if eq < 0 {
		return ""
	}
	return rest[:eq]
}

// isParentOfAnyOwn reports whether pkg equals, or is a strict
// proto-package-prefix parent of, any package in ownPackages.
func isParentOfAnyOwn(pkg string, ownPackages map[string]bool) bool {
	for own := range ownPackages {
		if own == pkg || strings.HasPrefix(own, pkg+".") {
			return true
		}
	}
	return false
}

// ResolveExternPathOptionsForReferences returns ResolveExternPathOptions plus
// self extern_path entries for the library's own proto packages whenever any
// of those packages is a strict sub-package of an imported (parent) package.
//
// Used by protoc-gen-prost-serde and protoc-gen-tonic. Both emit Rust code at
// crate-root using absolute crate-qualified paths; without a self-extern
// override prost's longest-prefix-wins matching would route a reference like
// ".pkg.sub.MyType" through the parent's external crate instead of resolving
// it to crate::pkg::sub::MyType.
func ResolveExternPathOptionsForReferences(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	parents := ResolveTransitiveExternPaths(r, from)
	owns := ownProtoPackages(r, from)
	selves := selfExternPathsForOverride(owns, parents)

	all := make([]string, 0, len(parents)+len(selves))
	all = append(all, parents...)
	all = append(all, selves...)
	sort.Strings(all)
	return mergeExternPathOptions(cfg, all)
}

// ResolveTransitiveExternPaths walks the transitive dependency graph of
// proto files and builds an extern_path option string for each dependency
// package. Self-extern overrides are NOT included — see
// ResolveExternPathOptionsForReferences for the variant that adds them.
func ResolveTransitiveExternPaths(r *rule.Rule, from label.Label) []string {
	lib := r.PrivateAttr(protoc.ProtoLibraryKey)
	if lib == nil {
		return nil
	}
	library := lib.(protoc.ProtoLibrary)
	libRule := library.Rule()

	if cached, ok := libRule.PrivateAttr(TransitiveExternPathsKey).([]string); ok {
		return cached
	}

	resolver := protoc.GlobalResolver()

	ownFiles := make(map[string]bool)
	for _, src := range library.Srcs() {
		ownFiles[path.Join(from.Pkg, src)] = true
	}

	seen := make(map[string]bool)
	stack := list.New()
	for _, src := range library.Srcs() {
		stack.PushBack(path.Join(from.Pkg, src))
	}

	externPathsByPackage := make(map[string]string)

	for stack.Len() > 0 {
		current := stack.Front()
		stack.Remove(current)

		protofile := current.Value.(string)
		if seen[protofile] {
			continue
		}
		seen[protofile] = true

		depends := resolver.Resolve("proto", "depends", protofile)
		for _, dep := range depends {
			depFile := path.Join(dep.Label.Pkg, dep.Label.Name)
			stack.PushBack(depFile)
		}

		if ownFiles[protofile] {
			continue
		}

		// Skip well-known types — prost ships these built-in.
		if strings.HasPrefix(protofile, "google/protobuf/") {
			continue
		}

		results := resolver.Resolve("proto", "prost_extern", protofile)
		if len(results) == 0 {
			continue
		}

		first := results[0]
		protoPackage := first.Label.Pkg
		crateName := first.Label.Name
		if protoPackage == "" {
			continue
		}
		if _, exists := externPathsByPackage[protoPackage]; exists {
			continue
		}

		// extern_path=.{proto_package}=::{crate_name}
		// The crate exposes all generated types at its root (see the
		// proto_rust_library Starlark macro), so no nested module path is
		// appended after the crate name.
		externPathsByPackage[protoPackage] = "extern_path=." + protoPackage + "=::" + crateName
	}

	result := make([]string, 0, len(externPathsByPackage))
	for _, ep := range externPathsByPackage {
		result = append(result, ep)
	}
	sort.Strings(result)

	libRule.SetPrivateAttr(TransitiveExternPathsKey, result)
	return result
}

// mergeExternPathOptions strips any pre-existing extern_path= entries from
// cfg.Options and returns the remainder concatenated with the supplied
// extern_path strings.
func mergeExternPathOptions(cfg *protoc.PluginConfiguration, externPaths []string) []string {
	options := make([]string, 0, len(cfg.Options)+len(externPaths))
	for _, opt := range cfg.Options {
		if !strings.HasPrefix(opt, "extern_path=") {
			options = append(options, opt)
		}
	}
	options = append(options, externPaths...)
	return options
}

// ownProtoPackages returns the set of proto packages the library itself
// contributes, computed from prost_extern resolver entries for each own
// proto file. Cached on the library rule.
func ownProtoPackages(r *rule.Rule, from label.Label) map[string]bool {
	lib := r.PrivateAttr(protoc.ProtoLibraryKey)
	if lib == nil {
		return nil
	}
	library := lib.(protoc.ProtoLibrary)
	libRule := library.Rule()

	if cached, ok := libRule.PrivateAttr(OwnProtoPackagesKey).(map[string]bool); ok {
		return cached
	}

	resolver := protoc.GlobalResolver()
	out := make(map[string]bool)
	for _, src := range library.Srcs() {
		ownFile := path.Join(from.Pkg, src)
		for _, ext := range resolver.Resolve("proto", "prost_extern", ownFile) {
			if ext.Label.Pkg != "" {
				out[ext.Label.Pkg] = true
			}
		}
	}

	libRule.SetPrivateAttr(OwnProtoPackagesKey, out)
	return out
}

// selfExternPathsForOverride returns "extern_path=.{ownPkg}=crate::..."
// entries for every own proto package whose path is a strict sub-package of
// any package present in parents. parents is the slice of dependency
// extern_path option strings (as returned by ResolveTransitiveExternPaths).
func selfExternPathsForOverride(ownPackages map[string]bool, parents []string) []string {
	if len(ownPackages) == 0 || len(parents) == 0 {
		return nil
	}
	parentPkgs := parentExternPackages(parents)
	out := make([]string, 0)
	for ownPkg := range ownPackages {
		if !hasParentInImports(ownPkg, parentPkgs) {
			continue
		}
		// All own types live at the crate root (flat convention), so the
		// self-extern override maps the proto sub-package to bare `crate`.
		out = append(out, "extern_path=."+ownPkg+"=crate")
	}
	return out
}

// parentExternPackages parses a slice of "extern_path=.{pkg}=..." strings
// and returns the set of proto packages they cover.
func parentExternPackages(opts []string) map[string]bool {
	out := make(map[string]bool, len(opts))
	const prefix = "extern_path=."
	for _, opt := range opts {
		if !strings.HasPrefix(opt, prefix) {
			continue
		}
		rest := opt[len(prefix):]
		eq := strings.IndexByte(rest, '=')
		if eq < 0 {
			continue
		}
		out[rest[:eq]] = true
	}
	return out
}

// hasParentInImports reports whether any of importedPackages is a proto-
// package-prefix parent of ownPkg (e.g. "a.b" is a parent of "a.b.c").
func hasParentInImports(ownPkg string, importedPackages map[string]bool) bool {
	for imp := range importedPackages {
		if strings.HasPrefix(ownPkg, imp+".") {
			return true
		}
	}
	return false
}
