package protorule

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Generate takes the given rule definition and writes the result files to the
// filesystem.
func Generate(rule *ProtoRule) error {
	if err := generateFile(rule, rule.ImplementationTmpl, rule.ImplementationFilename); err != nil {
		return err
	}
	if err := generateFile(rule, rule.WorkspaceExampleTmpl, rule.WorkspaceExampleFilename); err != nil {
		return err
	}
	if err := generateFile(rule, rule.BuildExampleTmpl, rule.BuildExampleFilename); err != nil {
		return err
	}
	if err := generateFile(rule, rule.MarkdownTmpl, rule.MarkdownFilename); err != nil {
		return err
	}
	if err := generateFile(rule, rule.DepsTmpl, rule.DepsFilename); err != nil {
		return err
	}
	return nil
}

func generateFile(rule *ProtoRule, templateName, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("count not generate %s: %w", filename, err)
	}
	defer f.Close()

	if err := rule.Templates.ExecuteTemplate(f, filepath.Base(templateName), &templateData{rule}); err != nil {
		return fmt.Errorf("count not generate %s: %w", filename, err)
	}

	return nil
}

// FromJSONFile constructs a ProtoRule struct from the given filename that
// contains a JSON.
func FromJSONFile(filename string) (*ProtoRule, error) {
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
	))

	return &rule, nil
}
