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
	// TransitiveExternPathsKey caches computed extern_path options on the
	// library rule's private attrs.
	TransitiveExternPathsKey = "_transitive_extern_paths"
)

// ResolveExternPathOptions filters existing extern_path= options from
// cfg.Options, resolves transitive extern paths, and returns the combined
// options list. This is the common implementation shared by prost, serde,
// and tonic plugins.
func ResolveExternPathOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	externPaths := ResolveTransitiveExternPaths(r, from)

	options := make([]string, 0)
	for _, opt := range cfg.Options {
		if !strings.HasPrefix(opt, "extern_path=") {
			options = append(options, opt)
		}
	}

	options = append(options, externPaths...)
	return options
}

// ResolveTransitiveExternPaths walks the transitive dependency graph of proto
// files and builds extern_path option strings for each dependency package.
func ResolveTransitiveExternPaths(r *rule.Rule, from label.Label) []string {
	lib := r.PrivateAttr(protoc.ProtoLibraryKey)
	if lib == nil {
		return nil
	}
	library := lib.(protoc.ProtoLibrary)
	libRule := library.Rule()

	// Check cache
	if cached, ok := libRule.PrivateAttr(TransitiveExternPathsKey).([]string); ok {
		return cached
	}

	resolver := protoc.GlobalResolver()

	// Build set of own proto files to exclude from extern_paths
	ownFiles := make(map[string]bool)
	for _, src := range library.Srcs() {
		ownFiles[path.Join(from.Pkg, src)] = true
	}

	// Build set of own proto packages. prost's extern_path matches by package
	// prefix, so any imported package that equals or is a parent of one of our
	// own packages would cause prost to rewrite our own types as references
	// into the imported crate. Collect own packages here and use them below
	// to filter such entries out.
	ownPackages := make(map[string]bool)
	for ownFile := range ownFiles {
		for _, ext := range resolver.Resolve("proto", "prost_extern", ownFile) {
			if ext.Label.Pkg != "" {
				ownPackages[ext.Label.Pkg] = true
			}
		}
	}

	// BFS over transitive proto file dependencies
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

		// Walk dependencies
		depends := resolver.Resolve("proto", "depends", protofile)
		for _, dep := range depends {
			depFile := path.Join(dep.Label.Pkg, dep.Label.Name)
			stack.PushBack(depFile)
		}

		// Skip own files
		if ownFiles[protofile] {
			continue
		}

		// Skip well-known types
		if strings.HasPrefix(protofile, "google/protobuf/") {
			continue
		}

		// Look up prost_extern data for this proto file
		results := resolver.Resolve("proto", "prost_extern", protofile)
		if len(results) == 0 {
			continue
		}

		first := results[0]
		protoPackage := first.Label.Pkg // proto package name
		crateName := first.Label.Name   // crate name (e.g., "v1beta1_rs")

		if protoPackage == "" {
			continue
		}

		// Skip extern_path entries that would shadow our own packages. prost's
		// extern_path matches by package prefix, so emitting one for a package
		// that equals or is a parent of one of our own packages would cause
		// prost to rewrite our own type references into the imported crate.
		if isOwnOrParentOfOwn(protoPackage, ownPackages) {
			continue
		}

		// Deduplicate by proto package
		if _, exists := externPathsByPackage[protoPackage]; exists {
			continue
		}

		// extern_path=.{proto_package}=::{crate_name}::{rust_module_path}
		rustModulePath := strings.ReplaceAll(protoPackage, ".", "::")
		externPath := "extern_path=." + protoPackage + "=::" + crateName + "::" + rustModulePath
		externPathsByPackage[protoPackage] = externPath
	}

	result := make([]string, 0, len(externPathsByPackage))
	for _, ep := range externPathsByPackage {
		result = append(result, ep)
	}
	sort.Strings(result)

	// Cache on the library rule
	libRule.SetPrivateAttr(TransitiveExternPathsKey, result)

	return result
}

// isOwnOrParentOfOwn reports whether protoPackage equals one of ownPackages
// or is a proto-package-prefix parent of one (e.g. "a.b" is a parent of
// "a.b.c"). Used to filter dependency extern_path entries that would
// otherwise shadow the current library's own type references through prost's
// prefix-matching extern_path semantics.
func isOwnOrParentOfOwn(protoPackage string, ownPackages map[string]bool) bool {
	for own := range ownPackages {
		if own == protoPackage || strings.HasPrefix(own, protoPackage+".") {
			return true
		}
	}
	return false
}
