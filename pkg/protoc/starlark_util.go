package protoc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

type errorReporter func(format string, args ...interface{}) error

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

func newStringListDict(in map[string]map[string]bool) *starlark.Dict {
	out := &starlark.Dict{}
	return out
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

func resolveStarlarkFilename(workDir, filename string) (string, error) {
	if filename == "" {
		return "", fmt.Errorf("filename is empty")
	}
	if _, err := os.Stat(filepath.Join(workDir, filename)); !errors.Is(err, os.ErrNotExist) {
		return filepath.Join(workDir, filename), nil
	}

	dirname := workDir
	// looking for a file named 'DO_NOT_BUILD_HERE' in a parent directory.  The
	// contents of this file names the sourceRoot directory.
	var sourceRoot string
	for dirname != "." {
		sourceRootFile := filepath.Join(dirname, "DO_NOT_BUILD_HERE")
		if _, err := os.Stat(sourceRootFile); errors.Is(err, os.ErrNotExist) {
			dirname = filepath.Dir(dirname)
		} else {
			data, err := ioutil.ReadFile(sourceRootFile)
			if err != nil {
				return "", fmt.Errorf("failed to read DO_NOT_BUILD_HERE file: %w", err)
			}
			sourceRoot = strings.TrimSpace(string(data))
			break
		}
	}
	if sourceRoot == "" {
		return "", fmt.Errorf("failed to find sourceRoot")
	}
	return filepath.Join(sourceRoot, filename), nil
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
		attrs := new(starlark.Dict)

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
				"attrs": attrs,
			},
		)

		return value, nil
	}
}

func newGazelleLoadInfoFunction() goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		symbols := new(starlark.List)
		after := new(starlark.List)

		if err := starlark.UnpackArgs("LoadInfo", args, kwargs,
			"name", &name,
			"symbols", &symbols,
			"after?", &after,
		); err != nil {
			return nil, err
		}

		value := starlarkstruct.FromStringDict(
			Symbol("LoadInfo"),
			starlark.StringDict{
				"name":    starlark.String(name),
				"symbols": symbols,
				"after":   after,
			},
		)

		return value, nil
	}
}

func newGazelleKindInfoFunction() goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var matchAny bool
		matchAttrs := new(starlark.List)
		nonEmptyAttrs := new(starlark.Dict)
		substituteAttrs := new(starlark.Dict)
		mergeableAttrs := new(starlark.Dict)
		resolveAttrs := new(starlark.Dict)

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
				"match_attrs":      matchAttrs,
				"non_empty_attrs":  nonEmptyAttrs,
				"substitute_attrs": substituteAttrs,
				"mergeable_attrs":  mergeableAttrs,
				"resolve_attrs":    resolveAttrs,
			},
		)

		return value, nil
	}
}

func structAttrStringSlice(in *starlarkstruct.Struct, name string, errorReporter errorReporter) []string {
	value, err := in.Attr(name)
	if err != nil {
		errorReporter("getting struct attr %s: %w", err)
		return nil
	}
	list, ok := value.(*starlark.List)
	if !ok {
		errorReporter("%s is not a list", name)
		return nil
	}
	return stringSlice(list, errorReporter)
}

func stringSlice(list *starlark.List, errorReporter errorReporter) (out []string) {
	for i := 0; i < list.Len(); i++ {
		value := list.Index(i)
		switch value := (value).(type) {
		case starlark.String:
			out = append(out, value.GoString())
		default:
			errorReporter("list[%d]: expected string, got %s", i, value.Type())
		}
	}
	return
}

func structAttrBool(in *starlarkstruct.Struct, name string, errorReporter errorReporter) (out bool) {
	value, err := in.Attr(name)
	if err != nil {
		errorReporter("getting struct attr %s: %w", err)
		return
	}
	if value == nil {
		return false
	}
	switch t := value.(type) {
	case starlark.Bool:
		out = bool(t.Truth())
	default:
		errorReporter("attr %q: want bool, got %T", name, value)
	}
	return
}

func structAttrString(in *starlarkstruct.Struct, name string, errorReporter errorReporter) string {
	value, err := in.Attr(name)
	if err != nil {
		errorReporter("getting struct attr %s: %w", err)
		return ""
	}
	switch value := value.(type) {
	case starlark.String:
		return value.GoString()
	default:
		errorReporter("%s is not a string (%T)", name, value)
		return ""
	}
}

func structAttrMapStringBool(in *starlarkstruct.Struct, name string, errorReporter errorReporter) (out map[string]bool) {
	value, err := in.Attr(name)
	if err != nil {
		if _, ok := err.(starlark.NoSuchAttrError); ok {
			return
		}
		errorReporter("%v", err)
		return
	}
	if value == nil {
		return
	}
	dict, ok := value.(*starlark.Dict)
	if !ok {
		errorReporter("%v.%s: value must have type starlark.Dict (got %T)", in.Constructor(), name, value)
		return
	}
	out = make(map[string]bool, dict.Len())
	for _, key := range dict.Keys() {
		k, ok := key.(starlark.String)
		if !ok {
			errorReporter("%v.%s: dict keys must have type string (got %T)", in.Constructor(), name, key)
			return
		}
		if value, ok, err := dict.Get(key); ok && err == nil {
			b, ok := value.(starlark.Bool)
			if !ok {
				errorReporter("%v.%s: dict value for %q must have type bool (got %T)", in.Constructor(), name, k.GoString(), value)
				return
			}
			out[k.GoString()] = bool(b.Truth())
		}
	}
	return
}
