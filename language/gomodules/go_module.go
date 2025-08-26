package gomodules

import (
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type goModule struct {
	index     map[label.Label]bool
	loadName  string
	targetPkg string
}

func newGoModule(loadName, targetPkg string) *goModule {
	return &goModule{
		index:     make(map[label.Label]bool),
		loadName:  loadName,
		targetPkg: targetPkg,
	}
}

func (m *goModule) kind() string {
	return "go_module"
}

func (m *goModule) loadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    m.loadName,
		Symbols: []string{m.kind()},
	}
}

func (m *goModule) kindInfo() rule.KindInfo {
	return rule.KindInfo{
		MatchAny: true,
		MergeableAttrs: map[string]bool{
			"srcs":       true,
			"importpath": true,
		},
		NonEmptyAttrs: map[string]bool{"srcs": true},
	}
}

func (m *goModule) generate(fromPkg string, args language.GenerateArgs) (*rule.Rule, bool) {
	var srcs []string
	var importpath string

	for _, r := range args.OtherGen {
		switch r.Kind() {
		case "proto_compile":
			if strings.HasSuffix(r.Name(), "_go_compile") {
				srcs = append(srcs, r.AttrStrings("outputs")...)
			}
		case "proto_go_library":
			importpath = r.AttrString("importpath")
		}
	}

	if len(srcs) == 0 {
		return nil, false
	}

	goModule := rule.NewRule(m.kind(), "go_module")
	goModule.SetAttr("importpath", importpath)
	goModule.SetAttr("visibility", []string{m.targetPkg + ":__pkg__"})
	goModule.SetAttr("srcs", srcs)

	m.index[label.New(args.Config.RepoName, fromPkg, goModule.Name())] = true

	return goModule, true
}
