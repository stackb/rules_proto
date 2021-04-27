package golden

import "testing"

func TestGoldens(t *testing.T) {
	NewTestHarness("cmd/protoc/golden", "gazelle").Run(t)
}
