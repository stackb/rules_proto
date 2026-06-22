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
	libraries := mergedLibraries(r)
	if len(libraries) == 0 {
		return nil
	}
	// Cache off the first library's underlying rule — merge can occur but the
	// post-merge PrivateAttrs travel with the first library's rule for back-
	// compat with the existing cache placement.
	cacheRule := libraries[0].Rule()
	if cached, ok := cacheRule.PrivateAttr(TransitiveExternPathsKey).([]string); ok {
		return cached
	}

	resolver := protoc.GlobalResolver()

	ownFiles := make(map[string]bool)
	for _, library := range libraries {
		for _, src := range library.Srcs() {
			ownFiles[path.Join(from.Pkg, src)] = true
		}
	}
	// Also treat any proto file whose registered prost_extern crate matches
	// one of our own as own. This is a belt-and-suspenders guard for cases
	// where ownFiles can't identify a merged-in file by path alone (the
	// caller's `from.Pkg` is the rule's bazel package, but a merged proto_-
	// library may sit in a different bazel package).
	ownCrates := ownCrateNames(libraries, resolver, from)
	// Per-file mode: sibling .proto files within the same proto package
	// resolve to *different* per-file crates, so the ownCrates check above
	// doesn't catch them. Tracking our own proto packages lets us suppress
	// the per-package extern_path emission for siblings — they're routed
	// per-type by perFileSiblingTypeExternPaths below. Without this, a
	// `.<our_pkg>=::<sibling_crate>` entry would (via prost's longest-
	// prefix matching) hijack every type reference in our own crate.
	ownPackages := ownProtoPackages(r, from)

	seen := make(map[string]bool)
	stack := list.New()
	for _, library := range libraries {
		for _, src := range library.Srcs() {
			stack.PushBack(path.Join(from.Pkg, src))
		}
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
		// Skip files whose crate is one of ours — they're part of the
		// merged library set, just couldn't be matched by file path.
		if ownCrates[crateName] {
			continue
		}
		// Skip same-proto-package siblings in per-file mode. Per-type
		// extern_paths handle them downstream; a per-package entry here
		// would hijack our own crate's references via prost's longest-
		// prefix matching.
		if ownPackages[protoPackage] {
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

	// In per-file mode, also emit per-type extern_path entries for sibling
	// per-file crates in the same proto package — without these, prost
	// emits a relative path for cross-file same-package references that
	// won't resolve once each file lives in its own crate. The siblings'
	// types were registered against PerFileTypeProvideKind keyed by proto
	// package; we look them up by our own packages and skip entries
	// whose crate is one of ours.
	result = append(result, perFileSiblingTypeExternPaths(resolver, ownPackages, ownCrates)...)

	sort.Strings(result)

	cacheRule.SetPrivateAttr(TransitiveExternPathsKey, result)
	return result
}

// perFileSiblingTypeExternPaths emits extern_path entries that route same-
// package cross-file type references through the correct sibling per-file
// crate. Returns nil if no per-file type entries are registered (the
// common case — only per-file-mode libraries register them).
//
// `ownPackages` is the set of proto packages we contribute (derived via
// `ownProtoPackages` from the resolver, not from File.Package() — the
// parsed proto-package field isn't reliably populated outside the real
// gazelle Generate pass, e.g. in unit tests).
func perFileSiblingTypeExternPaths(
	resolver protoc.ImportResolver,
	ownPackages map[string]bool,
	ownCrates map[string]bool,
) []string {
	var out []string
	seen := make(map[string]bool)
	for pkg := range ownPackages {
		for _, ent := range resolver.Resolve("proto", PerFileTypeProvideKind, pkg) {
			typeName := ent.Label.Pkg
			crateName := ent.Label.Name
			if typeName == "" || crateName == "" {
				continue
			}
			// Skip types declared in our own crate(s) — prost emits them
			// locally; an extern_path here would tell prost to skip
			// generating the definition entirely.
			if ownCrates[crateName] {
				continue
			}
			// Per-type extern_path keys are EXACT-matched by prost-build, which
			// returns the rust_path verbatim — without appending the trailing
			// type segment. Spelling the type name into the rust_path so the
			// generated code lands on `::<crate>::<TypeName>` rather than
			// `::<crate>` (the latter would refer to the crate itself, not
			// the type, and fails compilation as E0425/E0433).
			entry := "extern_path=." + pkg + "." + typeName + "=::" + crateName + "::" + typeName
			if seen[entry] {
				continue
			}
			seen[entry] = true
			out = append(out, entry)
		}
	}
	return out
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
// proto file across all merged proto_libraries. Cached on the rule.
func ownProtoPackages(r *rule.Rule, from label.Label) map[string]bool {
	libraries := mergedLibraries(r)
	if len(libraries) == 0 {
		return nil
	}
	cacheRule := libraries[0].Rule()
	if cached, ok := cacheRule.PrivateAttr(OwnProtoPackagesKey).(map[string]bool); ok {
		return cached
	}

	resolver := protoc.GlobalResolver()
	out := make(map[string]bool)
	for _, library := range libraries {
		for _, src := range library.Srcs() {
			ownFile := path.Join(from.Pkg, src)
			for _, ext := range resolver.Resolve("proto", "prost_extern", ownFile) {
				if ext.Label.Pkg != "" {
					out[ext.Label.Pkg] = true
				}
			}
		}
	}

	cacheRule.SetPrivateAttr(OwnProtoPackagesKey, out)
	return out
}

// mergedLibraries returns the full set of proto_libraries backing a proto_-
// compile / proto_compiled_sources rule. Prefers MergedProtoLibrariesKey
// (set by proto_compile.Rule for every merge — see protoc.MergedProtoLibrariesKey)
// and falls back to wrapping the single ProtoLibraryKey for callers that
// haven't migrated (e.g. proto_rust_library, hand-constructed test rules).
// Returns nil when the rule carries neither.
func mergedLibraries(r *rule.Rule) []protoc.ProtoLibrary {
	if libs, ok := r.PrivateAttr(protoc.MergedProtoLibrariesKey).([]protoc.ProtoLibrary); ok && len(libs) > 0 {
		return libs
	}
	if lib, ok := r.PrivateAttr(protoc.ProtoLibraryKey).(protoc.ProtoLibrary); ok && lib != nil {
		return []protoc.ProtoLibrary{lib}
	}
	return nil
}

// ownCrateNames returns the set of rust crate names registered (via
// prost_extern) for files belonging to any of the library's merged proto_-
// libraries. Used to recognise own-merged files in the dep walk even when
// their on-disk path doesn't share a prefix with from.Pkg (e.g. a
// proto_compiled_sources at //a:foo that merges in //b:bar_proto would
// otherwise see b/bar.proto as external).
func ownCrateNames(libraries []protoc.ProtoLibrary, resolver protoc.ImportResolver, from label.Label) map[string]bool {
	out := make(map[string]bool)
	for _, library := range libraries {
		for _, src := range library.Srcs() {
			ownFile := path.Join(from.Pkg, src)
			for _, ext := range resolver.Resolve("proto", "prost_extern", ownFile) {
				if ext.Label.Name != "" {
					out[ext.Label.Name] = true
				}
			}
		}
	}
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
