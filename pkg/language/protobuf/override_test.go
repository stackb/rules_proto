package protobuf

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/google/go-cmp/cmp"

	"github.com/stackb/rules_proto/pkg/protoc"
)

// TestOverrideRule demonstrates the shape of an override rule: as a carrier for
// ProtoLibrary instances in a PrivateAttr.  The proto_library rules inside it
// might have go_googleapis deps that we want to scrub out and replace with
// locally-resolved ones.  This is to get around gazelle's hardcoded resolver
// strategy for these labels.
func TestOverrideRule(t *testing.T) {
	// rel is the package path.  deps are the input deps on the proto_library
	// rule. imps are the imports of that proto_library rule's file(s). want are
	// the expected deps on the rule.  resolve is mapping from the imp to a
	// label.
	for name, tc := range map[string]struct {
		rel              string
		deps, imps, want []string
		known            map[string]label.Label
	}{
		// "empty": {
		// 	rel:  "",
		// 	want: []string{},
		// },
		// "no go_googleapis deps": {
		// 	rel:  "",
		// 	deps: []string{"@com_google_protobuf//:any_proto"},
		// 	want: []string{"@com_google_protobuf//:any_proto"},
		// },
		"has go_googleapis dep": {
			rel:  "",
			deps: []string{"@com_google_protobuf//:any_proto", "@go_googleapis//google/api:api_proto"},
			want: []string{"//google/api:http_proto", "@com_google_protobuf//:any_proto"},
			imps: []string{"google/api/http.proto"},
			known: map[string]label.Label{
				"google/api/http.proto": label.New("", "google/api", "http_proto"),
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := protoc.NewImportResolver(&protoc.ImportResolverOptions{
				Printf: t.Logf,
				Debug:  true,
			})
			for k, v := range tc.known {
				resolver.Provide("proto", "proto", k, v)
			}
			r := makeProtoLibraryRule("test_proto", tc.deps, tc.imps)
			lib := makeOtherProtoLibrary(r)
			overrideRule := makeProtoOverrideRule([]protoc.ProtoLibrary{lib})
			resolveOverrideRule(tc.rel, overrideRule, resolver)

			got := r.AttrStrings("deps")
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("resolverOverrideRule() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func makeProtoLibraryRule(name string, deps, imps []string) *rule.Rule {
	r := rule.NewRule("proto_library", name)
	r.SetAttr("deps", deps)
	r.SetPrivateAttr(config.GazelleImportsKey, imps)
	return r
}

func makeOtherProtoLibrary(r *rule.Rule) protoc.ProtoLibrary {
	f := rule.EmptyFile("", "")
	return protoc.NewOtherProtoLibrary(f, r)
}
