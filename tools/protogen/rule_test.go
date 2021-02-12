package protogen

import (
	"testing"
)

func TestCollectRuleDeps(t *testing.T) {

	cases := []struct {
		d             string
		rule          *ProtoRule
		expectedOrder []string
	}{
		{
			d:             "empty deps",
			rule:          &ProtoRule{},
			expectedOrder: []string{},
		},
		{
			d: "single rule deps",
			rule: &ProtoRule{
				Deps: []*ProtoDependency{
					{
						Name: "foo",
					},
				},
			},
			expectedOrder: []string{"foo"},
		},
		{
			d: "nodejs_grpc_library rule deps",
			rule: &ProtoRule{
				Name: "nodejs_grpc_library",
				Deps: []*ProtoDependency{
					{
						Name: "com_google_protobuf",
						Deps: []*ProtoDependency{
							{Name: "bazel_skylib"},
							{Name: "rules_python"},
							{Name: "zlib"},
						},
					},
				},
				Plugins: []*ProtoPlugin{
					{
						Name: "grpc",
						Deps: []*ProtoDependency{
							{
								Name: "grpc_js_node_modules",
								Deps: []*ProtoDependency{
									{Name: "build_bazel_rules_nodejs"},
								},
							},
							{
								Name: "grpc_tools_node_modules",
								Deps: []*ProtoDependency{
									{Name: "build_bazel_rules_nodejs"},
								},
							},
						},
					},
				},
			},
			expectedOrder: []string{
				"build_bazel_rules_nodejs",
				"grpc_js_node_modules",
				"grpc_tools_node_modules",
				"bazel_skylib",
				"rules_python",
				"zlib",
				"com_google_protobuf",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.d, func(t *testing.T) {
			actual := collectRuleDeps(tc.rule)
			for i, v := range actual {
				t.Logf("%d: %+v", i, v.Dep)
			}
			if len(actual) != len(tc.expectedOrder) {
				t.Fatalf("order len: want %d, got %d", len(tc.expectedOrder), len(actual))
			}

			for i, ruleDep := range actual {
				want := tc.expectedOrder[i]
				got := ruleDep.Dep.Name
				if want != got {
					t.Errorf("order %d: want %s, got %s", i, want, got)
				}
			}
		})
	}
}
