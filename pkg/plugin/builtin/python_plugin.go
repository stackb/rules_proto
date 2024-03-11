package builtin

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
	protoc.Plugins().MustRegisterPlugin(&PythonPlugin{})
}

// PythonPlugin implements Plugin for the built-in protoc python plugin.
type PythonPlugin struct{}

// Name implements part of the Plugin interface.
func (p *PythonPlugin) Name() string {
	return "builtin:python"
}

// Configure implements part of the Plugin interface.
func (p *PythonPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	flags := parsePythonPluginOptions(p.Name(), ctx.PluginConfig.GetFlags())

	pyFiles := protoc.FlatMapFiles(
		pythonGeneratedFileName(ctx.Rel),
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)

	pyOutputs := make([]string, 0, len(pyFiles))
	for _, pyFile := range pyFiles {
		if flags.excludeOutput[filepath.Base(pyFile)] {
			continue
		}
		pyOutputs = append(pyOutputs, pyFile)
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/builtin", "python"),
		Outputs: pyOutputs,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// pythonGeneratedFileName is a utility function that returns a function that
// computes the name of a predicted generated file having the given
// extension(s) relative to the given dir.
func pythonGeneratedFileName(reldir string) func(f *protoc.File) []string {
	return func(f *protoc.File) []string {
		name := strings.ReplaceAll(f.Name, "-", "_")
		if reldir != "" {
			name = path.Join(reldir, name)
		}
		return []string{name + "_pb2.py"}
	}
}

// pythonPluginOptions represents the parsed flag configuration for the
// ProtocGenTsProto implementation.
type pythonPluginOptions struct {
	excludeOutput map[string]bool
}

func parsePythonPluginOptions(kindName string, args []string) *pythonPluginOptions {
	flags := flag.NewFlagSet(kindName, flag.ExitOnError)

	var excludeOutput string
	flags.StringVar(&excludeOutput, "exclude_output", "", "--exclude_output=foo_pb2.py suppresses the file 'foo_pb2.py' from the output list")

	if err := flags.Parse(args); err != nil {
		log.Fatalf("failed to parse flags for %q: %v", kindName, err)
	}
	config := &pythonPluginOptions{
		excludeOutput: make(map[string]bool),
	}
	for _, value := range strings.Split(excludeOutput, ",") {
		config.excludeOutput[value] = true
	}

	return config
}
