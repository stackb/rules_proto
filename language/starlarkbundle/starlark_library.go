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
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/pathtools"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/buildtools/build"
)

const (
	starlarkLibraryKind = "starlark_library"
	fileType            = ".bzl"
)

var ignoreSuffix = suffixes{
	"_tests.bzl",
	"_test.bzl",
}

var starlarkLibraryKindInfo = map[string]rule.KindInfo{
	starlarkLibraryKind: {
		NonEmptyAttrs:  map[string]bool{"srcs": true, "deps": true},
		MergeableAttrs: map[string]bool{"srcs": true},
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

func starlarkLibraryRule(args language.GenerateArgs, f string) (*rule.Rule, []string) {
	name := strings.TrimSuffix(f, fileType)
	r := rule.NewRule(starlarkLibraryKind, name)

	r.SetAttr("srcs", []string{f})

	shouldSetVisibility := args.File == nil || !args.File.HasDefaultVisibility()
	if shouldSetVisibility {
		vis := checkInternalVisibility(args.Rel, "//visibility:public")
		r.SetAttr("visibility", []string{vis})
	}

	fullPath := filepath.Join(args.Dir, f)
	loads, err := getBzlFileLoads(fullPath)

	if err != nil {
		log.Printf("%s: contains syntax errors: %v", fullPath, err)
		// don't return early since it is reasonable to create a target even
		// without deps.
	}

	return r, loads
}

func starlarkLibraryImports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	srcs := r.AttrStrings("srcs")
	imports := make([]resolve.ImportSpec, 0, len(srcs))

	for _, src := range srcs {
		spec := resolve.ImportSpec{
			// Lang is the language in which the import string appears (this should
			// match Resolver.Name).
			Lang: languageName,
			// Imp is an import string for the library.
			Imp: fmt.Sprintf("//%s:%s", f.Pkg, src),
		}

		imports = append(imports, spec)
	}

	return imports
}

func starlarkLibraryResolve(c *config.Config, ix *resolve.RuleIndex, rc *repo.RemoteCache, r *rule.Rule, importsRaw interface{}, from label.Label) {
	imports := importsRaw.([]string)

	r.DelAttr("deps")

	if len(imports) == 0 {
		return
	}

	deps := make([]string, 0, len(imports))
	for _, imp := range imports {
		impLabel, err := label.Parse(imp)
		if err != nil {
			log.Printf("%s: import of %q is invalid: %v", from.String(), imp, err)
			continue
		}

		// the index only contains absolute labels, not relative
		impLabel = impLabel.Abs(from.Repo, from.Pkg)

		if impLabel.Repo == "bazel_tools" {
			// The @bazel_tools repo is tricky because it is a part of the
			// "shipped with bazel" core library for interacting with the
			// outside world. This means that it can not depend on skylib.
			// Fortunately there is a fairly simple workaround for this, which
			// is that you can add those bzl files as `deps` entries.
			deps = append(deps, imp)
			continue
		}

		if impLabel.Repo != "" || !c.IndexLibraries {
			// This is a dependency that is external to the current repo, or
			// indexing is disabled so take a guess at what the target name
			// should be.
			deps = append(deps, strings.TrimSuffix(imp, fileType))
			continue
		}

		res := resolve.ImportSpec{
			Lang: languageName,
			Imp:  impLabel.String(),
		}
		matches := ix.FindRulesByImportWithConfig(c, res, languageName)
		if len(matches) == 0 {
			log.Printf("%s: %q (%s) was not found in dependency index. Skipping. This may result in an incomplete deps section and require manual BUILD file intervention.\n", from.String(), imp, impLabel.String())
		}

		for _, m := range matches {
			depLabel := m.Label
			depLabel = depLabel.Rel(from.Repo, from.Pkg)
			deps = append(deps, depLabel.String())
		}
	}

	sort.Strings(deps)
	if len(deps) > 0 {
		r.SetAttr("deps", deps)
	}
}

func starlarkLibararyGenerate(args language.GenerateArgs) language.GenerateResult {
	var rules []*rule.Rule
	var imports []any

	for _, f := range append(args.RegularFiles, args.GenFiles...) {
		if !isBzlSourceFile(f) {
			continue
		}
		r, loads := starlarkLibraryRule(args, f)
		rules = append(rules, r)
		imports = append(imports, loads)
	}

	return language.GenerateResult{
		Gen:     rules,
		Imports: imports,
		Empty:   generateEmpty(args),
	}
}

func getBzlFileLoads(path string) ([]string, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile(%q) error: %v", path, err)
	}
	ast, err := build.ParseBuild(path, f)
	if err != nil {
		return nil, fmt.Errorf("build.Parse(%q) error: %v", f, err)
	}

	var loads []string
	build.WalkOnce(ast, func(expr *build.Expr) {
		n := *expr
		if l, ok := n.(*build.LoadStmt); ok {
			loads = append(loads, l.Module.Value)
		}
	})
	sort.Strings(loads)

	return loads, nil
}

func isBzlSourceFile(f string) bool {
	return strings.HasSuffix(f, fileType) && !ignoreSuffix.Matches(f)
}

// generateEmpty generates the list of rules that don't need to exist in the
// BUILD file any more. For each symbol_library rule in args.File that only has
// srcs that aren't in args.RegularFiles or args.GenFiles, add a symbol_library
// with no srcs or deps. That will let Gazelle delete symbol_library rules after
// the corresponding .bzl files are deleted.
func generateEmpty(args language.GenerateArgs) []*rule.Rule {
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
			srcsExist := false
			for _, f := range r.AttrStrings("srcs") {
				if exists[f] {
					srcsExist = true
					break
				}
			}
			if !srcsExist {
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
