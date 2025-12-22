package golden

import (
	"testing"

	"github.com/stackb/rules_proto/v4/pkg/goldentest"
)

func TestGoldens(t *testing.T) {
	goldentest.
		FromDir("example/golden").
		Run(t, "gazelle")
}
