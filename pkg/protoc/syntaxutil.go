package protoc

import (
	"sort"

	"github.com/bazelbuild/buildtools/build"
)

// MakeStringListDict converts the given string -> string_list into a build
// expression.
func MakeStringListDict(in map[string][]string) build.Expr {
	items := make([]*build.KeyValueExpr, 0)
	keys := make([]string, 0)
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		values := in[key]
		sort.Strings(values)
		value := &build.ListExpr{List: make([]build.Expr, len(values))}
		for i, val := range values {
			value.List[i] = &build.StringExpr{Value: val}
		}
		items = append(items, &build.KeyValueExpr{
			Key:   &build.StringExpr{Value: key},
			Value: value,
		})
	}
	return &build.DictExpr{List: items}
}

// MakeStringDict converts the given string -> string into a build expression.
func MakeStringDict(in map[string]string) build.Expr {
	dict := &build.DictExpr{}
	keys := make([]string, 0)
	for k := range in {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		value := in[key]
		dict.List = append(dict.List, &build.KeyValueExpr{
			Key:   &build.StringExpr{Value: key},
			Value: &build.StringExpr{Value: value},
		})
	}
	return dict
}
