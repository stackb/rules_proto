package protoc

import (
	"testing"

	"github.com/stackb/rules_proto/pkg/goldentest"
)

func TestProtoc(t *testing.T) {
	goldentest.FromDir("gazelle/protoc").Run(t, "gazelle-protoc")
}
