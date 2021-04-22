package protoc

// ProtoRuleOption is a function that modifies a ProtoRule upon construction.
type ProtoRuleOption func(*ProtoRule) *ProtoRule

// WithVisibility assigns the rule visibility
func WithVisibility(visibility []string) ProtoRuleOption {
	return func(rule *ProtoRule) *ProtoRule {
		rule.visibility = visibility
		return rule
	}
}

// WithPublicVisibility assigns the rule visibility to public
func WithPublicVisibility() ProtoRuleOption {
	return WithVisibility([]string{"//visibility:public"})
}

// WithPrivateVisibility assigns the rule visibility to private
func WithPrivateVisibility() ProtoRuleOption {
	return WithVisibility([]string{"//visibility:private"})
}

// WithComment assigns the rule comment.  It takes a varidic number of strings,
// each of which will be prepended with '# ' (as required by gazelle library).
func WithComment(lines ...string) ProtoRuleOption {
	return func(rule *ProtoRule) *ProtoRule {
		for _, line := range lines {
			rule.comment = append(rule.comment, "# "+line)
		}
		return rule
	}
}
