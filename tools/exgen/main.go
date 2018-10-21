// Wsifier, a tool to parse BUILD files and bzl files, generate tests cases and
// documentation.
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/urfave/cli"
)

// Language represents one directory in this repo
type Language struct {
	// Directory in the repo where this language is rooted.  Typically this is
	// the same as the name
	Dir string

	// Name of the language
	Name string

	// Workspace usage
	Usage string

	// List of rules
	Rules []*Rule

	// List of available plugins
	Plugins []*Plugin
}

func (l *Language) GetPluginByName(name string) *Plugin {
	for _, p := range l.Plugins {
		if p.Name == name {
			return p
		}
	}
	return nil
}

type Rule struct {
	Name           string
	Type           string
	Doc            string
	Usage          *template.Template
	Example        *template.Template
	Implementation *template.Template
	Attrs          []*Attr
	Plugins        []string
}

type Attr struct {
	Name      string
	Type      string
	Default   string
	Doc       string
	Mandatory bool
}

type Plugin struct {
	Name string
}

type ruleData struct {
	Lang *Language
	Rule *Rule
}

func main() {
	app := cli.NewApp()
	app.Name = "exgen"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "dir",
			Usage: "Directory to scan",
			Value: ".",
		},
	}
	app.Action = func(c *cli.Context) error {
		err := action(c)
		if err != nil {
			return cli.NewExitError("%v", 1)
		}
		return nil
	}

	app.Run(os.Args)
}

func action(c *cli.Context) error {
	dir := c.String("dir")
	if dir == "" {
		return fmt.Errorf("--dir required")
	}

	languages := []*Language{
		makeCpp(),
		makeClosure(),
		makeNode(),
		makeJava(),
		makeScala(),
		makePython(),
		makeGo(),
		makeRuby(),
		makeGithubComGrpcGrpcWeb(),
	}

	for _, lang := range languages {
		mustWriteReadme(dir, lang)
		mustWriteRules(dir, lang)
		mustWriteExamples(dir, lang)
	}

	mustWriteMakefile(dir, languages)

	return nil
}

var protoCompileAttrs = []*Attr{
	&Attr{
		Name:      "deps",
		Type:      "list<ProtoInfo>",
		Default:   "[]",
		Doc:       "List of labels that provide a `ProtoInfo` (`native.proto_library`)",
		Mandatory: true,
	},
	&Attr{
		Name:      "plugins",
		Type:      "list<ProtoPluginInfo>",
		Default:   "[]",
		Doc:       "List of labels that provide a `ProtoPluginInfo`",
		Mandatory: false,
	},
	&Attr{
		Name:      "plugin_options",
		Type:      "list<string>",
		Default:   "[]",
		Doc:       "List of additional 'global' plugin options (applies to all plugins)",
		Mandatory: false,
	},
	&Attr{
		Name:      "outputs",
		Type:      "list<generated file>",
		Default:   "[]",
		Doc:       "List of additional expected generated file outputs",
		Mandatory: false,
	},
	&Attr{
		Name:      "has_services",
		Type:      "bool",
		Default:   "False",
		Doc:       "If the proto files(s) have a service rpc, generate grpc outputs",
		Mandatory: false,
	},
	&Attr{
		Name:      "protoc",
		Type:      "executable file",
		Default:   "@com_google_protobuf//:protoc",
		Doc:       "The protocol compiler tool",
		Mandatory: false,
	},
	&Attr{
		Name:      "verbose",
		Type:      "int",
		Default:   "0",
		Doc:       "1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*",
		Mandatory: false,
	},
	&Attr{
		Name:      "include_imports",
		Type:      "bool",
		Default:   "True",
		Doc:       "Pass the --include_imports argument to the protoc_plugin",
		Mandatory: false,
	},
	&Attr{
		Name:      "include_source_info",
		Type:      "bool",
		Default:   "True",
		Doc:       "Pass the --include_source_info argument to the protoc_plugin",
		Mandatory: false,
	},
	&Attr{
		Name:      "transitive",
		Type:      "bool",
		Default:   "False",
		Doc:       "Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies",
		Mandatory: false,
	},
}

var compileRuleTemplate = mustTemplate(`load("//:compile.bzl", "proto_compile")

def {{ .Rule.Name }}(**kwargs):
    proto_compile(
		plugins = [{{ range .Rule.Plugins }}
			str(Label("{{ . }}")),{{ end}}
		],
        **kwargs
	)`)

var usageTemplate = mustTemplate(`# WORKSPACE
load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()`)

var grpcUsageTemplate = mustTemplate(`# WORKSPACE
load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()`)

var protoCompileExampleTemplate = mustTemplate(`# BUILD.bazel
load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
	name = "person_{{ .Lang.Name }}_proto",
	deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

var grpcCompileExampleTemplate = mustTemplate(`# BUILD.bazel
load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
	name = "greeter_{{ .Lang.Name }}_grpc",
	deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)`)

var protoLibraryExampleTemplate = mustTemplate(`# BUILD.bazel
load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
	name = "person_{{ .Lang.Name }}_library",
	deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

var grpcLibraryExampleTemplate = mustTemplate(`# BUILD.bazel
load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
	name = "greeter_{{ .Lang.Name }}_library",
	deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)`)

func mustWriteRules(dir string, lang *Language) {
	for _, rule := range lang.Rules {
		mustWriteRule(dir, lang, rule)
	}
}

func mustWriteRule(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	out.t(rule.Implementation, &ruleData{lang, rule})
	out.ln()
	out.MustWrite(path.Join(dir, lang.Dir, rule.Name+".bzl"))
}

func mustWriteExamples(dir string, lang *Language) {
	for _, rule := range lang.Rules {
		exampleDir := path.Join(dir, lang.Dir, "example", rule.Name)
		os.MkdirAll(exampleDir, os.ModePerm)
		mustWriteExampleWorkspace(exampleDir, lang, rule)
		mustWriteExampleBuildFile(exampleDir, lang, rule)
	}
}

func mustWriteExampleWorkspace(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	depth := strings.Split(lang.Dir, "/")
	// +2 as we are in the example/{rule} subdirectory
	relpath := strings.Repeat("../", len(depth)+2)

	out.w(`
http_archive(
	name = "bazel_toolchains",
	urls = [
		"https://mirror.bazel.build/github.com/bazelbuild/bazel-toolchains/archive/bc09b995c137df042bb80a395b73d7ce6f26afbe.tar.gz",
		"https://github.com/bazelbuild/bazel-toolchains/archive/bc09b995c137df042bb80a395b73d7ce6f26afbe.tar.gz",
	],
	strip_prefix = "bazel-toolchains-bc09b995c137df042bb80a395b73d7ce6f26afbe",
	sha256 = "4329663fe6c523425ad4d3c989a8ac026b04e1acedeceb56aa4b190fa7f3973c",
)

local_repository(
	name = "build_stack_rules_proto",
	path = "%s",
)
`, relpath)

	out.ln()
	out.t(rule.Usage, &ruleData{lang, rule})
	out.ln()
	out.MustWrite(path.Join(dir, "WORKSPACE"))
}

func mustWriteExampleBuildFile(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	out.t(rule.Example, &ruleData{lang, rule})
	out.ln()
	out.MustWrite(path.Join(dir, "BUILD.bazel"))
}

func mustWriteMakefile(dir string, languages []*Language) {
	out := &LineWriter{}

	for _, lang := range languages {
		ruleNames := make([]string, len(lang.Rules))
		for i, rule := range lang.Rules {
			ruleNames[i] = rule.Name

			out.w("%s: ", rule.Name)
			out.w("\t(cd %s && /home/pcj/.cache/bzl/release/0.17.2/bin/bazel --bazelrc /home/pcj/go/src/github.com/stackb/rules_proto/tools/bazelrc.remote build //...)", path.Join(lang.Dir, "example", rule.Name))
			out.ln()
		}
		out.w("%s: %s", lang.Name, strings.Join(ruleNames, " "))
		out.ln()
	}

	out.MustWrite(path.Join(dir, "Makefile.examples"))
}

func mustWriteReadme(dir string, lang *Language) {
	out := &LineWriter{}

	out.w("# Language `%s`", lang.Name)
	out.ln()

	for _, rule := range lang.Rules {
		out.w("## `%s`", rule.Name)
		out.ln()
		out.w(rule.Doc)
		out.ln()

		out.w("### Usage")
		out.ln()

		out.w("```python")
		out.t(rule.Usage, &ruleData{lang, rule})
		out.w("```")
		out.ln()

		out.w("### Example")
		out.ln()

		out.w("```python")
		out.t(rule.Example, &ruleData{lang, rule})
		out.w("```")
		out.ln()

		out.w("### Mandatory Attributes")
		out.ln()
		out.w("| Name | Type | Default | Description |")
		out.w("| ---: | :--- | ------- | ----------- |")
		for _, attr := range rule.Attrs {
			if attr.Mandatory {
				out.w("| %s   | `%s` | `%s`    | %s          |", attr.Name, attr.Type, attr.Default, attr.Doc)
			}
		}
		out.ln()

		out.w("### Optional Attributes")
		out.ln()
		out.w("| Name | Type | Default | Description |")
		out.w("| ---: | :--- | ------- | ----------- |")
		for _, attr := range rule.Attrs {
			if !attr.Mandatory {
				out.w("| %s   | `%s` | `%s`    | %s          |", attr.Name, attr.Type, attr.Default, attr.Doc)
			}
		}
		out.ln()

	}

	out.ln()

	out.MustWrite(path.Join(dir, lang.Dir, "README.md"))
}

// ********************************
// Utility types
// ********************************

type LineWriter struct {
	lines []string
}

func (w *LineWriter) w(s string, args ...interface{}) {
	w.lines = append(w.lines, fmt.Sprintf(s, args...))
}

func (w *LineWriter) t(t *template.Template, data interface{}) {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		log.Fatalf("%v", err)
	}
	w.lines = append(w.lines, buf.String())
}

func (w *LineWriter) ln() {
	w.lines = append(w.lines, "")
}

func (w *LineWriter) MustWrite(filepath string) {
	err := writeFile(filepath, strings.Join(w.lines, "\n"))
	if err != nil {
		log.Fatalf("FAIL %s: %v", filepath, err)
	}
}

// ********************************
// Utility functions
// ********************************

func mustTemplate(tpl string) *template.Template {
	return template.Must(template.New("").Option("missingkey=error").Parse(tpl))
}

func writeFile(filepath, content string) error {
	fo, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer fo.Close()

	_, err = io.Copy(fo, strings.NewReader(content))
	if err != nil {
		return err
	}

	log.Printf("Wrote %s", filepath)
	return nil
}
