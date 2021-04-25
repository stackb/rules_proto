package protoc

import "github.com/bazelbuild/bazel-gazelle/rule"

// ProtoRule is capable of taking a compilation and deriving another rule based
// on it.  For example, a java_proto_library ProtoRule implementation might
// collect all the emitted *.srcjar files from the protoc configuration and wrap
// it with a java_library.
type ProtoRule interface {
	// LoadInfo returns the gazelle LoadInfo.
	LoadInfo() rule.LoadInfo
	// KindInfo returns the gazelle KindInfo.
	KindInfo() rule.KindInfo
	// GenerateRule takes the given configration and compilation and emits a
	// RuleProvider
	GenerateRule(rc *ProtoRuleConfig, pc *ProtocConfiguration) RuleProvider
}
