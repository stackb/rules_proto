package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type RuleIndex interface {
	// Put records the association of the rule under the given label.
	Put(label.Label, *rule.Rule)
	// Get returns the rule under the given label, or nil if not known.
	Get(label.Label) *rule.Rule
}

// GlobalRuleIndex returns a reference to the global RuleIndex
func GlobalRuleIndex() RuleIndex {
	return globalRuleIndex
}

// globalRuleIndex is the default resolver singleton.
var globalRuleIndex = &ruleIndex{
	rules: make(map[label.Label]*rule.Rule),
}

// ruleIndex implements RuleIndex.
type ruleIndex struct {
	rules map[label.Label]*rule.Rule
}

func (r *ruleIndex) Put(from label.Label, rule *rule.Rule) {
	r.rules[from] = rule
}

func (r *ruleIndex) Get(from label.Label) *rule.Rule {
	return r.rules[from]
}
