/* Copyright 2020 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package symbol generates a `starlark_library` target for every `.bzl` file in
// each package.  At the root of the module, a single starlark_bundle is
// populated with deps that include all other symbol_libraries.
//
// The original code for this gazelle extension started from
// https://github.com/bazelbuild/bazel-skylib/blob/main/gazelle/bzl/gazelle.go.
package starlarkbundle

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/pathtools"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/buildtools/build"
)

const (
	starlarkLibraryKind       = "starlark_library"
	starlarkLibraryNamePrefix = "lib"
	fileType                  = ".bzl"
	visibilityPublic          = "//visibility:public"
)

var ignoreSuffix = suffixes{
	// "_tests.bzl",
	// "_test.bzl",
}

var starlarkLibraryKindInfo = map[string]rule.KindInfo{
	starlarkLibraryKind: {
		NonEmptyAttrs: map[string]bool{"src": true},
		ResolveAttrs:  map[string]bool{"deps": true},
	},
}

var starlarkLibraryLoadInfo = rule.LoadInfo{
	Name:    "@build_stack_rules_proto//rules:starlark_library.bzl",
	Symbols: []string{starlarkLibraryKind},
}

type suffixes []string

func (s suffixes) Matches(test string) bool {
	for _, v := range s {
		if strings.HasSuffix(test, v) {
			return true
		}
	}
	return false
}

func starlarkLibraryRule(args language.GenerateArgs, f string, ext *starlarkBundleLang) (*rule.Rule, []*build.LoadStmt) {
	fullPath := filepath.Join(args.Dir, f)
	name := starlarkLibraryNamePrefix + strings.TrimSuffix(f, fileType)

	ast, loads, err := getBzlFileLoadsStmts(fullPath)
	if err != nil {
		ext.logf("%s: contains syntax errors: %v", fullPath, err)
		// don't return early since it is reasonable to create a target even
		// without deps.
	}

	// always castrate fail() exprs
	if renameFailToPrint(ast) {
		ext.emitBzlFile(fullPath, ast)
	}

	r := rule.NewRule(starlarkLibraryKind, name)

	r.SetAttr("src", f)

	shouldSetVisibility := args.File == nil || !args.File.HasDefaultVisibility()
	if shouldSetVisibility {
		vis := checkInternalVisibility(args.Rel, visibilityPublic)
		r.SetAttr("visibility", []string{vis})
	}

	r.SetPrivateAttr("full_path", fullPath)
	r.SetPrivateAttr("ast", ast)

	return r, loads
}

func starlarkLibraryImports(_ *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	src := r.AttrString("src")
	return []resolve.ImportSpec{{
		// Lang is the language in which the import string appears (this should
		// match Resolver.Name).
		Lang: languageName,
		// Imp is an import string for the library.
		Imp: fmt.Sprintf("//%s:%s", f.Pkg, src),
	}}
}

func starlarkLibraryResolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, importsRaw interface{}, from label.Label, ext *starlarkBundleLang) {
	loads := importsRaw.([]*build.LoadStmt)
	fullPath := r.PrivateAttr("full_path").(string)
	ast := r.PrivateAttr("ast").(*build.File)

	ext.logf("%s: starlarkLibraryResolve: %s with %d loads", fullPath, from.String(), len(loads))

	r.DelAttr("deps")

	if len(loads) == 0 {
		ext.logf("  no loads to process")
		return
	}

	rewriteBzlSourceFile := false
	var unknownDeps []string
	var bazelToolsDeps []string

	deps := make([]string, 0, len(loads))
	for i, load := range loads {
		imp := load.Module.Value
		ext.logf("  processing load %d/%d: %q", i+1, len(loads), imp)

		impLabel, err := label.Parse(imp)
		if err != nil {
			ext.logf("    ERROR: import of %q is invalid: %v", imp, err)
			continue
		}

		// the index only contains absolute labels, not relative
		impLabel = impLabel.Abs(from.Repo, from.Pkg)
		ext.logf("    absolute label: %s (repo=%q, pkg=%q, name=%q)", impLabel.String(), impLabel.Repo, impLabel.Pkg, impLabel.Name)

		if impLabel.Repo == "bazel_tools" {
			// The @bazel_tools repo is tricky because it is a part of the
			// "shipped with bazel" core library for interacting with the
			// outside world. This means that it can not depend on skylib.
			// Fortunately there is a fairly simple workaround for this, which
			// is that you can add those bzl files as `deps` entries.
			//
			// the bazel source code gathers them up in filegroups but not
			// exposed publically.  For this to work, caller must be able to use
			// a modified version of bazel that adds public visibility to the targets (see `tools/build_defs/repo/BUILD.repo`)
			bazelToolsLabel := label.New(impLabel.Repo, impLabel.Pkg, "bzl_srcs")
			ext.logf("    bazel_tools dependency: adding to bazel_tools_deps")
			// deps = append(deps, imp)
			bazelToolsDeps = append(bazelToolsDeps, bazelToolsLabel.String())
			// unknownDeps = append(deps, imp)
			continue
		}

		if impLabel.Repo != "" || !c.IndexLibraries {
			// This is a dependency that is external to the current repo.
			// Rewrite the repo label to one suffixed by "_docs".  We expect to
			// find the starlark_library dependency that provides the file in
			// that repo.  Rewrite the load label because starlark_doc_extract
			// will also expect to load the symbol from that location.
			ext.logf("    external repo dependency: repo=%q, IndexLibraries=%v", impLabel.Repo, c.IndexLibraries)
			extRepo, known := ext.getModuleDependencyRepoName(impLabel.Repo)
			extLabel := label.New(
				extRepo,
				impLabel.Pkg,
				starlarkLibraryNamePrefix+strings.TrimSuffix(impLabel.Name, fileType),
			)

			loadLabel := label.New(extLabel.Repo, extLabel.Pkg, impLabel.Name)
			ext.logf("    rewriting load: %q -> %q", load.Module.Value, loadLabel.String())
			load.Module.Value = loadLabel.String()
			rewriteBzlSourceFile = true

			if known {
				ext.logf("    adding to deps: %s", extLabel.String())
				deps = append(deps, extLabel.String())
			} else {
				ext.logf("    adding to unknownDeps: %s", extLabel.String())
				unknownDeps = append(unknownDeps, extLabel.String())
			}

			continue
		}

		ext.logf("    internal dependency: looking up in index")
		res := resolve.ImportSpec{
			Lang: languageName,
			Imp:  impLabel.String(),
		}
		matches := ix.FindRulesByImportWithConfig(c, res, languageName)
		ext.logf("    found %d matches in index", len(matches))
		if len(matches) == 0 {
			ext.logf("    WARNING: %q (%s) was not found in dependency index", imp, impLabel.String())
			// unknownDeps = append(unknownDeps, impLabel.String())
		}

		for _, m := range matches {
			depLabel := m.Label
			ext.logf("    adding match to deps: %s", depLabel.String())
			// depLabel = depLabel.Rel(from.Repo, from.Pkg)
			deps = append(deps, depLabel.String())
		}
	}

	ext.logf("  resolution complete: %d deps, %d unknown, %d bazel_tools", len(deps), len(unknownDeps), len(bazelToolsDeps))

	if len(deps) > 0 {
		deps = deduplicateAndSort(deps)
		r.SetAttr("deps", deps)
		ext.logf("  set deps attribute: %v", deps)
	}
	if len(unknownDeps) > 0 {
		unknownDeps = deduplicateAndSort(unknownDeps)
		r.SetAttr("unknown_deps", unknownDeps)
		ext.logf("  set unknown_deps attribute: %v", unknownDeps)
	}
	if len(bazelToolsDeps) > 0 {
		bazelToolsDeps = deduplicateAndSort(bazelToolsDeps)
		r.SetAttr("bazel_tools_deps", bazelToolsDeps)
		ext.logf("  set bazel_tools_deps attribute: %v", bazelToolsDeps)
	}

	if rewriteBzlSourceFile {
		ext.emitBzlFile(fullPath, ast)
	}
}

func (ext *starlarkBundleLang) emitBzlFile(fullPath string, ast *build.File) {
	ext.logf("emitting source file: %s", fullPath)
	data := build.Format(ast)
	if err := os.WriteFile(fullPath, data, os.ModePerm); err != nil {
		ext.logf("  ERROR: failed to emit rewritten bzl file: %v", err)
	}
}

func starlarkLibraryGenerate(args language.GenerateArgs, starlarkLibraries map[label.Label]*rule.Rule, ext *starlarkBundleLang) language.GenerateResult {
	var rules []*rule.Rule
	var imports []any

	for _, f := range append(args.RegularFiles, args.GenFiles...) {
		if !isBzlSourceFile(f) {
			continue
		}
		r, loads := starlarkLibraryRule(args, f, ext)
		if r == nil {
			continue
		}

		rules = append(rules, r)
		imports = append(imports, loads)

		// populate the map so the bundle rule can use them later.
		if isVisibilityPublic(r.AttrStrings("visibility")) {
			from := label.New(args.Config.RepoName, args.Rel, r.Name())
			starlarkLibraries[from] = r
		}
	}

	return language.GenerateResult{
		Gen:     rules,
		Imports: imports,
		Empty:   starlarkLibraryEmptyRules(args),
	}
}

func getBzlFileLoadsStmts(path string) (*build.File, []*build.LoadStmt, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("os.ReadFile(%q) error: %v", path, err)
	}
	ast, err := build.ParseBzl(path, f)
	if err != nil {
		return nil, nil, fmt.Errorf("build.Parse(%q) error: %v", f, err)
	}

	var loads []*build.LoadStmt
	build.WalkOnce(ast, func(expr *build.Expr) {
		n := *expr
		if l, ok := n.(*build.LoadStmt); ok {
			loads = append(loads, l)
		}
	})

	return ast, loads, nil
}

func renameFailToPrint(ast *build.File) bool {
	modified := false
	build.Walk(ast, func(expr build.Expr, stack []build.Expr) {
		if call, ok := expr.(*build.CallExpr); ok {
			if callName, ok := call.X.(*build.Ident); ok {
				if callName.Name == "fail" {
					callName.Name = "print"
					modified = true
				}
			}
		}
	})
	return modified
}

func isBzlSourceFile(f string) bool {
	return strings.HasSuffix(f, fileType) && !ignoreSuffix.Matches(f)
}

func isVisibilityPublic(vis []string) bool {
	return len(vis) == 1 && vis[0] == visibilityPublic
}

// starlarkLibraryEmptyRules generates the list of rules that don't need to
// exist in the BUILD file any more. For each symbol_library rule in args.File
// that only has srcs that aren't in args.RegularFiles or args.GenFiles, add a
// symbol_library with no srcs or deps. That will let Gazelle delete
// symbol_library rules after the corresponding .bzl files are deleted.
func starlarkLibraryEmptyRules(args language.GenerateArgs) []*rule.Rule {
	var ret []*rule.Rule
	if args.File == nil {
		return ret
	}
	for _, r := range args.File.Rules {
		if r.Kind() != starlarkLibraryKind {
			continue
		}
		name := r.AttrString("name")

		exists := make(map[string]bool)
		for _, f := range args.RegularFiles {
			exists[f] = true
		}
		for _, f := range args.GenFiles {
			exists[f] = true
		}
		for _, r := range args.File.Rules {
			srcExist := exists[r.AttrString("src")]
			if !srcExist {
				ret = append(ret, rule.NewRule(starlarkLibraryKind, name))
			}
		}
	}
	return ret
}

// checkInternalVisibility overrides the given visibility if the package is
// internal.
func checkInternalVisibility(rel, visibility string) string {
	if i := pathtools.Index(rel, "internal"); i > 0 {
		visibility = fmt.Sprintf("//%s:__subpackages__", rel[:i-1])
	} else if i := pathtools.Index(rel, "private"); i > 0 {
		visibility = fmt.Sprintf("//%s:__subpackages__", rel[:i-1])
	} else if pathtools.HasPrefix(rel, "internal") || pathtools.HasPrefix(rel, "private") {
		visibility = "//:__subpackages__"
	}
	return visibility
}

func sanitizeName(name string) string {
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ReplaceAll(name, "+", "_")
	return name
}

// deduplicateAndSort removes duplicate entries and sorts the list
func deduplicateAndSort(in []string) (out []string) {
	if len(in) == 0 {
		return in
	}
	seen := make(map[string]bool)
	for _, v := range in {
		if seen[v] {
			continue
		}
		seen[v] = true
		out = append(out, v)
	}
	sort.Strings(out)
	return
}
