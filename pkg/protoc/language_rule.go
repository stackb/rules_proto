package protoc

import "github.com/bazelbuild/bazel-gazelle/rule"

// LanguageRule is capable of taking a compilation and deriving another rule
// based on it.  For example, a java_proto_library LanguageRule implementation
// might collect all the emitted *.srcjar files from the protoc configuration
// and wrap it with a java_library.
type LanguageRule interface {
	// Name returns the name of the rule
	Name() string
	// LoadInfo returns the gazelle LoadInfo.
	LoadInfo() rule.LoadInfo
	// KindInfo returns the gazelle KindInfo.
	KindInfo() rule.KindInfo
	// ProvideRule takes the given configration and compilation and emits a
	// RuleProvider.  If the state of the ProtocConfiguration is such that the
	// rule should not be emitted, implementation should return nil.
	ProvideRule(rc *LanguageRuleConfig, pc *ProtocConfiguration) RuleProvider
}
