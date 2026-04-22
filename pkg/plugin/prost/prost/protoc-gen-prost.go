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
	ProtocGenProstPluginName = "neoeinstein:prost:protoc-gen-prost"

	// TransitiveExternPathsKey caches computed extern_path options on the
	// library rule's private attrs.
	TransitiveExternPathsKey = "_transitive_extern_paths"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenProstPlugin{})
}

// ProtocGenProstPlugin implements Plugin for protoc-gen-prost.
type ProtocGenProstPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenProstPlugin) Name() string {
	return ProtocGenProstPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenProstPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}

	outputs := p.outputs(ctx.ProtoLibrary)
	if len(outputs) == 0 {
		return nil
	}

	p.registerExternPaths(ctx.ProtoLibrary)

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/prost/prost", "protoc-gen-prost"),
		Outputs: outputs,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// ResolvePluginOptions implements the PluginOptionsResolver interface.
// It computes extern_path options based on transitive proto file dependencies.
func (p *ProtocGenProstPlugin) ResolvePluginOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	externPaths := p.resolveTransitiveExternPaths(r, from)

	options := make([]string, 0)
	for _, opt := range cfg.Options {
		if !strings.HasPrefix(opt, "extern_path=") {
			options = append(options, opt)
		}
	}

	options = append(options, externPaths...)
	return options
}

// shouldApply returns true if the library has files with messages or enums.
func (p *ProtocGenProstPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// outputs computes the output files for the plugin. Prost generates one .rs
// file per proto package, named {proto_package}.rs. The path includes the
// file's directory so that mergeSources can handle the rel stripping.
func (p *ProtocGenProstPlugin) outputs(lib protoc.ProtoLibrary) []string {
	seen := make(map[string]bool)
	outputs := make([]string, 0)

	for _, f := range lib.Files() {
		if !(f.HasMessages() || f.HasEnums()) {
			continue
		}
		pkg := f.Package()
		if pkg.Name == "" {
			continue
		}
		if seen[pkg.Name] {
			continue
		}
		seen[pkg.Name] = true

		filename := pkg.Name + ".rs"
		if f.Dir != "" {
			filename = path.Join(f.Dir, filename)
		}
		outputs = append(outputs, filename)
	}

	sort.Strings(outputs)
	return outputs
}

// registerExternPaths records prost extern_path data in the GlobalResolver for
// each proto file in the library. This data is later consumed by
// resolveTransitiveExternPaths when computing extern_path options for dependent
// packages.
//
// The label encodes: Pkg = proto package name, Name = crate name.
func (p *ProtocGenProstPlugin) registerExternPaths(lib protoc.ProtoLibrary) {
	crateName := lib.BaseName() + "_rs"

	for _, f := range lib.Files() {
		pkg := f.Package()
		if pkg.Name == "" {
			continue
		}

		protoFile := path.Join(f.Dir, f.Basename)
		protoc.GlobalResolver().Provide(
			"proto",
			"prost_extern",
			protoFile,
			label.New("", pkg.Name, crateName),
		)
	}
}

// resolveTransitiveExternPaths walks the transitive dependency graph of proto
// files and builds extern_path option strings for each dependency package.
func (p *ProtocGenProstPlugin) resolveTransitiveExternPaths(r *rule.Rule, from label.Label) []string {
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
