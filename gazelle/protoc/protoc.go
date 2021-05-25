package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/language"

	extension "github.com/stackb/rules_proto/pkg/gazelle/protoc"

	_ "github.com/stackb/rules_proto/pkg/builtin"
)

// NewLanguage is called by Gazelle to install this language extension in a
// binary.
func NewLanguage() language.Language {
	return extension.NewProtoc("protoc")
}
