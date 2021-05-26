package main

import (
	"fmt"
	"os"
	"text/template"
)

func generate(c *Config) error {
	f, err := os.Create(c.Out)
	if err != nil {
		return fmt.Errorf("create %s: %v", c.Out, err)
	}
	defer f.Close()

	t := template.Must(template.New(c.Name + "_deps.bzl").Option("missingkey=error").Parse(depsBzl))

	data := &templateData{
		Name: c.Name,
		Deps: collectDeps(c.Deps),
	}

	if err := t.Execute(f, data); err != nil {
		return fmt.Errorf("template %s: %v", c.Out, err)
	}

	return nil
}
