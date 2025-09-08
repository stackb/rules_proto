package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
)

type LineWriter struct {
	lines []string
}

func (w *LineWriter) w(s string, args ...interface{}) {
	w.lines = append(w.lines, fmt.Sprintf(s, args...))
}

func (w *LineWriter) t(t *template.Template, data interface{}) {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		log.Fatalf("%v", err)
	}
	w.lines = append(w.lines, buf.String())
}

func (w *LineWriter) tpl(filename string, data interface{}) {
	tpl, err := template.ParseFiles(filename)
	if err != nil {
		log.Fatalf("Failed to parse %s: %v", filename, err)
	}
	w.t(tpl, data)
}

func (w *LineWriter) ln() {
	w.lines = append(w.lines, "")
}

func (w *LineWriter) MustWrite(filepath string) {
	err := os.WriteFile(filepath, []byte(strings.Join(w.lines, "\n")), 0666)
	if err != nil {
		log.Fatalf("FAIL %s: %v", filepath, err)
	}
	log.Printf("Wrote %s", filepath)
}

func (w *LineWriter) Write(filepath string) error {
	if err := os.WriteFile(filepath, []byte(strings.Join(w.lines, "\n")), 0666); err != nil {
		return fmt.Errorf("could not write %s: %w", filepath, err)
	}
	log.Printf("Wrote %s", filepath)
	return nil
}

func mustTemplate(tpl string) *template.Template {
	return template.Must(template.New("").Option("missingkey=error").Parse(tpl))
}
