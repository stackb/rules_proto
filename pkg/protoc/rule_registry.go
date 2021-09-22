package protoc

import "errors"

// ErrUnknownRule is the error returned when a rule is not known.
var ErrUnknownRule = errors.New("unknown rule")

// RuleRegistry represents a library of rule implementations.
type RuleRegistry interface {
	// RuleNames returns a sorted list of rule names.
	RuleNames() []string
	// LookupRule returns the implementation under the given name.  If the rule
	// is not found, ErrUnknownRule is returned.
	LookupRule(name string) (LanguageRule, error)
	// MustRegisterRule installs a LanguageRule implementation under the given
	// name in the global rule registry.  Panic will occur if the same rule is
	// registered multiple times.
	MustRegisterRule(name string, rule LanguageRule) RuleRegistry
}
