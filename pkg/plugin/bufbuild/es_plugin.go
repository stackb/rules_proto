package bufbuild

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
	protoc.Plugins().MustRegisterPlugin(&EsProto{})
}

// EsProto implements Plugin for the bufbuild/protoc-gen-es plugin.
type EsProto struct{}

// Name implements part of the Plugin interface.
func (p *EsProto) Name() string {
	return "bufbuild:es"
}

// Configure implements part of the Plugin interface.
func (p *EsProto) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	flags := parseEsProtoOptions(p.Name(), ctx.PluginConfig.GetFlags())
	imports := make(map[string]bool)
	for _, file := range ctx.ProtoLibrary.Files() {
		for _, imp := range file.Imports() {
			imports[imp.Filename] = true
		}
	}
	// TODO: get target option from directive
	var options = []string{"keep_empty_files=true", "target=ts"}
	tsFiles := make([]string, 0)
	for _, file := range ctx.ProtoLibrary.Files() {
		// TODO: outputs should be conditional on which target= value is used
		tsFile := file.Name + "_pb.ts"
		if flags.excludeOutput[filepath.Base(tsFile)] {
			continue
		}
		if ctx.Rel != "" {
			tsFile = path.Join(ctx.Rel, tsFile)
		}
		tsFiles = append(tsFiles, tsFile)
	}

	pc := &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/bufbuild", "es"),
		Outputs: protoc.DeduplicateAndSort(tsFiles),
		Options: protoc.DeduplicateAndSort(options),
	}
	if len(pc.Outputs) == 0 {
		pc.Outputs = nil
	}
	return pc
}

// EsProtoOptions represents the parsed flag configuration for the
// EsProto implementation.
type EsProtoOptions struct {
	excludeOutput map[string]bool
}

func parseEsProtoOptions(kindName string, args []string) *EsProtoOptions {
	flags := flag.NewFlagSet(kindName, flag.ExitOnError)

	var excludeOutput string
	flags.StringVar(&excludeOutput, "exclude_output", "", "--exclude_output=foo.ts suppresses the file 'foo.ts' from the output list")

	if err := flags.Parse(args); err != nil {
		log.Fatalf("failed to parse flags for %q: %v", kindName, err)
	}
	config := &EsProtoOptions{
		excludeOutput: make(map[string]bool),
	}
	for _, value := range strings.Split(excludeOutput, ",") {
		config.excludeOutput[value] = true
	}

	return config
}
