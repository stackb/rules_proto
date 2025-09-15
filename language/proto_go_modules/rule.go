package proto_go_modules

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type protoGoModules struct {
	cfg  Config
	deps map[label.Label]bool
}

func newProtoGoModules(cfg Config) *protoGoModules {
	return &protoGoModules{
		cfg:  cfg,
		deps: make(map[label.Label]bool),
	}
}

func (m *protoGoModules) kind() string {
	return "proto_go_modules"
}

func (m *protoGoModules) loadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    m.cfg.LoadName(),
		Symbols: []string{m.kind()},
	}
}

func (m *protoGoModules) kindInfo() rule.KindInfo {
	return rule.KindInfo{
		MatchAny:      true,
		ResolveAttrs:  map[string]bool{"deps": true},
		NonEmptyAttrs: map[string]bool{"deps": true},
	}
}

func (m *protoGoModules) generateRule(args language.GenerateArgs) (*rule.Rule, bool) {
	indexKinds := m.cfg.IndexKinds()

	for _, r := range args.OtherGen {
		if indexKinds[r.Kind()] {
			dep := label.New(args.Config.RepoName, args.Rel, r.Name())
			m.deps[dep] = true
		}
	}

	if args.Rel != "" {
		return nil, false
	}

	newRule := rule.NewRule(m.kind(), m.kind())
	newRule.SetAttr("visibility", []string{"//visibility:public"})

	return newRule, true
}

func (m *protoGoModules) resolve(_ label.Label, r *rule.Rule) bool {
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
