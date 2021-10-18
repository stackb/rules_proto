package protobuf

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

// The following methods are implemented to satisfy the
// https://pkg.go.dev/github.com/bazelbuild/bazel-gazelle/resolve?tab=doc#Resolver
// interface, but are otherwise unused.
func (pl *protobufLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	fs.StringVar(&pl.configFiles, "proto_configs", "", "optional config.yaml file(s) that provide preconfiguration")
	fs.StringVar(&pl.importsInFiles, "proto_imports_in", "", "index files to parse and load symbols from")
	fs.StringVar(&pl.importsOutFile, "proto_imports_out", "", "filename where index should be written")
	fs.StringVar(&pl.repoName, "proto_repo_name", "", "external name of this repository")
	fs.BoolVar(&pl.overrideGoGooleapis, "override_go_googleapis", false, "if true, remove hardcoded proto_library deps on go_googleapis")

	registerWellKnownProtos(protoc.GlobalResolver())
}

func (pl *protobufLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
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
			if err := protoc.GlobalResolver().LoadFile(filename); err != nil {
				return fmt.Errorf("loading %s: %w", filename, err)
			}
		}
	}

	return nil
}

func (*protobufLang) KnownDirectives() []string {
	return []string{
		protoc.LanguageDirective,
		protoc.PluginDirective,
		protoc.RuleDirective,
	}
}

// Configure implements config.Configurer
func (pl *protobufLang) Configure(c *config.Config, rel string, f *rule.File) {
	if rel == "" {
		// if this is the root BUILD file, we are beginning the configuration
		// sequence.  Perform the equivalent of writing relevant
		// 'gazelle:resolve proto IMP LABEL` entries.
		protoc.GlobalResolver().Install(c)

		// some special handling for certain directives
		if f != nil {
			for _, d := range f.Directives {
				switch d.Key {
				case "prefix":
					// encode the prefix in the resolver.  The name is not used, but
					// the string 'go' is used to reflect the language of origin.
					protoc.GlobalResolver().Provide("gazelle", "directive", "prefix", label.New("", d.Value, "go"))
				}
			}
		}
	}

	if f == nil {
		return
	}

	if err := pl.getOrCreatePackageConfig(c).ParseDirectives(rel, f.Directives); err != nil {
		log.Fatalf("error while parsing rule directives in package %q: %v", rel, err)
	}
}

// getOrCreatePackageConfig either inserts a new config into the map under the
// language name or replaces it with a clone.
func (pl *protobufLang) getOrCreatePackageConfig(config *config.Config) *protoc.PackageConfig {
	var cfg *protoc.PackageConfig
	if existingExt, ok := config.Exts[pl.name]; ok {
		cfg = existingExt.(*protoc.PackageConfig).Clone()
	} else {
		cfg = protoc.NewPackageConfig(config)
	}
	config.Exts[pl.name] = cfg
	return cfg
}

func registerWellKnownProtos(resolver protoc.ImportResolver) {
	for k, v := range map[string]label.Label{
		"google/protobuf/any.proto":             label.New("com_google_protobuf", "", "any_proto"),
		"google/protobuf/api.proto":             label.New("com_google_protobuf", "", "api_proto"),
		"google/protobuf/compiler/plugin.proto": label.New("com_google_protobuf", "", "compiler_plugin_proto"),
		"google/protobuf/descriptor.proto":      label.New("com_google_protobuf", "", "descriptor_proto"),
		"google/protobuf/duration.proto":        label.New("com_google_protobuf", "", "duration_proto"),
		"google/protobuf/empty.proto":           label.New("com_google_protobuf", "", "empty_proto"),
		"google/protobuf/field_mask.proto":      label.New("com_google_protobuf", "", "field_mask_proto"),
		"google/protobuf/source_context.proto":  label.New("com_google_protobuf", "", "source_context_proto"),
		"google/protobuf/struct.proto":          label.New("com_google_protobuf", "", "struct_proto"),
		"google/protobuf/timestamp.proto":       label.New("com_google_protobuf", "", "timestamp_proto"),
		"google/protobuf/type.proto":            label.New("com_google_protobuf", "", "type_proto"),
		"google/protobuf/wrappers.proto":        label.New("com_google_protobuf", "", "wrappers_proto"),
	} {
		resolver.Provide("proto", "proto", k, v)
	}
}
