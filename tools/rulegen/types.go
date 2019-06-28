package main

import (
	"text/template"
)


// Language represents one directory in this repo
type Language struct {
	// Directory in the repo where this language is rooted.  Typically this is
	// the same as the name
	Dir string

	// Name of the language
	Name string

	// Workspace usage
	WorkspaceExample string

	// List of rules
	Rules []*Rule

	// Additional nodes about the language
	Notes *template.Template

	// List of available plugins
	Plugins map[string]*Plugin

	// Bazel build flags required / suggested
	Flags []*Flag

	// Additional CI-specific env vars in the form "K=V"
	PresubmitEnvVars map[string]string
}


type Rule struct {
	// Name of the rule
	Name string

	// Base name of the rule (typically the lang name)
	Base string

	// Kind of the rule (proto|grpc)
	Kind string

	// Description
	Doc string

	// Temmplate for workspace
	WorkspaceExample *template.Template

	// Template for build file
	BuildExample *template.Template

	// Template for bzl file
	Implementation *template.Template

	// List of attributes
	Attrs []*Attr

	// List of plugins
	Plugins []string

	// Not expected to be functional
	Experimental bool

	// Not compatible with remote execution
	RemoteIncompatible bool

	// Bazel build flags required / suggested
	Flags []*Flag

	// Additional CI-specific env vars in the form "K=V"
	PresubmitEnvVars map[string]string

	// If not the empty string, one-word reason why excluded from bazelci
	// configuration
	BazelCIExclusionReason string
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


type Plugin struct {
	Tool    string
	Options []string
}


type ruleData struct {
	Lang *Language
	Rule *Rule
}
