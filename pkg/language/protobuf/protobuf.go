package protobuf

import (
	"flag"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

// NewProtobufLanguage create a new ProtobufLanguage Gazelle extension implementation.
func NewProtobufLanguage(name string) *ProtobufLanguage {
	return &ProtobufLanguage{
		name:      name,
		rules:     protoc.Rules(),
		providers: make(map[label.Label]protoc.RuleProvider),
	}
}

// ProtobufLanguage implements language.Language.
type ProtobufLanguage struct {
	// name of the extension
	name string
	// the rule registry
	rules protoc.RuleRegistry
	// configFiles contains yconfig yaml files to parse.  May be comma-separated.
	configFiles string
	// providers is a mapping from label -> the provider that produced the rule.
	// we save this in the config such that we can retrieve the association
	// later in the resolve step.
	providers map[label.Label]protoc.RuleProvider
	// repoName is the name (if this an external repository)
	repoName string
	// importsOutFile is the name of the file to create.  If "", skip writing
	// the file.
	importsOutFile string
	// importsInFiles is a comma-separated list of files that contains proto
	// index csv content.
	importsInFiles string
	// overrideGoGooleapis performs special processing for go_googleapis deps
	overrideGoGooleapis bool
}

// Name returns the name of the language. This should be a prefix of the kinds
// of rules generated by the language, e.g., "go" for the Go extension since it
// generates "go_library" rules.
func (pl *ProtobufLanguage) Name() string { return pl.name }

// The following methods are implemented to satisfy the
// https://pkg.go.dev/github.com/bazelbuild/bazel-gazelle/resolve?tab=doc#Resolver
// interface, but are otherwise unused.
func (pl *ProtobufLanguage) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	fs.StringVar(&pl.configFiles, "proto_configs", "", "optional config.yaml file(s) that provide preconfiguration")
	fs.StringVar(&pl.importsInFiles, "proto_imports_in", "", "index files to parse and load symbols from")
	fs.StringVar(&pl.importsOutFile, "proto_imports_out", "", "filename where index should be written")
	fs.StringVar(&pl.repoName, "proto_repo_name", "", "external name of this repository")
	fs.BoolVar(&pl.overrideGoGooleapis, "override_go_googleapis", false, "if true, remove hardcoded proto_library deps on go_googleapis")
}

func (pl *ProtobufLanguage) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	cfg := protoc.NewPackageConfig(c)
	c.Exts[pl.name] = cfg

	if pl.configFiles != "" {
		for _, filename := range strings.Split(pl.configFiles, ",") {
			if err := protoc.LoadYConfigFile(c, cfg, filename); err != nil {
				return fmt.Errorf("loading -proto_configs %s: %w", filename, err)
			}
		}
	}

	if pl.importsInFiles != "" {
		for _, filename := range strings.Split(pl.importsInFiles, ",") {
			if err := protoc.GlobalResolver().LoadImportsFile(filename); err != nil {
				return fmt.Errorf("loading %s: %w", filename, err)
			}
		}
	}

	return nil
}

func (*ProtobufLanguage) KnownDirectives() []string {
	return []string{
		protoc.LanguageDirective,
		protoc.PluginDirective,
		protoc.RuleDirective,
	}
}

// Configure implements config.Configurer
func (pl *ProtobufLanguage) Configure(c *config.Config, rel string, f *rule.File) {
	if rel == "" {
		// if this is the root BUILD file, we are beginning the configuration
		// sequence.  Perform the equivalent of writing relevant
		// 'gazelle:resolve proto IMP LABEL` entries.
		protoc.GlobalResolver().InstallResolveOverrides(c)
	}

	if f == nil {
		return
	}
	if err := pl.getOrCreatePackageConfig(c).ParseDirectives(rel, f.Directives); err != nil {
		log.Fatalf("error while parsing rule directives in package %q: %v", rel, err)
	}
}

// Kinds returns a map of maps rule names (kinds) and information on how to
// match and merge attributes that may be found in rules of those kinds. All
// kinds of rules generated for this language may be found here.
func (*ProtobufLanguage) Kinds() map[string]rule.KindInfo {
	registry := protoc.Rules()

	kinds := make(map[string]rule.KindInfo)
	kinds[overrideKindName] = overrideKind

	for _, name := range registry.RuleNames() {
		rule, err := registry.LookupRule(name)
		if err != nil {
			log.Fatal("Kinds:", err)
		}
		kinds[rule.Name()] = rule.KindInfo()
	}

	return kinds
}

// Loads returns .bzl files and symbols they define. Every rule generated by
// GenerateRules, now or in the past, should be loadable from one of these
// files.
func (pl *ProtobufLanguage) Loads() []rule.LoadInfo {
	// Merge symbols
	symbolsByLoadName := make(map[string][]string)
	for _, name := range pl.rules.RuleNames() {
		rule, err := pl.rules.LookupRule(name)
		if err != nil {
			log.Fatal(err)
		}
		load := rule.LoadInfo()
		symbolsByLoadName[load.Name] = append(symbolsByLoadName[load.Name], load.Symbols...)
	}

	// Ensure names are sorted otherwise order of load statements can be
	// non-deterministic
	keys := make([]string, 0)
	for name := range symbolsByLoadName {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	// Build final load list
	loads := make([]rule.LoadInfo, 0)
	for _, name := range keys {
		symbols := symbolsByLoadName[name]
		sort.Strings(symbols)
		loads = append(loads, rule.LoadInfo{
			Name:    name,
			Symbols: symbols,
		})
	}
	return loads
}

// Fix repairs deprecated usage of language-specific rules in f. This is called
// before the file is indexed. Unless c.ShouldFix is true, fixes that delete or
// rename rules should not be performed.
func (*ProtobufLanguage) Fix(c *config.Config, f *rule.File) {}

// Imports returns a list of ImportSpecs that can be used to import the rule r.
// This is used to populate RuleIndex.
//
// If nil is returned, the rule will not be indexed. If any non-nil slice is
// returned, including an empty slice, the rule will be indexed.
func (pl *ProtobufLanguage) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	srcs := r.AttrStrings("srcs")
	imports := make([]resolve.ImportSpec, len(srcs))

	for i, src := range srcs {
		imports[i] = resolve.ImportSpec{
			// Lang is the language in which the import string appears (this
			// should match Resolver.Name).
			Lang: pl.name,
			// Imp is an import string for the library.
			Imp: fmt.Sprintf("//%s:%s", f.Pkg, src),
		}
	}

	return imports
}

// Embeds returns a list of labels of rules that the given rule embeds. If a
// rule is embedded by another importable rule of the same language, only the
// embedding rule will be indexed. The embedding rule will inherit the imports
// of the embedded rule. Since SkyLark doesn't support embedding this should
// always return nil.
func (*ProtobufLanguage) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }

// Resolve translates imported libraries for a given rule into Bazel
// dependencies. Information about imported libraries is returned for each rule
// generated by language.GenerateRules in language.GenerateResult.Imports.
// Resolve generates a "deps" attribute (or the appropriate language-specific
// equivalent) for each import according to language-specific rules and
// heuristics.
func (pl *ProtobufLanguage) Resolve(
	c *config.Config,
	ix *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	importsRaw interface{},
	from label.Label,
) {
	if r.Kind() == overrideKindName {
		resolveOverrideRule(from.Pkg, r)
		return
	}

	if provider, ok := pl.providers[from]; ok {
		cfg := pl.getOrCreatePackageConfig(c)
		if imports, ok := importsRaw.([]string); ok {
			provider.Resolve(cfg, r, imports, from)
		} else {
			log.Panicf("warning: resolve imports: expected []string, got %T", importsRaw)
		}
	}
}

func (*ProtobufLanguage) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	return protoc.GlobalResolver().CrossResolve(c, ix, imp, lang)
}

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested in
// depth-first post-order.
//
// args contains the arguments for GenerateRules. This is passed as a struct to
// avoid breaking implementations in the future when new fields are added.
//
// A GenerateResult struct is returned. Optional fields may be added to this
// type in the future.
//
// Any non-fatal errors this function encounters should be logged using
// log.Print.
func (pl *ProtobufLanguage) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	cfg := pl.getOrCreatePackageConfig(args.Config)

	files := make(map[string]*protoc.File)
	for _, f := range args.RegularFiles {
		if !protoc.IsProtoFile(f) {
			continue
		}
		file := protoc.NewFile(args.Rel, f)
		if err := file.Parse(); err != nil {
			log.Fatalf("unparseable proto file dir=%s, file=%s: %v", args.Dir, file.Basename, err)
		}
		files[f] = file
	}

	protoLibraries := make([]protoc.ProtoLibrary, 0)
	for _, r := range args.OtherGen {
		if r.Kind() != "proto_library" {
			continue
		}
		srcs := r.AttrStrings("srcs")
		srcLabels := make([]label.Label, len(srcs))
		for i, src := range srcs {
			srcLabel, err := label.Parse(src)
			if err != nil {
				log.Fatalf("%s %q: unparseable source label %q: %v", r.Kind(), r.Name(), src, err)
			}
			srcLabels[i] = srcLabel
		}
		lib := protoc.NewOtherProtoLibrary(args.File, r, matchingFiles(files, srcLabels)...)
		protoLibraries = append(protoLibraries, lib)
	}

	pkg := protoc.NewPackage(args.Rel, cfg, protoLibraries...)

	for _, provider := range pkg.RuleProviders() {
		labl := label.New(args.Config.RepoName, args.Rel, provider.Name())
		pl.providers[labl] = provider
		// TODO: if needed allow FileVisitor to mutate the rule.File here.
	}

	rules := pkg.Rules()

	otherLibs := make([]*rule.Rule, 0)

	// Seek out proto_library rules and register the sources they provide
	for _, r := range args.OtherGen {
		if r.Kind() == "proto_library" {
			otherLibs = append(otherLibs, r)
			for _, src := range r.AttrStrings("srcs") {
				protoc.GlobalResolver().Provides(
					"proto_library",
					path.Join(args.Rel, src),
					label.New("", args.Rel, r.Name()))
			}
		}
	}

	// special case if we want to override go_googleapis deps.
	if pl.overrideGoGooleapis && len(otherLibs) > 0 {
		rules = append(rules, makeProtoOverrideRule(otherLibs))
	}

	imports := make([]interface{}, len(rules))
	for i, rule := range rules {
		imports[i] = rule.PrivateAttr(config.GazelleImportsKey)
	}

	// special case if this is the root BUILD file and the user requested to
	// write the imports file.
	if args.Rel == "" && pl.importsOutFile != "" {
		if err := protoc.GlobalResolver().SaveImportsFile(pl.importsOutFile, pl.repoName); err != nil {
			log.Printf("error saving import file: %s: %v", pl.importsOutFile, err)
		}
	}

	return language.GenerateResult{
		Gen:     rules,
		Imports: imports,
		Empty:   pkg.Empty(),
	}
}

// getOrCreatePackageConfig either inserts a new config into the map under the
// language name or replaces it with a clone.
func (pl *ProtobufLanguage) getOrCreatePackageConfig(config *config.Config) *protoc.PackageConfig {
	var cfg *protoc.PackageConfig
	if existingExt, ok := config.Exts[pl.name]; ok {
		cfg = existingExt.(*protoc.PackageConfig).Clone()
	} else {
		cfg = protoc.NewPackageConfig(config)
	}
	config.Exts[pl.name] = cfg
	return cfg
}

func matchingFiles(files map[string]*protoc.File, srcs []label.Label) []*protoc.File {
	matching := make([]*protoc.File, 0)
	for _, src := range srcs {
		if file, ok := files[src.Name]; ok {
			matching = append(matching, file)
		}
	}
	return matching
}
