package protoc

import "errors"

// ErrUnknownRule is the error returned when a rule is not known.
var ErrUnknownRule = errors.New("unknown ProtoRule")

// ProtoRuleRegistry represents a library of rule implementations.
type ProtoRuleRegistry interface {
	// RuleNames returns a sorted list of rule names.
	RuleNames() []string
	// LookupProtoRule returns the implementation under the given name.
	LookupProtoRule(name string) (ProtoRule, error)
	// MustRegisterProtoRule installs a ProtoRule implementation under the given
	// name in the global rule registry.  Panic will occur of the same rule is
	// registered multiple times.
	MustRegisterProtoRule(name string, rule ProtoRule)
}

// RuleNames for the .
func RuleNames() []string {
	return protoc.RuleNames()
}

// LookupProtoRule returns the implementation under the given name.
func LookupProtoRule(name string) (ProtoRule, error) {
	return protoc.LookupProtoRule(name)
}

// MustRegisterProtoRule installs a ProtoRule implementation under the given
// name in the global rule registry.  Panic will occur of the same rule is
// registered multiple times.
func MustRegisterProtoRule(name string, rule ProtoRule) {
	protoc.MustRegisterProtoRule(name, rule)
}
