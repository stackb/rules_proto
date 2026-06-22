package prost

import (
	"path"
	"sort"
	"strings"

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

	perFile := ctx.IsProtoFileMode()
	outputs := p.outputs(ctx.ProtoLibrary, perFile)
	if len(outputs) == 0 {
		return nil
	}

	p.registerExternPaths(ctx.ProtoLibrary, perFile)

	options := ctx.PluginConfig.GetOptions()
	if perFile {
		// In per-file mode the plugin needs to be told to emit one .rs
		// per .proto file with the file-stemmed naming convention.
		// `per_file=true` triggers the matching branch in our vendored
		// protoc-gen-prost; see bazel_tools/rust/vendor/README.md in the
		// downstream repo for the rationale.
		options = append(options, "per_file=true")
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/neoeinstein/prost", "protoc-gen-prost"),
		Outputs: outputs,
		Options: options,
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

// HasGeneratableContent reports whether `f` defines any top-level construct
// that would cause prost / pbjson / tonic to emit non-trivial output.
//
// Specifically: real messages (not `extend` blocks — emicklei/proto reports
// `extend` as a `proto.Message` with `IsExtend=true`, but downstream
// generators emit nothing for them), enums, or services.
//
// Used by the per-file output computations to avoid declaring a serde/tonic
// output for a file that only contains `extend` blocks: pbjson-build /
// tonic-build would produce no file, and the missing artifact later trips
// proto_compile.bzl's `mv` step.
func HasGeneratableContent(f *protoc.File) bool {
	if f.HasEnums() || f.HasServices() {
		return true
	}
	for _, msg := range f.Messages() {
		if !msg.IsExtend {
			return true
		}
	}
	return false
}

// outputs computes the output files for the plugin.
//
// In per-package mode (the default), prost generates one .rs file per proto
// package, named {proto_package}.rs.
//
// In per-file mode (used when the surrounding bazel package opts into
// `gazelle:proto file`), prost generates one .rs file per .proto file in
// `file_to_generate`, named `<file_stem>.<proto_package>.rs` so that
// multiple files in the same package don't collide. The vendored
// protoc-gen-prost honours this naming when `per_file=true` is set on
// its options.
//
// Service-only files (no messages/enums) are included — prost still emits
// a stub .rs containing the @@protoc_insertion_point(module) marker, which
// tonic relies on to inject client/server code via append_to_file.
func (p *ProtocGenProstPlugin) outputs(lib protoc.ProtoLibrary, perFile bool) []string {
	outputs := make([]string, 0)

	if perFile {
		for _, f := range lib.Files() {
			if !(f.HasMessages() || f.HasEnums() || f.HasServices()) {
				continue
			}
			pkg := f.Package()
			if pkg.Name == "" {
				continue
			}
			stem := strings.TrimSuffix(f.Basename, ".proto")
			filename := stem + "." + pkg.Name + ".rs"
			if f.Dir != "" {
				filename = path.Join(f.Dir, filename)
			}
			outputs = append(outputs, filename)
		}
		sort.Strings(outputs)
		return outputs
	}

	seen := make(map[string]bool)
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
// each proto file in the library. The data is later consumed by
// ResolveTransitiveExternPaths when computing extern_path options for
// dependent packages.
//
// Two flavours are registered:
//
//   - "prost_extern" / <proto file path> → (proto pkg, **façade** crate name).
//     Used to map a file-import to the rust crate that consumers should depend
//     on. Always the per-package façade — not a per-file crate — because
//     cross-package extern_path entries are prefix-matched by prost-build,
//     which then runs `to_snake_case` on each segment of the rust_path. That
//     normalization collapses `__` → `_`, so a per-file rust_path like
//     `::<facade>__<stem>` would land as the nonexistent `::<facade>_<stem>`.
//     The façade transparently re-exports the per-file crates' contents
//     (see proto_rust_library.bzl's facade emission), so the consumer's
//     code path stays identical.
//
//   - "prost_extern_type" / <proto package> → (type name, **per-file** crate
//     name). Only populated in per-file mode. Used by
//     ResolveTransitiveExternPaths to emit per-type extern_path entries for
//     same-package cross-file references — these are EXACT-match resolved by
//     prost-build (no to_snake) so the double-underscored per-file crate name
//     survives intact. Without these, prost emits a relative `super::...`
//     path that won't resolve once each .proto lives in its own crate.
//
// Crate-name choices: per-package mode uses `RustCrateName(pkg.Name)` —
// e.g. `omnistac.spok.message` → `omnistac_spok_message`. The per-file
// crates (registered only under PerFileTypeProvideKind) add a `__<file_stem>`
// suffix — `omnistac_spok_message__order`. That suffix matches the
// proto_rust_library names the gazelle plugin emits for per-file crates.
func (p *ProtocGenProstPlugin) registerExternPaths(lib protoc.ProtoLibrary, perFile bool) {
	resolver := protoc.GlobalResolver()
	for _, f := range lib.Files() {
		pkg := f.Package()
		if pkg.Name == "" {
			continue
		}

		facadeCrate := protoc.RustCrateName(pkg.Name)
		protoFile := path.Join(f.Dir, f.Basename)

		resolver.Provide(
			"proto",
			"prost_extern",
			protoFile,
			label.New("", pkg.Name, facadeCrate),
		)

		// In per-file mode, also register per-type entries so sibling per-file
		// crates in the same proto package can map their cross-file type
		// references to the correct sibling crate. Without this, prost emits
		// a relative path like `Other` assuming nested-module layout, which
		// fails to resolve once Other lives in a different crate.
		//
		// The key is the proto package name (not the fully-qualified type
		// path), so that ResolveTransitiveExternPaths can enumerate every
		// type in a package by a single Resolve() call. The label's `Pkg`
		// holds the type name; `Name` holds the crate that defines it.
		if !perFile {
			continue
		}
		stem := strings.TrimSuffix(f.Basename, ".proto")
		perFileCrate := facadeCrate + "__" + stem
		for _, msg := range f.Messages() {
			resolver.Provide(
				"proto",
				PerFileTypeProvideKind,
				pkg.Name,
				label.New("", msg.Name, perFileCrate),
			)
		}
		for _, enum := range f.Enums() {
			resolver.Provide(
				"proto",
				PerFileTypeProvideKind,
				pkg.Name,
				label.New("", enum.Name, perFileCrate),
			)
		}
	}
}

// PerFileTypeProvideKind is the GlobalResolver `kind` arg under which
// per-file mode registers each message/enum a per-file proto_library
// contributes. Keyed by proto package, the resulting Resolve() call returns
// the (type name, crate name) for every type declared anywhere in that
// proto package across sibling per-file libraries.
const PerFileTypeProvideKind = "prost_extern_per_file_type"
