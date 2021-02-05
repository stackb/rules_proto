package main

import (
	"text/template"
)

type Rule struct {
	// Name of the rule
	Name string

	// Package of the rule (label package)
	Package string

	// Kind of the rule (proto|grpc)
	Kind string

	// Description
	Doc string

	// ImplementationFilename of the rule (full path)
	ImplementationFilename string

	// Filename of the implementation template
	ImplementationTmpl string

	// Temmplate for implementation
	Implementation *template.Template

	// Temmplate for workspace
	WorkspaceExample *template.Template

	// Filename of the workspace template
	WorkspaceExampleTmpl string

	// WorkspaceExampleFilename of the rule (full path)
	WorkspaceExampleFilename string

	// BuildExampleFilename of the rule (full path)
	BuildExampleFilename string

	// Filename of the BuildExample template
	BuildExampleTmpl string

	// Temmplate for build example
	BuildExample *template.Template

	// TestFilename of the rule (full path)
	TestFilename string

	// Filename of the Test template
	TestTmpl string

	// Temmplate for Test
	Test *template.Template

	// List of attributes
	Attrs []*Attr

	// List of plugins
	Plugins []string

	// Not expected to be functional
	Experimental bool

	// Bazel build flags required / suggested
	Flags []*Flag

	// Additional CI-specific env vars in the form "K=V"
	PresubmitEnvVars map[string]string

	// Platforms for which to skip testing this rule, overrides language level
	// The special value 'all' will skip app platforms
	SkipTestPlatforms []string

	// Flag indicating if the merge_directories flag should be set to false for
	// the generated rule
	SkipDirectoriesMerge bool
}

// Flag captures information about a bazel build flag.
type Flag struct {
	Category    string
	Name        string
	Value       string
	Description string
}

type Attr struct {
	Name      string
	Type      string
	Default   string
	Doc       string
	Mandatory bool
}

// templateData is the type used by templates
type templateData struct {
	Rule *Rule
}
