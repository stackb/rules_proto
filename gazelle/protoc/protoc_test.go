package protoc

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/golden"
)

func TestProtoc(t *testing.T) {
	golden.NewTestDataDir("gazelle/protoc").Run(t, "protoc")
}
