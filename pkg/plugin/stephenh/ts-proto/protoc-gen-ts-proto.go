package ts_proto

import (
	"flag"
	"log"
	"path"
	"path/filepath"
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
	flags := parseProtocGenTsProtoOptions(p.Name(), ctx.PluginConfig.GetFlags())
	imports := make(map[string]bool)
	for _, file := range ctx.ProtoLibrary.Files() {
		for _, imp := range file.Imports() {
			imports[imp.Filename] = true
		}
	}
	var emitImportedFiles bool
	var options []string
	for _, option := range ctx.PluginConfig.GetOptions() {
		// options may be configured to include many "M=" options, but only
		// include the relevant ones to avoid BUILD file clutter.
		if strings.HasPrefix(option, "M=") {
			keyVal := option[len("M="):]
			parts := strings.SplitN(keyVal, "=", 2)
			filename := parts[0]
			if !imports[filename] {
				continue
			}
		}
		if option == "emitImportedFiles=true" {
			emitImportedFiles = true
		}
		options = append(options, option)
	}

	tsFiles := make([]string, 0)
	for _, file := range ctx.ProtoLibrary.Files() {
		if emitImportedFiles {
			for _, imp := range file.Imports() {
				if !strings.HasPrefix(imp.Filename, "google/protobuf") {
					continue
				}
				tsFiles = append(tsFiles, strings.TrimSuffix(imp.Filename, ".proto")+".ts")
			}
		}

		tsFile := file.Name + ".ts"
		if flags.excludeOutput[filepath.Base(tsFile)] {
			continue
		}
		if ctx.Rel != "" {
			tsFile = path.Join(ctx.Rel, tsFile)
		}
		tsFiles = append(tsFiles, tsFile)
	}

	pc := &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/stephenh/ts-proto", "protoc-gen-ts-proto"),
		Outputs: protoc.DeduplicateAndSort(tsFiles),
		Options: protoc.DeduplicateAndSort(options),
	}
	if len(pc.Outputs) == 0 {
		pc.Outputs = nil
	}
	return pc
}

// protocGenTsProtoOptions represents the parsed flag configuration for the
// ProtocGenTsProto implementation.
type protocGenTsProtoOptions struct {
	excludeOutput map[string]bool
}

func parseProtocGenTsProtoOptions(kindName string, args []string) *protocGenTsProtoOptions {
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
