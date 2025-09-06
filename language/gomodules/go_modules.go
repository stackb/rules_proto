package gomodules

import (
	"log"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type goModules struct {
	cfg  Config
	deps map[label.Label]bool
}

func newGoModules(cfg Config) *goModules {
	return &goModules{
		cfg:  cfg,
		deps: make(map[label.Label]bool),
	}
}

func (m *goModules) kind() string {
	return "go_modules"
}

func (m *goModules) loadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    m.cfg.LoadName(),
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

func (m *goModules) generateRule(args language.GenerateArgs) (*rule.Rule, bool) {
	log.Println("generateRule:", args.Rel)
	for _, r := range args.OtherGen {
		switch r.Kind() {
		case "proto_go_library":
			dep := label.New(args.Config.RepoName, args.Rel, r.Name())
			m.deps[dep] = true
			log.Println("indexed dep:", dep)
		}
	}

	if !(args.Rel == m.cfg.TargetPkg() || (args.Rel == "" && m.cfg.TargetPkg() == "ROOT")) {
		return nil, false
	}

	goModules := rule.NewRule(m.kind(), m.kind())
	goModules.SetAttr("visibility", []string{"//visibility:public"})

	return goModules, true
}

func (m *goModules) resolve(_ label.Label, r *rule.Rule) bool {
	if r.Kind() != m.kind() {
		return false
	}

	deps := make([]string, 0)
	for dep, ok := range m.deps {
		if !ok {
			continue
		}
		deps = append(deps, dep.String())
	}
	sort.Strings(deps)

	r.SetAttr("deps", deps)

	return false
}
