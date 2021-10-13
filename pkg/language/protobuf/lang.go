package protobuf

import (
	"github.com/stackb/rules_proto/pkg/protoc"
)

// NewprotobufLang create a new protobufLang Gazelle extension implementation.
func NewprotobufLang(name string) *protobufLang {
	return &protobufLang{
		name:  name,
		rules: protoc.Rules(),
	}
}

// protobufLang implements language.Language.
type protobufLang struct {
	// name of the extension
	name string
	// the rule registry
	rules protoc.RuleRegistry
	// configFiles contains yconfig yaml files to parse.  May be comma-separated.
	configFiles string
	// repoName is the name (if this an external repository)
	repoName string
	// importsOutFile is the name of the file to create.  If "", skip writing
	// the file.
	importsOutFile string
	// importsInFiles is a comma-separated list of files that contains proto
	// index csv content.
	importsInFiles string
	// overrideGoGooleapis performs special processing for go_googleapis deps
	overrideGoGooleapis bool
}

// Name implements part of the language.Language interface.
func (pl *protobufLang) Name() string { return pl.name }
