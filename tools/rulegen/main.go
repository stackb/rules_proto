// Wsifier, a tool to parse BUILD files and bzl files, generate tests cases and
// documentation.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/urfave/cli"
)


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

	mustWriteReadme(dir, c.String("header"), c.String("footer"), struct {
		Ref, Sha256 string
	}{
		Ref:    ref,
		Sha256: sha256,
	}, languages)

	mustWriteBazelciPresubmitYml(dir, struct {
		Ref, Sha256 string
	}{
		Ref:    ref,
		Sha256: sha256,
	}, languages, []string{})

	mustWriteExamplesMakefile(dir, languages)
	mustWriteTestWorkspacesMakefile(dir)

	return nil
}


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
		exampleDir := path.Join(dir, "example", lang.Dir, rule.Name)
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
	out.t(rule.WorkspaceExample, &ruleData{lang, rule})
	out.ln()
	out.MustWrite(path.Join(dir, "WORKSPACE"))
}


func mustWriteLanguageExampleBuildFile(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	out.t(rule.BuildExample, &ruleData{lang, rule})
	out.ln()
	out.MustWrite(path.Join(dir, "BUILD.bazel"))
}


func mustWriteLanguageExampleBazelrcFile(dir string, lang *Language, rule *Rule) {
	out := &LineWriter{}
	for _, f := range lang.Flags {
		if f.Description != "" {
			out.w("# %s", f.Description)
		} else {
			out.w("#")
		}
		out.w("%s --%s=%s", f.Category, f.Name, f.Value)
	}
	for _, f := range rule.Flags {
		if f.Description != "" {
			out.w("# %s", f.Description)
		} else {
			out.w("#")
		}
		out.w("%s --%s=%s", f.Category, f.Name, f.Value)
	}
	out.ln()
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
		out.t(rule.WorkspaceExample, &ruleData{lang, rule})
		out.w("```")
		out.ln()

		out.w("### `BUILD.bazel`")
		out.ln()

		out.w("```python")
		out.t(rule.BuildExample, &ruleData{lang, rule})
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

	out.MustWrite(path.Join(dir, lang.Dir, "README.md"))
}


func mustWriteReadme(dir, header, footer string, data interface{}, languages []*Language) {
	out := &LineWriter{}

	badgeImageURL := "https://badge.buildkite.com/5980cc1d55f96e721bd9a7bd5dc1e40a096a7c30bc13117910.svg?branch=master"

	out.tpl(header, data)
	out.ln()

	out.w("## Rules")
	out.ln()

	out.w("| Status | Lang | Rule | Description")
	out.w("| ---    | ---: | :--- | :--- |")
	for _, lang := range languages {
		for _, rule := range lang.Rules {
			ciLink := fmt.Sprintf("[![Build Status](%s)](https://buildkite.com/bazel/rules-proto)", badgeImageURL)
			if rule.BazelCIExclusionReason != "" {
				ciLink = rule.BazelCIExclusionReason
			}
			dirLink := fmt.Sprintf("[%s](/%s)", lang.Name, lang.Dir)
			ruleLink := fmt.Sprintf("[%s](/%s#%s)", rule.Name, lang.Dir, rule.Name)
			exampleLink := fmt.Sprintf("[example](/example/%s/%s)", lang.Dir, rule.Name)
			out.w("| %s | %s | %s | %s (%s) |", ciLink, dirLink, ruleLink, rule.Doc, exampleLink)
		}
	}
	out.ln()

	out.tpl(footer, data)

	out.MustWrite(path.Join(dir, "README.md"))
}


func mustWriteBazelciPresubmitYml(dir string, data interface{}, languages []*Language, envVars []string) {
	out := &LineWriter{}
	platforms := []string{"ubuntu1604", "ubuntu1804", "windows", "macos"}

	// Write header
	out.w("---")
	out.w("tasks:")

	//
	// Write tasks for main code
	//
	for _, platform := range platforms {
		out.w("  main_%s:", platform)
		out.w("    name: build & test all")
		out.w("    platform: %s", platform)
		out.w("    environment:")
		out.w("      CC: clang")
		out.w("    test_flags:")
		out.w(`    - "--test_output=errors"`)
		out.w("    test_targets:")
		out.w(`    - "//example/routeguide/..."`)
		out.w("    build_targets:")
		for _, lang := range languages {
			// Skip experimental or excluded
			if lang.Name == "php" || lang.Name == "swift" || lang.Name == "csharp" || stringInSlice(platform, lang.BazelCIExcludePlatforms) {
				continue
			}
			out.w(`    - "//%s/..."`, lang.Dir)
		}
	}

	//
	// Write tasks for examples
	//
	for _, lang := range languages {
		for _, rule := range lang.Rules {
			if rule.BazelCIExclusionReason != "" {
				continue
			}

			exampleDir := path.Join(dir, "example", lang.Dir, rule.Name)

			for _, platform := range platforms {
				if stringInSlice(platform, rule.BazelCIExcludePlatforms) || stringInSlice(platform, lang.BazelCIExcludePlatforms) {
					continue
				}

				out.w("  %s_%s_%s:", lang.Name, rule.Name, platform)
				out.w("    name: '%s: %s'", lang.Name, rule.Name)
				out.w("    platform: %s", platform)
				out.w("    build_targets:")
				out.w(`      - "//..."`)
				out.w("    working_directory: %s", exampleDir)

				if len(lang.PresubmitEnvVars) > 0 || len(rule.PresubmitEnvVars) > 0 {
					out.w("    environment:")
					for k, v := range lang.PresubmitEnvVars {
						out.w("      %s: %s", k, v)
					}
					for k, v := range rule.PresubmitEnvVars {
						out.w("      %s: %s", k, v)
					}
				}
			}
		}
	}

	// Add test workspaces
	for _, testWorkspace := range findTestWorkspaceNames(dir) {
		for _, platform := range platforms {
			if platform == "windows" && (testWorkspace == "python2_grpc" || testWorkspace == "python3_grpc" || testWorkspace == "python_deps") {
				continue // Don't run python grpc test workspaces on windows
			}
			out.w("  test_workspace_%s_%s:", testWorkspace, platform)
			out.w("    name: 'test workspace: %s'", testWorkspace)
			out.w("    platform: %s", platform)
			out.w("    test_flags:")
			out.w(`    - "--test_output=errors"`)
			out.w("    test_targets:")
			out.w(`      - "//..."`)
			out.w("    working_directory: %s", path.Join(dir, "test_workspaces", testWorkspace))
		}
	}

	out.ln()
	out.MustWrite(path.Join(dir, ".bazelci", "presubmit.yml"))
}


func mustWriteExamplesMakefile(dir string, languages []*Language) {
	out := &LineWriter{}
	slashRegex := regexp.MustCompile("/")

	var allNames []string
	for _, lang := range languages {
		var langNames []string

		// Calculate depth of lang dir
		langDepth := len(slashRegex.FindAllStringIndex(lang.Dir, -1))

		// Create rules for each example
		for _, rule := range lang.Rules {
			exampleDir := path.Join(dir, "example", lang.Dir, rule.Name)

			var name = fmt.Sprintf("%s_%s_example", lang.Name, rule.Name)
			allNames = append(allNames, name)
			langNames = append(langNames, name)
			out.w("%s:", name)
			out.w("	cd %s; \\", exampleDir)
			out.w("	bazel build --disk_cache=%s../../bazel-disk-cache //... ; \\", strings.Repeat("../", langDepth))
			out.w("	bazel shutdown")
			out.ln()
		}

		// Create grouped rules for each language
		out.w("%s_examples: %s", lang.Name, strings.Join(langNames, " "))
		out.ln()
	}

	// Write all examples rule
	out.w("all_examples: %s", strings.Join(allNames, " "))

	out.ln()
	out.MustWrite(path.Join(dir, "example", "Makefile.mk"))
}


func mustWriteTestWorkspacesMakefile(dir string) {
	out := &LineWriter{}

	// For each test workspace, add makefile rule
	var allNames []string
	for _, testWorkspace := range findTestWorkspaceNames(dir) {
		var name = fmt.Sprintf("test_workspace_%s", testWorkspace)
		allNames = append(allNames, name)
		out.w("%s:", name)
		out.w("	cd %s; \\", path.Join(dir, "test_workspaces", testWorkspace))
		out.w("	bazel test --disk_cache=../bazel-disk-cache --test_output=errors //... ; \\")
		out.w("	bazel shutdown")
		out.ln()
	}

	// Write all test workspaces rule
	out.w("all_test_workspaces: %s", strings.Join(allNames, " "))

	out.ln()
	out.MustWrite(path.Join(dir, "test_workspaces", "Makefile.mk"))
}


func findTestWorkspaceNames(dir string) []string {
	files, err := ioutil.ReadDir(path.Join(dir, "test_workspaces"))
	if err != nil {
		log.Fatal(err)
	}

	var testWorkspaces []string
	for _, file := range files {
		if file.IsDir() && !strings.HasPrefix(file.Name(), ".") && !strings.HasPrefix(file.Name(), "bazel-") {
			testWorkspaces = append(testWorkspaces, file.Name())
		}
	}

	return testWorkspaces
}
