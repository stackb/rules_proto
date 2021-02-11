package protogen

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"
)

// NewProtoRuleFromJSONFile constructs a ProtoRule struct from the given
// filename that contains a JSON.
func NewProtoRuleFromJSONFile(filename string) (*ProtoRule, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}

	var rule ProtoRule
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}

	rule.Templates = template.Must(template.ParseFiles(
		rule.ImplementationTmpl,
		rule.WorkspaceExampleTmpl,
		rule.BuildExampleTmpl,
		rule.MarkdownTmpl,
		rule.DepsTmpl,
		rule.BazelTestTmpl,
	))

	return &rule, nil
}

// Generate takes the given rule definition and writes the result files to the
// filesystem.
func (rule *ProtoRule) Generate() error {
	data := &ruleTemplateData{
		Rule:     rule,
		RuleDeps: collectRuleDeps(rule),
	}
	if err := generateFile(rule.Templates, rule.ImplementationTmpl, rule.ImplementationFilename, data); err != nil {
		return err
	}
	if err := generateFile(rule.Templates, rule.WorkspaceExampleTmpl, rule.WorkspaceExampleFilename, data); err != nil {
		return err
	}
	if err := generateFile(rule.Templates, rule.BuildExampleTmpl, rule.BuildExampleFilename, data); err != nil {
		return err
	}
	if err := generateFile(rule.Templates, rule.MarkdownTmpl, rule.MarkdownFilename, data); err != nil {
		return err
	}
	if err := generateFile(rule.Templates, rule.DepsTmpl, rule.DepsFilename, data); err != nil {
		return err
	}
	if err := generateFile(rule.Templates, rule.BazelTestTmpl, rule.BazelTestFilename, data); err != nil {
		return err
	}
	return nil
}

// collectRuleDeps accumulates the transitive dependencies of the given rule,
// eliminating duplicates but maintaining DFS ordering.
func collectRuleDeps(rule *ProtoRule) (deps []*ruleDependency) {
	seen := make(map[string]bool)

	var visit func(string, *ProtoDependency)
	visit = func(parentName string, dep *ProtoDependency) {
		if seen[dep.Name] {
			return
		}
		seen[dep.Name] = true
		deps = append(deps, &ruleDependency{
			ParentName: parentName,
			Dep:        dep,
		})
		for _, child := range dep.Deps {
			visit(dep.Name, child)
		}
	}

	for _, dep := range rule.Deps {
		if seen[dep.Name] {
			continue
		}
		visit("rule "+rule.Name, dep)
	}

	for _, plugin := range rule.Plugins {
		for _, dep := range plugin.Deps {
			if seen[dep.Name] {
				continue
			}
			visit("plugin "+plugin.Label, dep)
		}
	}

	reverseDeps(deps)

	return
}

func reverseDeps(deps []*ruleDependency) {
	for i, j := 0, len(deps)-1; i < j; i, j = i+1, j-1 {
		deps[i], deps[j] = deps[j], deps[i]
	}
}

// ruleTemplateData is the type used by rules templates
type ruleTemplateData struct {
	Rule     *ProtoRule
	RuleDeps []*ruleDependency
}

type ruleDependency struct {
	ParentName string
	Dep        *ProtoDependency
}
