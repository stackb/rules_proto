package ts_proto

import (
	"flag"
	"log"
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenTsProto{})
}

// ProtocGenTsProto implements Plugin for the built-in protoc js/library plugin.
type ProtocGenTsProto struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenTsProto) Name() string {
	return "stephenh:ts-proto:protoc-gen-ts-proto"
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenTsProto) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	options := parseProtoTsLibraryOptions(p.Name(), ctx.PluginConfig.GetFlags())
	tsFiles := make([]string, 0)
	for _, file := range ctx.ProtoLibrary.Files() {
		tsFile := file.Name + ".ts"
		if options.excludeOutput[tsFile] {
			continue
		}
		if ctx.Rel != "" {
			tsFile = path.Join(ctx.Rel, tsFile)
		}
		tsFiles = append(tsFiles, tsFile)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/stephenh/ts-proto", "protoc-gen-ts-proto"),
		Outputs: tsFiles,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// protocGenTsProtoOptions represents the parsed flag configuration for the
// ProtocGenTsProto implementation.
type protocGenTsProtoOptions struct {
	excludeOutput map[string]bool
}

func parseProtoTsLibraryOptions(kindName string, args []string) *protocGenTsProtoOptions {
	flags := flag.NewFlagSet(kindName, flag.ExitOnError)

	var excludeOutput string
	flags.StringVar(&excludeOutput, "exclude_output", "", "--exclude_output=foo.ts suppresses the file 'foo.ts' from the output list")

	if err := flags.Parse(args); err != nil {
		log.Fatalf("failed to parse flags for %q: %v", kindName, err)
	}
	config := &protocGenTsProtoOptions{
		excludeOutput: make(map[string]bool),
	}
	for _, value := range strings.Split(excludeOutput, ",") {
		config.excludeOutput[value] = true
	}

	return config
}
