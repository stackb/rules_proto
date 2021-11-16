package protoc

import (
	"strings"
)

// intent represents an action for an attribute name or "value" optionally
// prefixed by a '+' or '-'.  If the prefix is missing, the intent is not
// negative.
type intent struct {
	Value string
	Want  bool
}

func parseIntent(value string) *intent {
	value = strings.TrimSpace(value)
	negative := strings.HasPrefix(value, "-")
	positive := strings.HasPrefix(value, "+")
	if negative || positive {
		value = value[1:]
	}
	return &intent{Value: value, Want: !negative}
}
