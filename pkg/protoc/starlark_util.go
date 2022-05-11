package protoc

import (
	"go.starlark.net/starlark"
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

func newStringList(in []string) *starlark.List {
	values := make([]starlark.Value, len(in))
	for i, v := range in {
		values[i] = starlark.String(v)
	}
	return starlark.NewList(values)
}
