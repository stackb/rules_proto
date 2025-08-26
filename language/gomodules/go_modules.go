package gomodules

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type goModules struct {
	loadName  string
	targetPkg string
}

func newGoModules(loadname, targetPkg string) *goModules {
	return &goModules{
		loadName:  loadname,
		targetPkg: targetPkg,
	}
}

func (m *goModules) kind() string {
	return "go_modules"
}

func (m *goModules) loadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    m.loadName,
		Symbols: []string{m.kind()},
	}
}

func (m *goModules) kindInfo() rule.KindInfo {
	return rule.KindInfo{
		MatchAny:      true,
		ResolveAttrs:  map[string]bool{"deps": true},
		NonEmptyAttrs: map[string]bool{"deps": true},
	}
}

func (m *goModules) generate(fromPkg string) (*rule.Rule, bool) {
	if fromPkg != m.targetPkg {
		return nil, false
	}
	goModules := rule.NewRule(m.kind(), m.kind())
	return goModules, true
}

func (m *goModules) resolve(_ label.Label, r *rule.Rule, index map[label.Label]bool) bool {
	if r.Kind() != m.kind() {
		return false
	}

	deps := make([]string, len(index))
	i := 0
	for lbl := range index {
		deps[i] = lbl.String()
		i++
	}
	sort.Strings(deps)

	r.SetAttr("deps", deps)

	return false
}
