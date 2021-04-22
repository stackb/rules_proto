package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// ProtoDescriptorSet implements a ruleProvider for the @rules_proto
// "proto_descriptor_set" rule.
type ProtoDescriptorSet struct {
	library ProtoLibrary
}

// NewProtoDescriptorSet constructs a ProtoDescriptorSet based on the
// ProtoLibrary rule on which it depends.
func NewProtoDescriptorSet(library ProtoLibrary) *ProtoDescriptorSet {
	return &ProtoDescriptorSet{
		library: library,
	}
}

// Kind implements part of the ruleProvider interface
func (s *ProtoDescriptorSet) Kind() string {
	return "proto_descriptor_set"
}

// Rule implements part of the ruleProvider interface
func (s *ProtoDescriptorSet) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())
	newRule.SetAttr("visibility", s.Visibility())
	newRule.SetAttr("deps", s.Deps())
	return newRule
}

// KindInfo implements part of the ruleProvider interface
func (s *ProtoDescriptorSet) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		NonEmptyAttrs:  map[string]bool{"deps": true},
		MergeableAttrs: map[string]bool{},
	}
}

// Imports implements part of the ruleProvider interface
func (s *ProtoDescriptorSet) Imports() []string {
	return []string{s.Kind()}
}

// Name implements part of the ruleProvider interface
func (s *ProtoDescriptorSet) Name() string {
	return s.library.Name() + "_descriptor"
}

// Visibility implements part of the ruleProvider interface
func (s *ProtoDescriptorSet) Visibility() []string {
	return []string{"//visibility:public"}
}

// Deps computes the dependencies of the rule.
func (s *ProtoDescriptorSet) Deps() []string {
	return []string{":" + s.library.Name()}
}

// Resolve implements part of the RuleProvider interface.
func (s *ProtoDescriptorSet) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
}
