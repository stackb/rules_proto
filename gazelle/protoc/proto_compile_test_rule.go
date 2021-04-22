package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// ProtoCompileTest implements a ruleProvider for the @build_stack_rules_proto
// "proto_compile_test".  When under "bazel test TARGET", the rule asserts that
// the (generated) srcs match the generated ones.  When under "bazel run
// TARGET.update", the generated files are copied over to the source location.
type ProtoCompileTest struct {
	ProtoRule *ProtoRule
}

// Kind implements part of the ruleProvider interface.
func (s *ProtoCompileTest) Kind() string {
	return "proto_compile_test"
}

// Name implements part of the ruleProvider interface.
func (s *ProtoCompileTest) Name() string {
	return s.ProtoRule.Name() + "_test"
}

// Visibility implements part of the ruleProvider interface.
func (s *ProtoCompileTest) Visibility() []string {
	return []string{"//visibility:private"}
}

// Imports implements part of the ruleProvider interface.
func (s *ProtoCompileTest) Imports() []string {
	return []string{s.Kind()}
}

// Rule implements part of the ruleProvider interface.
func (s *ProtoCompileTest) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())
	newRule.SetAttr("rule", ":"+s.ProtoRule.Name())
	newRule.SetAttr("srcs", s.Srcs())
	newRule.SetAttr("visibility", s.Visibility())
	return newRule
}

// KindInfo implements part of the ruleProvider interface.
func (s *ProtoCompileTest) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		NonEmptyAttrs:  map[string]bool{"deps": true},
		MergeableAttrs: map[string]bool{},
	}
}

// Srcs provides the input sources.  In this case they are provided by the
// underlying ProtoRule.
func (s *ProtoCompileTest) Srcs() []string {
	return s.ProtoRule.GeneratedSrcs()
}

// Resolve implements part of the RuleProvider interface.
func (s *ProtoCompileTest) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
}
