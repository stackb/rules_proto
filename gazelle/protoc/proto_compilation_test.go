package protoc

import "testing"

type getfilesTestCase struct {
	rel      string
	gensrcs  map[string][]string
	srcs     []string
	mappings map[string]map[string]string
}

func TestGenfiles(t *testing.T) {
	for name, tc := range map[string]genfileTestCase{} {
		t.Run(name, func(t *testing.T) {

		})
	}
}
