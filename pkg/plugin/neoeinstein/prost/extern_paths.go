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
//
// Per-file routing is gated on the **upstream** package's mode (i.e. whether
// it has PerFileTypeProvideKind entries registered). Consumer mode is
// irrelevant: per-package consumers of per-file upstreams also need the
// per-type extern_paths because the upstream's facade label no longer
// exists.
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

		// Suppress the per-package facade extern_path whenever the
		// upstream package is in per-file mode (has any
		// PerFileTypeProvideKind entries). With the facade gone, that
		// label no longer exists; references route through per-type
		// extern_paths emitted by perFileCrossPackageTypeExternPaths
		// below instead. This applies whether the consumer itself is
		// per-file or per-package — a per-package consumer of a per-
		// file upstream still needs per-type routing because there's
		// no facade to depend on.
		if hasPerFileTypeEntries(resolver, protoPackage) {
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

	// Per-type extern_path entries for sibling per-file crates in the same
	// proto package — without these, prost emits a relative path for
	// cross-file same-package references that won't resolve once each
	// file lives in its own crate. The siblings' types were registered
	// against PerFileTypeProvideKind keyed by proto package; we look them
	// up by our own packages and skip entries whose crate is one of ours.
	result = append(result, perFileSiblingTypeExternPaths(resolver, ownPackages, ownCrates)...)

	// Per-type entries for cross-package per-file upstreams. Emitted
	// unconditionally — `perFileCrossPackageTypeExternPaths` walks the
	// transitive `seen` set and filters per-upstream-mode internally.
	result = append(result, perFileCrossPackageTypeExternPaths(resolver, ownPackages, seen)...)

	sort.Strings(result)

	cacheRule.SetPrivateAttr(TransitiveExternPathsKey, result)
	return result
}

// hasPerFileTypeEntries reports whether the resolver has any
// PerFileTypeProvideKind entry for the given proto package. Used as the
// per-file-mode detector for upstream packages.
func hasPerFileTypeEntries(resolver protoc.ImportResolver, protoPackage string) bool {
	return len(resolver.Resolve("proto", PerFileTypeProvideKind, protoPackage)) > 0
}

// perFileCrossPackageTypeExternPaths emits per-type extern_path entries for
// every cross-package per-file upstream that appears in the transitive
// import graph walked above (the `seen` set). For each upstream proto file
// we look up its package via prost_extern, then enumerate the package's
// PerFileTypeProvideKind entries to emit one extern_path per declared type.
//
// Same-package types are not emitted here (they're handled by
// perFileSiblingTypeExternPaths).
func perFileCrossPackageTypeExternPaths(
	resolver protoc.ImportResolver,
	ownPackages map[string]bool,
	seen map[string]bool,
) []string {
	out := make([]string, 0)
	emitted := make(map[string]bool)
	emittedPackages := make(map[string]bool)

	for protofile := range seen {
		externs := resolver.Resolve("proto", "prost_extern", protofile)
		if len(externs) == 0 {
			continue
		}
		protoPackage := externs[0].Label.Pkg
		if protoPackage == "" || ownPackages[protoPackage] || emittedPackages[protoPackage] {
			continue
		}
		entries := resolver.Resolve("proto", PerFileTypeProvideKind, protoPackage)
		if len(entries) == 0 {
			continue
		}
		emittedPackages[protoPackage] = true
		for _, ent := range entries {
			typeName := ent.Label.Pkg
			crateName := ent.Label.Name
			if typeName == "" || crateName == "" {
				continue
			}
			// `typeName` is a dotted proto type path (e.g. `Outer.Inner`)
			// for nested types. Convert through protoTypePathToRustPath so
			// the rust_path becomes `outer::Inner`, matching prost-build's
			// nested-module emission. The previous direct toUpperCamel
			// call left the dot literal in the rust_path, producing the
			// syntactically invalid `::<crate>::Outer.Inner`.
			entry := "extern_path=." + protoPackage + "." + typeName + "=::" + crateName + "::" + protoTypePathToRustPath(typeName)
			if emitted[entry] {
				continue
			}
			emitted[entry] = true
			out = append(out, entry)
		}
	}
	return out
}

// toUpperCamel converts a proto type name to the Rust UpperCamel form prost
// emits. Matches prost-build's `to_upper_camel` (which delegates to heck's
// `ToUpperCamelCase`): consecutive runs of uppercase letters are treated as
// a single word, with a new word starting when an uppercase letter is
// followed by a lowercase letter (e.g. `PIKInfo` → `PikInfo`, `URLLoader` →
// `UrlLoader`).
//
// The transform also splits on `_` and `-` and capitalizes each resulting
// word's first character.
func toUpperCamel(name string) string {
	runes := []rune(name)
	if len(runes) == 0 {
		return name
	}
	var words []string
	start := 0
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '_' || r == '-' {
			if start < i {
				words = append(words, string(runes[start:i]))
			}
			start = i + 1
			continue
		}
		if i == 0 {
			continue
		}
		prev := runes[i-1]
		// lower → upper boundary: e.g. `MyType` splits at `T`.
		if !unicodeUpper(prev) && unicodeUpper(r) {
			words = append(words, string(runes[start:i]))
			start = i
			continue
		}
		// Acronym tail: prev=upper, curr=upper, next=lower splits before curr.
		// e.g. `PIKInfo` at i=3 ('I'): prev='K'(upper), curr='I'(upper),
		// next='n'(lower) → boundary before 'I' so the word is `PIK`.
		if unicodeUpper(prev) && unicodeUpper(r) && i+1 < len(runes) && unicodeLower(runes[i+1]) {
			words = append(words, string(runes[start:i]))
			start = i
		}
	}
	if start < len(runes) {
		words = append(words, string(runes[start:]))
	}

	var b strings.Builder
	for _, w := range words {
		if w == "" {
			continue
		}
		wr := []rune(w)
		b.WriteRune(unicodeToUpper(wr[0]))
		for _, r := range wr[1:] {
			b.WriteRune(unicodeToLower(r))
		}
	}
	return b.String()
}

func unicodeUpper(r rune) bool   { return r >= 'A' && r <= 'Z' }
func unicodeLower(r rune) bool   { return r >= 'a' && r <= 'z' }
func unicodeToUpper(r rune) rune { return rune(strings.ToUpper(string(r))[0]) }
func unicodeToLower(r rune) rune { return rune(strings.ToLower(string(r))[0]) }

// protoTypePathToRustPath converts a dotted proto type path like
// `Outer.Inner.Leaf` to the matching Rust path `outer::inner::Leaf` —
// every segment except the last becomes a snake-cased module name
// (matching prost-build's nested-message module layout), and the last
// segment becomes the upper-camel-cased struct/enum name. A single-
// segment input like `MyMessage` is returned as `MyMessage`.
func protoTypePathToRustPath(typePath string) string {
	parts := strings.Split(typePath, ".")
	if len(parts) == 1 {
		return toUpperCamel(parts[0])
	}
	var b strings.Builder
	for i, p := range parts {
		if i > 0 {
			b.WriteString("::")
		}
		if i == len(parts)-1 {
			b.WriteString(toUpperCamel(p))
		} else {
			// Module name: snake_case of the upper-camel form, matching
			// prost-build's `push_mod(&to_snake(...))` convention for
			// nested messages.
			b.WriteString(toSnakeFromCamel(p))
		}
	}
	return b.String()
}

// toSnakeFromCamel converts an UpperCamel/lowerCamel/UNDERSCORE_CASE/`PIKInfo`
// identifier to snake_case. Mirrors heck's `ToSnakeCase` behavior used by
// prost-build for nested-message module names.
func toSnakeFromCamel(s string) string {
	// Reuse our word-boundary detection (it's the same as toUpperCamel's),
	// then lowercase each word and join with `_`.
	runes := []rune(s)
	if len(runes) == 0 {
		return s
	}
	var words []string
	start := 0
	for i := 0; i < len(runes); i++ {
		r := runes[i]
		if r == '_' || r == '-' {
			if start < i {
				words = append(words, string(runes[start:i]))
			}
			start = i + 1
			continue
		}
		if i == 0 {
			continue
		}
		prev := runes[i-1]
		if !unicodeUpper(prev) && unicodeUpper(r) {
			words = append(words, string(runes[start:i]))
			start = i
			continue
		}
		if unicodeUpper(prev) && unicodeUpper(r) && i+1 < len(runes) && unicodeLower(runes[i+1]) {
			words = append(words, string(runes[start:i]))
			start = i
		}
	}
	if start < len(runes) {
		words = append(words, string(runes[start:]))
	}

	var b strings.Builder
	for i, w := range words {
		if w == "" {
			continue
		}
		if i > 0 {
			b.WriteRune('_')
		}
		for _, r := range w {
			b.WriteRune(unicodeToLower(r))
		}
	}
	return b.String()
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
			//
			// Nested types arrive as dotted paths (e.g. `Outer.Inner`). Each
			// non-final segment becomes a snake-cased Rust module
			// (`outer::Inner`), matching prost-build's nested-message layout.
			// The leaf segment is upper-camel-cased — proto names like
			// `PIKInfo` become `PikInfo` in generated Rust.
			entry := "extern_path=." + pkg + "." + typeName + "=::" + crateName + "::" + protoTypePathToRustPath(typeName)
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
// extern_path strings, keeping only the first mapping for each protobuf path.
// Callers sort the paths first, making collisions deterministic.
func mergeExternPathOptions(cfg *protoc.PluginConfiguration, externPaths []string) []string {
	options := make([]string, 0, len(cfg.Options)+len(externPaths))
	seen := make(map[string]bool, len(externPaths))
	for _, opt := range cfg.Options {
		if !strings.HasPrefix(opt, "extern_path=") {
			options = append(options, opt)
		}
	}
	for _, opt := range externPaths {
		key := opt
		const prefix = "extern_path="
		if strings.HasPrefix(opt, prefix) {
			if eq := strings.IndexByte(opt[len(prefix):], '='); eq >= 0 {
				key = opt[:len(prefix)+eq]
			}
		}
		if seen[key] {
			continue
		}
		seen[key] = true
		options = append(options, opt)
	}
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
