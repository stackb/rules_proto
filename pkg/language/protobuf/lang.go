package protobuf

import (
	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

// NewProtobufLang create a new protobufLang Gazelle extension implementation.
func NewProtobufLang(name string) *protobufLang {
	return &protobufLang{
		name:     name,
		rules:    protoc.Rules(),
		packages: make(map[string]*protoc.Package),
		resolver: protoc.GlobalResolver(),
	}
}

// protobufLang implements language.Language.
type protobufLang struct {
	// name of the extension
	name string
	// the rule registry
	rules protoc.RuleRegistry
	// the packages that we've generated
	packages map[string]*protoc.Package
	// configFiles contains yconfig yaml files to parse.  May be comma-separated.
	configFiles string
	// repoName is the name (if this an external repository)
	repoName string
	// importsOutFile is the name of the file to create.  If "", skip writing
	// the file.
	importsOutFile string
	// importsInFiles is a comma-separated list of files that contains proto
	// index csv content.
	importsInFiles string
	// reresolveKnownProtoImports performs an additional resolve step for go_googleapis deps
	reresolveKnownProtoImports bool
	// the resolver instance used for cross-resolution
	resolver protoc.ImportResolver
	// starlarkRules stores custom starlark proto rule names in the form filename%rulename
	starlarkRules arrayFlags
	// starlarkPlugins stores custom starlark proto plugin names in the form filename%pluginname
	starlarkPlugins arrayFlags
	// protoRustLibraryPackages collects the workspace-relative path of each
	// package-level proto_rust_library. Per-file crates remain standalone path
	// dependencies under excluded _rust directories. Populated in
	// GenerateRules and consumed in DoneGeneratingRules to update the root
	// Cargo.toml [workspace] members list.
	protoRustLibraryPackages []string
	// protoRustPerFilePackageDirs collects the workspace-relative _rust
	// directories that contain standalone per-file crates. These directories
	// are excluded from the root Cargo workspace so each generated manifest can
	// be used directly or as a path dependency.
	protoRustPerFilePackageDirs []string
	// protoRustPerFilePackages collects the package names and workspace-relative
	// paths of standalone per-file crates. They are listed under the root
	// [patch.crates-io] table so Bazel's Cargo lockfile parser can resolve
	// every transitive path package without making these crates workspace members.
	protoRustPerFilePackages []cargoPathDependency
	// vendorAssetLabels collects bazel labels of every generated rule that
	// provides ProtoCompileInfo and should appear in the root
	// `proto_compile_assets` aggregator. Populated in GenerateRules and
	// consumed in DoneGeneratingRules to update the deps list of the
	// vendoring target between the vendor_proto_sources_deps markers.
	vendorAssetLabels []string
	// repoRoot is captured from the first GenerateRules call so
	// DoneGeneratingRules (which receives no config) can locate the root
	// Cargo.toml and BUILD.bazel.
	repoRoot string
}

// Name implements part of the language.Language interface.
func (pl *protobufLang) Name() string { return pl.name }
