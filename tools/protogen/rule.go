package protogen

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
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
		Rule: rule,
		Deps: collectDeps(rule),
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

func collectDeps(rule *ProtoRule) (allDeps []*ProtoDependency) {
	all := make(map[string]*ProtoDependency)

	for _, dep := range rule.Deps {
		all[dep.Name] = dep
	}

	for _, plugin := range rule.Plugins {
		for _, dep := range plugin.Deps {
			all[dep.Name] = dep
		}
	}

	for _, dep := range all {
		allDeps = append(allDeps, dep)
	}

	return
}

// ruleTemplateData is the type used by rules templates
type ruleTemplateData struct {
	Rule *ProtoRule
	Deps []*ProtoDependency
}
