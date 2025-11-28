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

// Package symbol generates a `starlark_bundle` target for every `.bzl` file in
// each package.  At the root of the module, a single starlark_bundle is
// populated with deps that include all other symbol_libraries.
//
// The original code for this gazelle extension started from
// https://github.com/bazelbuild/bazel-skylib/blob/main/gazelle/bzl/gazelle.go.
package starlarkbundle

import (
	"log"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	starlarkBundleKind = "starlark_bundle"
	starlarkBundleName = "bundle"
)

var starlarkBundleKindInfo = map[string]rule.KindInfo{
	starlarkBundleKind: {
		NonEmptyAttrs: map[string]bool{"deps": true},
		ResolveAttrs:  map[string]bool{"deps": true},
	},
}

var starlarkBundleLoadInfo = rule.LoadInfo{
	Name:    "@build_stack_rules_proto//rules:starlark_bundle.bzl",
	Symbols: []string{starlarkBundleKind},
}

func starlarkBundleRule() (*rule.Rule, []string) {
	r := rule.NewRule(starlarkBundleKind, starlarkBundleName)

	r.SetAttr("visibility", []string{"//visibility:public"})

	return r, []string{}
}

func starlarkBundleImports(_ *config.Config, _ *rule.Rule, _ *rule.File) []resolve.ImportSpec {
	return nil
}

// starlarkBundleResolveResolve iterates the set of stalark_library rules and
// adds them to it's deps list.
func starlarkBundleResolve(r *rule.Rule, starlarkLibraries map[label.Label]*rule.Rule) {
	deps := make([]string, 0, len(starlarkLibraries))
	for dep := range starlarkLibraries {
		deps = append(deps, dep.String())
	}

	if len(deps) > 0 {
		sort.Strings(deps)
		r.SetAttr("deps", deps)
	}
}

// starlarkBundleResolveCoarse iterates the set of stalark_library rules and
// adds them to it's deps list.  It also takes out all external dependencies
// from the library rules and collects them into itself, but only as the :bundle
// target.  This is because we cannot ensure that the fine-grained dependency
// structure is correct. However, we can (hopefully) ensure that the
// coarse-grained dependency graph of bundles is correct.
func starlarkBundleResolveCoarse(r *rule.Rule, starlarkLibraries map[label.Label]*rule.Rule, moduleDeps map[string]string) {
	deps := make([]string, 0, len(starlarkLibraries))
	for dep, r := range starlarkLibraries {
		deps = append(deps, dep.String())

		ruleDeps := r.AttrStrings("deps")
		p1, _ := partitionExternalDeps(ruleDeps)
		r.SetAttr("deps", p1)
	}

	if false {
		for _, module := range moduleDeps {
			dep := label.New(module, "", "bundle")
			deps = append(deps, dep.String())
		}
	}

	if len(deps) > 0 {
		sort.Strings(deps)
		r.SetAttr("deps", deps)
	}
}

func partitionExternalDeps(deps []string) ([]string, []label.Label) {
	var p1 []string      // first and third-party deps
	var p3 []label.Label // first and third-party deps

	for _, dep := range deps {
		from, err := label.Parse(dep)
		if err != nil {
			log.Panicf("invalid dep label: %v", dep)
		}
		if from.Repo != "" {
			p3 = append(p3, from)
		} else {
			p1 = append(p1, dep)
		}
	}

	return p1, p3
}

func starlarkBundleGenerate(args language.GenerateArgs, root *string) (result language.GenerateResult) {
	if root == nil || *root != args.Rel {
		return result
	}

	r, loads := starlarkBundleRule()

	return language.GenerateResult{
		Gen:     []*rule.Rule{r},
		Imports: []any{loads},
	}
}
