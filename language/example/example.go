package example

import (
	"github.com/bazelbuild/bazel-gazelle/language"

	"github.com/stackb/rules_proto/pkg/language/noop"
	// Put your own imports here!
	// _ "github.com/org/repo/pkg/plugin/foo"
)

// NewLanguage is called by Gazelle to install this language extension in a
// binary.
func NewLanguage() language.Language {
	return noop.NewNoOpLanguage("my-project-name")
}
