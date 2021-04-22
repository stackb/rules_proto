package protoc

import (
	"path"
)

type protoRuleConfig struct {
	pattern string
	exclude bool
}

func (c *protoRuleConfig) IsRuleIncluded(name string) bool {
	// ignoring error here as we would have already panic'ed if the pattern was bad.
	match, _ := path.Match(c.pattern, name)
	if c.exclude {
		return !match
	}
	return match
}

func (c *protoRuleConfig) IsRuleExcluded(name string) bool {
	// ignoring error here as we would have already panic'ed if the pattern was bad.
	match, _ := path.Match(c.pattern, name)
	if c.exclude {
		return match
	}
	return !match
}
