package laze

import (
	"github.com/bazelbuild/bazel-gazelle/language"
)

// plugin satisfies the language.Language interface. It is the Gazelle extension
// for plugin rules.
type plugin struct {
	Configurer
	Resolver
}

// NewLanguage initializes a new plugin that satisfies the language.Language
// interface. This is the entrypoint for the extension initialization.
func NewLanguage() language.Language {
	return &plugin{}
}
