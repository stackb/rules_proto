package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

func generate(rule *Rule) error {
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

	workspaceDir, ok := os.LookupEnv("BUILD_WORKSPACE_DIRECTORY")
	if !ok {
		return fmt.Errorf("BUILD_WORKSPACE_DIRECTORY not set")
	}
	if err := os.Chdir(workspaceDir); err != nil {
		return err
	}

	if err := copyToWorkspace(rule, workspaceDir, rule.ImplementationFilename); err != nil {
		return err
	}

	return nil
}

func generateFile(rule *Rule, tmpl *template.Template, filename string) error {
	out := &LineWriter{}
	out.t(tmpl, &templateData{rule})
	out.ln()
	if err := out.Write(filename); err != nil {
		return fmt.Errorf("could not output %s: %w", filename, err)
	}
	return nil
}

func copyToWorkspace(rule *Rule, workspaceDir, src string) error {
	dst := filepath.Join(workspaceDir, rule.Package, filepath.Base(src))
	if _, err := copy(src, dst); err != nil {
		return fmt.Errorf("could not copy %s -> %s: %w", src, dst, err)
	}
	return nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
