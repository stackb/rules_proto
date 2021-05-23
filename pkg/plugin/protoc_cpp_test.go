package plugin

import (
	"testing"
)

func TestProtocCppPlugin(t *testing.T) {
	PluginTestCases(t, &ProtocCppPlugin{}, map[string]PluginTestCase{
		"empty file": {},
		"message in root package": {
			Input:       "message Test{}",
			ShouldApply: true,
			Outputs:     []string{"test.pb.cc", "test.pb.h"},
		},
		"message with package": {
			Input:       "package p;\nmessage Test{}",
			ShouldApply: true,
			Outputs:     []string{"p/test.pb.cc", "p/test.pb.h"},
		},
		"only services": {
			Input:       "service Test{}",
			ShouldApply: false,
		},
	})
}
