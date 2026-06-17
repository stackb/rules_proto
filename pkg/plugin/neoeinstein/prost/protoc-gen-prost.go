package prost

import (
	"path"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

const (
	ProtocGenProstPluginName = "neoeinstein:prost:protoc-gen-prost"
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
		Label:   label.New("build_stack_rules_proto", "plugin/neoeinstein/prost", "protoc-gen-prost"),
		Outputs: outputs,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// ResolvePluginOptions implements the PluginOptionsResolver interface. It
// computes extern_path options based on transitive proto file dependencies,
// and prepends compile_well_known_types=true whenever:
//
//   - one of the library's own proto packages is google.protobuf (the library
//     is itself compiling the well-known types — without this flag prost
//     skips them and emits a stub), or
//   - the resolved extern_path set references .google.protobuf (the library
//     consumes well-known types via a foreign crate; prost-build registers a
//     default extern path to ::prost_types unless this flag clears it, which
//     would collide with our extern_path entry and error out as "duplicate
//     extern Protobuf path").
//
// Only the prost plugin needs this — protoc-gen-prost-serde (pbjson-build)
// doesn't consult the flag and its codegen doesn't hit the collision path.
func (p *ProtocGenProstPlugin) ResolvePluginOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	opts := ResolveExternPathOptions(cfg, r, from)
	if needsCompileWellKnownTypes(opts, ownProtoPackages(r, from)) {
		opts = append([]string{"compile_well_known_types=true"}, opts...)
	}
	return opts
}

const wellKnownTypesProtoPackage = "google.protobuf"

// needsCompileWellKnownTypes reports whether the prost plugin should emit
// compile_well_known_types=true for a library based on its computed options
// and own-package set. See ResolvePluginOptions for the rationale.
func needsCompileWellKnownTypes(opts []string, ownPackages map[string]bool) bool {
	if ownPackages[wellKnownTypesProtoPackage] {
		return true
	}
	for _, opt := range opts {
		if opt == "extern_path=."+wellKnownTypesProtoPackage+"=::"+protoc.RustCrateName(wellKnownTypesProtoPackage) {
			return true
		}
		// More general guard: any extern_path mapping for .google.protobuf
		// suppresses prost-build's default registration to avoid collision.
		const prefix = "extern_path=." + wellKnownTypesProtoPackage + "="
		if len(opt) > len(prefix) && opt[:len(prefix)] == prefix {
			return true
		}
	}
	return false
}

// shouldApply returns true if the library has any file with messages, enums,
// or services. Services are included because protoc-gen-tonic's generated
// code is inserted into prost's {package}.rs at the @@protoc_insertion_point
// marker (see protoc-gen-prost's append_to_file mechanism). If prost isn't
// invoked, the file doesn't exist and tonic's insert fails with
// "Tried to insert into file that doesn't exist".
func (p *ProtocGenProstPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() || f.HasServices() {
			return true
		}
	}
	return false
}

// outputs computes the output files for the plugin. Prost generates one .rs
// file per proto package, named {proto_package}.rs. The path includes the
// file's directory so that mergeSources can handle the rel stripping.
//
// Packages contributed by service-only files (no messages/enums) are
// included — prost still emits a stub .rs containing the
// @@protoc_insertion_point(module) marker, which tonic relies on to inject
// its client/server code via append_to_file.
func (p *ProtocGenProstPlugin) outputs(lib protoc.ProtoLibrary) []string {
	seen := make(map[string]bool)
	outputs := make([]string, 0)

	for _, f := range lib.Files() {
		if !(f.HasMessages() || f.HasEnums() || f.HasServices()) {
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
// ResolveTransitiveExternPaths when computing extern_path options for dependent
// packages.
//
// The label encodes: Pkg = proto package name, Name = crate name. The crate
// name comes from protoc.RustCrateName so it matches the rust_library target
// name produced by RustLibrary.Name() — without this alignment, downstream
// extern_path entries would point at a non-existent crate and rustc would
// fail to resolve types.
func (p *ProtocGenProstPlugin) registerExternPaths(lib protoc.ProtoLibrary) {
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
			label.New("", pkg.Name, protoc.RustCrateName(pkg.Name)),
		)
	}
}
