package protoc

import (
	"fmt"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type goStarlarkFunction func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

// Symbol is the type of a Starlark constructor symbol.  It prints more
// favorably than a starlark.String.
type Symbol string

func (s Symbol) String() string             { return string(s) }
func (s Symbol) GoString() string           { return string(s) }
func (s Symbol) Type() string               { return "symbol" }
func (s Symbol) Freeze()                    {} // immutable
func (s Symbol) Truth() starlark.Bool       { return len(s) > 0 }
func (s Symbol) Hash() (uint32, error)      { return starlark.String(s).Hash() }
func (s Symbol) Len() int                   { return len(s) } // bytes
func (s Symbol) Index(i int) starlark.Value { return s[i : i+1] }

func newStringBoolDict(in map[string]bool) *starlark.Dict {
	out := &starlark.Dict{}
	for k, v := range in {
		out.SetKey(starlark.String(k), starlark.Bool(v))
	}
	return out
}

func newStringStringDict(in map[string]string) *starlark.Dict {
	out := &starlark.Dict{}
	for k, v := range in {
		out.SetKey(starlark.String(k), starlark.String(v))
	}
	return out
}

func newStringList(in []string) *starlark.List {
	values := make([]starlark.Value, len(in))
	for i, v := range in {
		values[i] = starlark.String(v)
	}
	return starlark.NewList(values)
}

func newPredeclared(plugins, rules map[string]*starlarkstruct.Struct) starlark.StringDict {
	protoc := &starlarkstruct.Module{
		Name: "protoc",
		Members: starlark.StringDict{
			"Plugin":              starlark.NewBuiltin("Plugin", newStarlarkPluginFunction(plugins)),
			"Rule":                starlark.NewBuiltin("Rule", newStarlarkLanguageRuleFunction(rules)),
			"PluginConfiguration": starlark.NewBuiltin("PluginConfiguration", newStarlarkPluginConfiguration()),
		},
	}

	gazelle := &starlarkstruct.Module{
		Name: "gazelle",
		Members: starlark.StringDict{
			"Rule":     starlark.NewBuiltin("Rule", newGazelleRuleFunction()),
			"LoadInfo": starlark.NewBuiltin("LoadInfo", newGazelleLoadInfoFunction()),
			"KindInfo": starlark.NewBuiltin("KindInfo", newGazelleKindInfoFunction()),
		},
	}

	return starlark.StringDict{
		protoc.Name:  protoc,
		gazelle.Name: gazelle,
		"struct":     starlark.NewBuiltin("struct", starlarkstruct.Make),
	}
}

func loadStarlarkProgram(filename string, src interface{}, predeclared starlark.StringDict, reporter func(msg string), errorReporter func(err error)) (*starlark.StringDict, *starlark.Thread, error) {
	newErrorf := func(msg string, args ...interface{}) error {
		err := fmt.Errorf(filename+": "+msg, args...)
		errorReporter(err)
		return err
	}

	_, program, err := starlark.SourceProgram(filename, src, predeclared.Has)
	if err != nil {
		return nil, nil, newErrorf("source program error: %v", err)
	}

	thread := new(starlark.Thread)
	thread.Print = func(thread *starlark.Thread, msg string) {
		reporter(msg)
	}
	globals, err := program.Init(thread, predeclared)
	if err != nil {
		return nil, nil, newErrorf("eval: %w", err)
	}

	return &globals, thread, nil
}

func newGazelleRuleFunction() goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name, kind string
		var attrs starlark.Dict

		if err := starlark.UnpackArgs("Rule", args, kwargs,
			"name", &name,
			"kind", &kind,
			"attrs?", &attrs,
		); err != nil {
			return nil, err
		}

		value := starlarkstruct.FromStringDict(
			Symbol("Rule"),
			starlark.StringDict{
				"name":  starlark.String(name),
				"kind":  starlark.String(kind),
				"attrs": &attrs,
			},
		)

		return value, nil
	}
}

func newGazelleLoadInfoFunction() goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var symbols, after starlark.List

		if err := starlark.UnpackArgs("LoadInfo", args, kwargs,
			"name", &name,
			"symbols", &symbols,
			"after", &after,
		); err != nil {
			return nil, err
		}

		value := starlarkstruct.FromStringDict(
			Symbol("LoadInfo"),
			starlark.StringDict{
				"name":    starlark.String(name),
				"symbols": &symbols,
				"after":   &after,
			},
		)

		return value, nil
	}
}

func newGazelleKindInfoFunction() goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var matchAny bool
		var matchAttrs, nonEmptyAttrs, substituteAttrs, mergeableAttrs, resolveAttrs starlark.Dict

		if err := starlark.UnpackArgs("KindInfo", args, kwargs,
			"match_any?", &matchAny,
			"match_attrs?", &matchAttrs,
			"non_empty_attrs?", &nonEmptyAttrs,
			"substitute_attrs?", &substituteAttrs,
			"mergeable_attrs?", &mergeableAttrs,
			"resolve_attrs?", &resolveAttrs,
		); err != nil {
			return nil, err
		}

		value := starlarkstruct.FromStringDict(
			Symbol("KindInfo"),
			starlark.StringDict{
				"match_any":        starlark.Bool(matchAny),
				"match_attrs":      &matchAttrs,
				"non_empty_attrs":  &nonEmptyAttrs,
				"substitute_attrs": &substituteAttrs,
				"mergeable_attrs":  &mergeableAttrs,
				"resolve_attrs":    &resolveAttrs,
			},
		)

		return value, nil
	}
}
