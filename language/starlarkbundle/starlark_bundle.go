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
	"sort"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	starlarkBundleKind = "starlark_bundle"
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
	r := rule.NewRule(starlarkBundleKind, starlarkBundleKind)

	r.SetAttr("visibility", []string{"//visibility:public"})

	return r, []string{}
}

func starlarkBundleImports(_ *config.Config, _ *rule.Rule, _ *rule.File) []resolve.ImportSpec {
	return nil
}

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
