package protorule

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"text/template"
)

// Generate takes the given rule definition and writes the result files to the
// filesystem.
func Generate(rule *ProtoRule) error {
	if err := generateFile(rule, rule.Implementation, rule.ImplementationFilename); err != nil {
		return err
	}
	if err := generateFile(rule, rule.WorkspaceExample, rule.WorkspaceExampleFilename); err != nil {
		return err
	}
	if err := generateFile(rule, rule.BuildExample, rule.BuildExampleFilename); err != nil {
		return err
	}
	if err := generateFile(rule, rule.Test, rule.TestFilename); err != nil {
		return err
	}
	return nil
}

func generateFile(rule *ProtoRule, tmpl *template.Template, filename string) error {
	out := &LineWriter{}
	out.t(tmpl, &templateData{rule})
	out.ln()
	if err := out.Write(filename); err != nil {
		return fmt.Errorf("could not output %s: %w", filename, err)
	}
	return nil
}

func FromJSONFile(filename string) (*ProtoRule, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}

	var rule ProtoRule
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}

	return &rule, nil
}

func ParseRuleTemplates(rule *ProtoRule) error {
	tpl, err := template.ParseFiles(rule.ImplementationTmpl)
	if err != nil {
		return fmt.Errorf("could not prepare rule %w", err)
	}
	rule.Implementation = tpl

	tpl, err = template.ParseFiles(rule.WorkspaceExampleTmpl)
	if err != nil {
		return fmt.Errorf("could not prepare rule %w", err)
	}
	rule.WorkspaceExample = tpl

	tpl, err = template.ParseFiles(rule.BuildExampleTmpl)
	if err != nil {
		return fmt.Errorf("could not prepare rule %w", err)
	}
	rule.BuildExample = tpl

	tpl, err = template.ParseFiles(rule.TestTmpl)
	if err != nil {
		return fmt.Errorf("could not prepare rule %w", err)
	}
	rule.Test = tpl

	return nil
}
