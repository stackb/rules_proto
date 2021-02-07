package protogen

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
)

// Generate takes the given rule definition and writes the result files to the
// filesystem.
func (lang *ProtoLanguage) Generate() error {
	data := &langTemplateData{lang}

	if err := generateFile(lang.Templates, lang.MarkdownTmpl, lang.MarkdownFilename, data); err != nil {
		return err
	}
	if err := generateFile(lang.Templates, lang.RulesTmpl, lang.RulesFilename, data); err != nil {
		return err
	}

	return nil
}

// NewProtoLanguageFromJSONFile constructs a ProtoLanguage struct from the given
// filename that contains a JSON.
func NewProtoLanguageFromJSONFile(filename string) (*ProtoLanguage, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not make rule %w", err)
	}

	var lang ProtoLanguage
	if err := json.Unmarshal(data, &lang); err != nil {
		return nil, fmt.Errorf("could not make lang %w", err)
	}

	lang.Templates = template.Must(template.ParseFiles(
		lang.MarkdownTmpl,
		lang.RulesTmpl,
	))

	return &lang, nil
}

// langTemplateData is the type used by language templates
type langTemplateData struct {
	Lang *ProtoLanguage
}
