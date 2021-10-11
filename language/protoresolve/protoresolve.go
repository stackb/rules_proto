package protoresolve

import (
	"github.com/bazelbuild/bazel-gazelle/language"

	"github.com/stackb/rules_proto/pkg/language/protoresolve"
)

// NewLanguage is called by Gazelle to install this language extension in a
// binary.
func NewLanguage() language.Language {
	return protoresolve.NewProtoIndexLanguage("protoresolve")
}
