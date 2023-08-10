package protoc

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func TestStructAttrMapStringBool(t *testing.T) {
	for name, tc := range map[string]struct {
		in      starlarkstruct.Struct
		name    string
		want    map[string]bool
		wantErr string
	}{
		"degenerate": {
			in: *starlarkstruct.FromStringDict(
				starlarkstruct.Default,
				starlark.StringDict{},
			),
		},
		"not found": {
			in: *starlarkstruct.FromStringDict(
				starlarkstruct.Default,
				starlark.StringDict{
					"a": starlark.False,
				},
			),
			name: "b",
		},
		"wrong value type": {
			in: *starlarkstruct.FromStringDict(
				starlarkstruct.Default,
				starlark.StringDict{
					"a": starlark.False,
				},
			),
			name:    "a",
			wantErr: `"struct".a: value must have type starlark.Dict (got starlark.Bool)`,
		},
		"wrong dict key": {
			in: *starlarkstruct.FromStringDict(
				starlarkstruct.Default,
				starlark.StringDict{
					"a": func() *starlark.Dict {
						d := starlark.NewDict(1)
						d.SetKey(starlark.False, starlark.False)
						return d
					}(),
				},
			),
			name:    "a",
			want:    map[string]bool{},
			wantErr: `"struct".a: dict keys must have type string (got starlark.Bool)`,
		},
		"wrong dict value": {
			in: *starlarkstruct.FromStringDict(
				starlarkstruct.Default,
				starlark.StringDict{
					"a": func() *starlark.Dict {
						d := starlark.NewDict(1)
						d.SetKey(starlark.String("merge_attrs"), starlark.String("true"))
						return d
					}(),
				},
			),
			name:    "a",
			want:    map[string]bool{},
			wantErr: `"struct".a: dict value for "merge_attrs" must have type bool (got starlark.String)`,
		},
		"correct dict value": {
			in: *starlarkstruct.FromStringDict(
				starlarkstruct.Default,
				starlark.StringDict{
					"a": func() *starlark.Dict {
						d := starlark.NewDict(1)
						d.SetKey(starlark.String("merge_attrs"), starlark.True)
						return d
					}(),
				},
			),
			name: "a",
			want: map[string]bool{"merge_attrs": true},
		},
	} {
		t.Run(name, func(t *testing.T) {
			var gotErr string
			got := structAttrMapStringBool(&tc.in, tc.name, func(format string, args ...interface{}) error {
				gotErr += fmt.Sprintf(format, args...)
				return nil
			})
			if diff := cmp.Diff(tc.wantErr, gotErr); diff != "" {
				t.Errorf("error (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("result (-want +got):\n%s", diff)
			}
		})
	}
}
