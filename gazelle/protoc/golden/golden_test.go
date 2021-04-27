package golden

import "testing"

func TestGoldens(t *testing.T) {
	NewTestHarness("gazelle/protoc/golden", "gazelle").Run(t)
}
