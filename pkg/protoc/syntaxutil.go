package protoc

import (
	"sort"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/buildtools/build"
)

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

// GetKeptFileRuleAttrString returns the value of the rule attribute IFF the
// backing File rule attribute has a '# keep' comment on it.
func GetKeptFileRuleAttrString(file *rule.File, r *rule.Rule, name string) string {
	if file == nil {
		return ""
	}
	assign := getRuleAssignExpr(file.File, r.Kind(), r.Name(), name)
	if assign == nil {
		return ""
	}
	if !rule.ShouldKeep(assign) {
		return ""
	}
	str, ok := assign.RHS.(*build.StringExpr)
	if !ok {
		return ""
	}
	return str.Value
}

// getRuleAssignExpr seeks through the file looking for call expressions having
// the given kind and rule name.  If found, the assignment expression having the
// given name is returned.  Otherwise, return nil.
func getRuleAssignExpr(file *build.File, kind string, name string, attrName string) *build.AssignExpr {
	if file == nil {
		return nil
	}
	for _, stmt := range file.Stmt {
		call, ok := stmt.(*build.CallExpr)
		if !ok {
			continue
		}
		ident, ok := call.X.(*build.Ident)
		if !ok {
			continue
		}
		if ident.Name != kind {
			continue
		}
		var ruleName string
		var want *build.AssignExpr
		for _, arg := range call.List {
			attr, ok := arg.(*build.AssignExpr)
			if !ok {
				continue
			}
			key := attr.LHS.(*build.Ident)
			switch key.Name {
			case "name":
				str, ok := attr.RHS.(*build.StringExpr)
				if ok {
					ruleName = str.Value
				}
			case attrName:
				want = attr
			}
		}
		if ruleName != name {
			continue
		}
		return want
	}
	return nil
}
