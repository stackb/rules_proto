// Wsifier, a tool to parse BUILD files and bzl files, generate tests cases and
// documentation.
package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
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

	// Additional nodes about the language
	Notes *template.Template

	// List of available plugins
	Plugins map[string]*Plugin

	// Does the langaguage has a routeguide server?  If so, this is the bazel target to run it.
	RouteGuideServer, RouteGuideClient string

	// If not the empty string, one-word reason why excluded from TravisCI
	// configuration
	TravisExclusionReason string

	// Additional travis-specific env vars in the form "K=V"
	TravisEnvVars []string
}

type Rule struct {
	// Name of the rule
	Name string

	// Base name of the rule (typically the lang name)
	Base string

	// Kind of the rule (proto|grpc)
	Kind string

	// Description
	Doc string

	// Temmplate for workspace
	Usage *template.Template

	// Template for build file
	Example *template.Template

	// Template for bzl file
	Implementation *template.Template

	// List of attributes
	Attrs []*Attr

	// List of plugins
	Plugins []string

	// Not expected to be functional
	Experimental bool

	// Not compatible with remote execution
	RemoteIncompatible bool

	// Bazel build flags
	Flags []*Flag

	// If not the empty string, one-word reason why excluded from TravisCI
	// configuration
	TravisExclusionReason string

	// Additional travis-specific env vars in the form "K=V"
	TravisEnvVars []string
}

// Flag captures information about a bazel build flag.
type Flag struct {
	Category    string
	Name        string
	Value       string
	Description string
}

type Attr struct {
	Name      string
	Type      string
	Default   string
	Doc       string
	Mandatory bool
}

type Plugin struct {
	Tool    string
	Options []string
}

type ruleData struct {
	Lang *Language
	Rule *Rule
}

func main() {
	app := cli.NewApp()
	app.Name = "rulegen"
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:  "dir",
			Usage: "Directory to scan",
			Value: ".",
		},
		&cli.StringFlag{
			Name:  "header",
			Usage: "Template for the main readme header",
			Value: "tools/rulegen/README.header.md",
		},
		&cli.StringFlag{
			Name:  "footer",
			Usage: "Template for the main readme footer",
			Value: "tools/rulegen/README.footer.md",
		},
		&cli.StringFlag{
			Name:  "travis_header",
			Usage: "Template for the travis header",
			Value: "tools/rulegen/travis.header.yml",
		},
		&cli.StringFlag{
			Name:  "travis_footer",
			Usage: "Template for the travis footer",
			Value: "tools/rulegen/travis.footer.yml",
		},
		&cli.StringFlag{
			Name:  "ref",
			Usage: "Version ref to use for main readme",
			Value: "{GIT_COMMIT_ID}",
		},
		&cli.StringFlag{
			Name:  "sha256",
			Usage: "Sha256 value to use for main readme",
			Value: "{ARCHIVE_TAR_GZ_SHA256}",
		},
		&cli.StringFlag{
			Name:  "github_url",
			Usage: "URL for github download",
			Value: "https://github.com/stackb/rules_proto/archive/{ref}.tar.gz",
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

	ref := c.String("ref")
	sha256 := c.String("sha256")
	githubURL := c.String("github_url")

	// Autodetermine sha256 if we have a real commit and templated sha256 value
	if ref != "{GIT_COMMIT_ID}" && sha256 == "{ARCHIVE_TAR_GZ_SHA256}" {
		sha256 = mustGetSha256(strings.Replace(githubURL, "{ref}", ref, 1))
	}

	languages := []*Language{
		makeAndroid(),
		makeClosure(),
		makeCpp(),
		makeCsharp(),
		makeD(),
		makeDart(),
		makeGo(),
		makeJava(),
		makeNode(),
		makeObjc(),
		makePhp(),
		makePython(),
		makeRuby(),
		makeRust(),
		makeScala(),
		makeSwift(),

		makeGogo(),
		makeGrpcGateway(),
		makeGrpcJs(),
		makeGithubComGrpcGrpcWeb(),
	}

	for _, lang := range languages {
		mustWriteLanguageReadme(dir, lang)
		mustWriteLanguageRules(dir, lang)
		mustWriteLanguageExamples(dir, lang)
	}

	bazelVersions := []string{
		"BAZEL=0.24.1",
	}

	mustWriteReadme(dir, c.String("header"), c.String("footer"), struct {
		Ref, Sha256 string
	}{
		Ref:    ref,
		Sha256: sha256,
	}, languages, bazelVersions)

	mustWriteTravisYml(dir, c.String("travis_header"), c.String("travis_footer"), struct {
		Ref, Sha256 string
	}{
		Ref:    ref,
		Sha256: sha256,
	}, languages, bazelVersions)

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
            str(Label("{{ . }}")),{{ end }}
        ],
        **kwargs
    )`)

var usageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()`)

var grpcUsageTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:deps.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()`)

var protoCompileExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

var grpcCompileExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "greeter_{{ .Lang.Name }}_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)`)

var protoLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "person_{{ .Lang.Name }}_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)`)

var grpcLibraryExampleTemplate = mustTemplate(`load("@build_stack_rules_proto//{{ .Lang.Dir }}:{{ .Rule.Name }}.bzl", "{{ .Rule.Name }}")

{{ .Rule.Name }}(
    name = "greeter_{{ .Lang.Name }}_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)`)

func mustWriteLanguageRules(dir string, lang *Language) {
	for _, rule := range lang.Rules {
		mustWriteLanguageRule(dir, lang, rule)
	}
}

func mustWriteLanguageRule(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	out.t(rule.Implementation, &ruleData{lang, rule})
	out.ln()
	out.MustWrite(path.Join(dir, lang.Dir, rule.Name+".bzl"))
}

func mustWriteLanguageExamples(dir string, lang *Language) {
	for _, rule := range lang.Rules {
		exampleDir := path.Join(dir, lang.Dir, "example", rule.Name)
		os.MkdirAll(exampleDir, os.ModePerm)
		mustWriteLanguageExampleWorkspace(exampleDir, lang, rule)
		mustWriteLanguageExampleBuildFile(exampleDir, lang, rule)
		mustWriteLanguageExampleBazelrcFile(exampleDir, lang, rule)
	}
}

func mustWriteLanguageExampleWorkspace(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	depth := strings.Split(lang.Dir, "/")
	// +2 as we are in the example/{rule} subdirectory
	relpath := strings.Repeat("../", len(depth)+2)

	out.w(`load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

local_repository(
    name = "build_stack_rules_proto",
    path = "%s",
)`, relpath)

	out.ln()
	out.t(rule.Usage, &ruleData{lang, rule})
	out.ln()
	out.MustWrite(path.Join(dir, "WORKSPACE"))
}

func mustWriteLanguageExampleBuildFile(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	out.t(rule.Example, &ruleData{lang, rule})
	out.ln()
	out.MustWrite(path.Join(dir, "BUILD.bazel"))
}

func mustWriteLanguageExampleBazelrcFile(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	for _, f := range rule.Flags {
		out.w("# %s", f.Description)
		out.w("%s --%s=%s", f.Category, f.Name, f.Value)
	}
	out.MustWrite(path.Join(dir, ".bazelrc"))
}

func mustWriteLanguageReadme(dir string, lang *Language) {
	out := &LineWriter{}

	out.w("# `%s`", lang.Name)
	out.ln()

	if lang.Notes != nil {
		out.t(lang.Notes, lang)
		out.ln()
	}

	out.w("| Rule | Description |")
	out.w("| ---: | :--- |")
	for _, rule := range lang.Rules {
		out.w("| [%s](#%s) | %s |", rule.Name, rule.Name, rule.Doc)
	}
	out.ln()

	for _, rule := range lang.Rules {
		out.w(`---`)
		out.ln()
		out.w("## `%s`", rule.Name)
		out.ln()

		if rule.Experimental {
			out.w(`> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!`)
			out.ln()
		}
		out.w(rule.Doc)
		out.ln()

		out.w("### `WORKSPACE`")
		out.ln()

		out.w("```python")
		out.t(rule.Usage, &ruleData{lang, rule})
		out.w("```")
		out.ln()

		out.w("### `BUILD.bazel`")
		out.ln()

		out.w("```python")
		out.t(rule.Example, &ruleData{lang, rule})
		out.w("```")
		out.ln()

		if len(rule.Flags) > 0 {
			out.w("### `Flags`")
			out.ln()

			out.w("| Category | Flag | Value | Description |")
			out.w("| --- | --- | --- | --- |")
			for _, f := range rule.Flags {
				out.w("| %s | %s | %s | %s |", f.Category, f.Name, f.Value, f.Description)
			}
			out.ln()
		}

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

func mustWriteReadme(dir, header, footer string, data interface{}, languages []*Language, versions []string) {
	out := &LineWriter{}

	headVersion := versions[0]

	out.tpl(header, data)
	out.ln()

	out.w("## Rules")
	out.ln()

	out.w("| Status | Lang | Rule | Description")
	out.w("| ---    | ---: | :--- | :--- |")
	for _, lang := range languages {
		travisExclusionReason := lang.TravisExclusionReason
		for _, rule := range lang.Rules {
			travisLink := fmt.Sprintf("[![%s](https://travis-ci.org/stackb/rules_proto.svg?branch=master)](https://travis-ci.org/stackb/rules_proto)", headVersion)
			if travisExclusionReason == "" {
				travisExclusionReason = rule.TravisExclusionReason
			}
			if travisExclusionReason != "" {
				travisLink = travisExclusionReason
			}
			dirLink := fmt.Sprintf("[%s](/%s)", lang.Name, lang.Dir)
			ruleLink := fmt.Sprintf("[%s](/%s#%s)", rule.Name, lang.Dir, rule.Name)
			exampleLink := fmt.Sprintf("[example](/%s/example/%s)", lang.Name, rule.Name)
			out.w("| %s | %s | %s | %s (%s) |", travisLink, dirLink, ruleLink, rule.Doc, exampleLink)
		}
	}
	out.ln()

	out.tpl(footer, data)
	out.ln()

	out.MustWrite(path.Join(dir, "README.md"))
}

func mustWriteTravisYml(dir, header, footer string, data interface{}, languages []*Language, envVars []string) {
	out := &LineWriter{}

	out.tpl(header, data)

	for _, lang := range languages {
		if lang.TravisExclusionReason != "" {
			continue
		}
		for _, rule := range lang.Rules {
			if rule.TravisExclusionReason != "" {
				continue
			}
			env := make([]string, 0)
			for _, v := range envVars {
				env = append(env, v)
			}
			for _, v := range lang.TravisEnvVars {
				env = append(env, v)
			}
			for _, v := range rule.TravisEnvVars {
				env = append(env, v)
			}
			env = append(env, "LANG="+lang.Dir)
			env = append(env, "RULE="+rule.Name)

			out.w("  - %s", strings.Join(env, " "))
		}
	}
	out.ln()

	out.tpl(footer, data)
	out.ln()

	out.MustWrite(path.Join(dir, ".travis.yml"))
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

func (w *LineWriter) tpl(filename string, data interface{}) {
	tpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", filename, err)
	}
	w.t(tpl, data)
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

func mustGetSha256(url string) string {
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	h := sha256.New()
	if _, err := io.Copy(h, response.Body); err != nil {
		log.Fatal(err)
	}

	sha256 := fmt.Sprintf("%x", h.Sum(nil))

	log.Printf("sha256 for %s is %q", url, sha256)

	return sha256
}
