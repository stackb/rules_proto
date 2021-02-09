package protogen

import (
	"fmt"
	"text/template"
	"os"
	"path/filepath"
)

func generateFile(templates *template.Template, templateName, filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("count not generate %s: %w", filename, err)
	}
	defer f.Close()

	if err := templates.ExecuteTemplate(f, filepath.Base(templateName), data); err != nil {
		return fmt.Errorf("count not generate %s: %w", filename, err)
	}

	return nil
}
