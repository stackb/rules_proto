package protoc

import (
	"fmt"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/emicklei/proto"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func isStarlarkPlugin(filename string) bool {
	return strings.HasSuffix(filename, ".starlark")
}

func loadStarlarkPlugin(name, filename string, src interface{}, reporter func(msg string), errorReporter func(err error)) (Plugin, error) {

	newErrorf := func(msg string, args ...interface{}) error {
		err := fmt.Errorf(filename+": "+msg, args...)
		errorReporter(err)
		return err
	}

	plugins := make(map[string]*starlarkstruct.Struct)

	module := &starlarkstruct.Module{
		Name: "protoc",
		Members: starlark.StringDict{
			"Plugin":              starlark.NewBuiltin("Plugin", newStarlarkPlugin(plugins)),
			"PluginConfiguration": starlark.NewBuiltin("PluginConfiguration", newStarlarkPluginConfiguration()),
		},
	}

	predeclared := starlark.StringDict{
		"protoc": module,
	}

	_, program, err := starlark.SourceProgram(filename, src, predeclared.Has)
	if err != nil {
		return nil, err
	}

	thread := new(starlark.Thread)
	thread.Print = func(thread *starlark.Thread, msg string) {
		reporter(msg)
	}
	_, err = program.Init(thread, predeclared)
	if err != nil {
		return nil, newErrorf("eval: %w", err)
	}

	if plugin, ok := plugins[name]; !ok {
		return nil, newErrorf("plugin %q was never declared", name)
	} else {
		return &starlarkPlugin{
			name:          name,
			plugin:        plugin,
			reporter:      thread.Print,
			errorReporter: newErrorf,
		}, nil
	}
}

// starlarkPlugin is a starlark plugin that implements the protoc plugin interface.
type starlarkPlugin struct {
	name          string
	reporter      func(thread *starlark.Thread, msg string)
	errorReporter func(msg string, args ...interface{}) error
	plugin        *starlarkstruct.Struct
}

func (p *starlarkPlugin) Name() string {
	return p.name
}

func (p *starlarkPlugin) Configure(ctx *PluginContext) *PluginConfiguration {

	var result *PluginConfiguration

	configure, err := p.plugin.Attr("configure")
	if err != nil {
		p.errorReporter("plugin %q has no configure function", p.name)
		return nil
	}

	thread := new(starlark.Thread)
	thread.Print = p.reporter
	value, err := starlark.Call(thread, configure, starlark.Tuple{
		newPluginContextStruct(ctx),
	}, []starlark.Tuple{})
	if err != nil {
		p.errorReporter("plugin %q configure failed: %w", p.name, err)
		return nil
	}

	switch value := value.(type) {
	case *starlarkstruct.Struct:
		labelValue, err := value.Attr("label")
		if err != nil {
			p.errorReporter("PluginConfiguration.label get value: %v", err)
			return nil
		}
		lbl, err := label.Parse(labelValue.(starlark.String).GoString())
		if err != nil {
			p.errorReporter("PluginConfiguration.label parse: %v", err)
			return nil
		}

		outputsValue, err := value.Attr("outputs")
		if err != nil {
			p.errorReporter("PluginConfiguration.outputs get value: %v", err)
		}
		outputsList := outputsValue.(*starlark.List)
		outputs := make([]string, outputsList.Len())
		for i := 0; i < outputsList.Len(); i++ {
			outputs[i] = outputsList.Index(i).(starlark.String).GoString()
		}
		result = &PluginConfiguration{
			Label:   lbl,
			Outputs: outputs,
		}
	default:
		p.errorReporter("plugin %q configure returned invalid type: %T", p.name, value)
		return nil
	}

	return result
}

type goStarlarkFunction func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error)

func newStarlarkPluginConfiguration() goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var labelStr string
		var out string
		outputs := &starlark.List{}
		options := &starlark.List{}

		if err := starlark.UnpackArgs("PluginConfiguration", args, kwargs,
			"label", &labelStr,
			"outputs", &outputs,
			"out?", &out,
			"options?", &options,
		); err != nil {
			return nil, err
		}

		lbl, err := label.Parse(labelStr)
		if err != nil {
			return nil, fmt.Errorf("invalid label: %w", err)
		}

		return starlarkstruct.FromStringDict(
			Symbol("PluginConfiguration"),
			starlark.StringDict{
				"label":   starlark.String(lbl.String()),
				"outputs": outputs,
				"out":     starlark.String(out),
				"options": options,
			},
		), nil
	}
}

func newStarlarkPlugin(plugins map[string]*starlarkstruct.Struct) goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var configure starlark.Callable

		if err := starlark.UnpackArgs("Plugin", args, kwargs,
			"name", &name,
			"configure", &configure,
		); err != nil {
			return nil, err
		}

		plugin := starlarkstruct.FromStringDict(
			Symbol("Plugin"),
			starlark.StringDict{
				"name":      starlark.String(name),
				"configure": configure,
			},
		)

		plugins[name] = plugin
		return plugin, nil
	}
}

func newPluginContextStruct(ctx *PluginContext) *starlarkstruct.Struct {

	return starlarkstruct.FromStringDict(
		Symbol("PluginContext"),
		starlark.StringDict{
			"rel":            starlark.String(ctx.Rel),
			"plugin_config":  newPluginConfigStruct(ctx.PluginConfig),
			"package_config": newPackageConfigStruct(ctx.PackageConfig),
			"proto_library":  newProtoLibraryStruct(ctx.ProtoLibrary),
		},
	)
}

func newPluginConfigStruct(cfg LanguagePluginConfig) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("LanguagePluginConfig"),
		starlark.StringDict{
			"name":           starlark.String(cfg.Name),
			"implementation": starlark.String(cfg.Implementation),
			"label":          starlark.String(cfg.Label.String()),
			"options":        newStringBoolDict(cfg.Options),
			"deps":           newStringBoolDict(cfg.Deps),
			"enabled":        starlark.Bool(cfg.Enabled),
		},
	)
}

func newPackageConfigStruct(cfg PackageConfig) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("PackageConfig"),
		starlark.StringDict{
			"config": newConfigStruct(cfg.Config),
		},
	)
}

func newConfigStruct(c *config.Config) *starlarkstruct.Struct {
	if c == nil {
		return starlarkstruct.FromStringDict(
			Symbol("Config"),
			starlark.StringDict{
				"work_dir":  starlark.String(""),
				"repo_root": starlark.String(""),
				"repo_name": starlark.String(""),
			},
		)
	}
	return starlarkstruct.FromStringDict(
		Symbol("Config"),
		starlark.StringDict{
			"work_dir":  starlark.String(c.WorkDir),
			"repo_root": starlark.String(c.RepoRoot),
			"repo_name": starlark.String(c.RepoName),
		},
	)
}

func newProtoLibraryStruct(p ProtoLibrary) *starlarkstruct.Struct {
	if p == nil {
		return starlarkstruct.FromStringDict(
			Symbol("ProtoLibrary"),
			starlark.StringDict{
				"name":                starlark.String(""),
				"base_name":           starlark.String(""),
				"strip_import_prefix": starlark.String(""),
				"srcs":                &starlark.List{},
				"deps":                &starlark.List{},
				"imports":             &starlark.List{},
				"files":               &starlark.List{},
			},
		)
	}
	return starlarkstruct.FromStringDict(
		Symbol("ProtoLibrary"),
		starlark.StringDict{
			"name":                starlark.String(p.Name()),
			"base_name":           starlark.String(p.BaseName()),
			"strip_import_prefix": starlark.String(p.StripImportPrefix()),
			"srcs":                newStringList(p.Srcs()),
			"deps":                newStringList(p.Deps()),
			"imports":             newStringList(p.Imports()),
			"files":               newProtoFileList(p.Files()),
		},
	)
}

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

func newProtoFileList(in []*File) *starlark.List {
	values := make([]starlark.Value, 0, len(in))
	for _, v := range in {
		if v == nil {
			continue
		}
		values = append(values, newProtoFileStruct(*v))
	}
	return starlark.NewList(values)
}

func newProtoFileStruct(f File) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("ProtoFile"),
		starlark.StringDict{
			"dir":          starlark.String(f.Dir),
			"basename":     starlark.String(f.Basename),
			"name":         starlark.String(f.Name),
			"pkg":          newProtoPackageStruct(f.pkg),
			"imports":      newProtoImportList(f.imports),
			"options":      newProtoOptionList(f.options),
			"messages":     newProtoMessageList(f.messages),
			"services":     newProtoServiceList(f.services),
			"enums":        newProtoEnumList(f.enums),
			"enum_options": newProtoEnumOptionList(f.enumOptions),
		},
	)
}

func newProtoPackageStruct(p proto.Package) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("ProtoPackage"),
		starlark.StringDict{
			"name": starlark.String(p.Name),
		},
	)
}

func newProtoImportList(in []proto.Import) *starlark.List {
	values := make([]starlark.Value, len(in))
	for i, v := range in {
		values[i] = newProtoImportStruct(v)
	}
	return starlark.NewList(values)
}

func newProtoImportStruct(i proto.Import) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("ProtoImport"),
		starlark.StringDict{
			"filename": starlark.String(i.Filename),
			"kind":     starlark.String(i.Kind),
		},
	)
}

func newProtoOptionList(in []proto.Option) *starlark.List {
	values := make([]starlark.Value, len(in))
	for i, v := range in {
		values[i] = newProtoOptionStruct(v)
	}
	return starlark.NewList(values)
}

func newProtoOptionStruct(o proto.Option) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("ProtoOption"),
		starlark.StringDict{
			"name":     starlark.String(o.Name),
			"constant": starlark.String(o.Constant.Source),
		},
	)
}

func newProtoMessageList(in []proto.Message) *starlark.List {
	values := make([]starlark.Value, len(in))
	for i, v := range in {
		values[i] = newProtoMessageStruct(v)
	}
	return starlark.NewList(values)
}

func newProtoMessageStruct(m proto.Message) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("ProtoMessage"),
		starlark.StringDict{
			"name":      starlark.String(m.Name),
			"is_extend": starlark.Bool(m.IsExtend),
		},
	)
}

func newProtoServiceList(in []proto.Service) *starlark.List {
	values := make([]starlark.Value, len(in))
	for i, v := range in {
		values[i] = newProtoServiceStruct(v)
	}
	return starlark.NewList(values)
}

func newProtoServiceStruct(s proto.Service) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("ProtoService"),
		starlark.StringDict{
			"name": starlark.String(s.Name),
		},
	)
}

func newProtoEnumList(in []proto.Enum) *starlark.List {
	values := make([]starlark.Value, len(in))
	for i, v := range in {
		values[i] = newProtoEnumStruct(v)
	}
	return starlark.NewList(values)
}

func newProtoEnumStruct(e proto.Enum) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("ProtoEnum"),
		starlark.StringDict{
			"name": starlark.String(e.Name),
		},
	)
}

func newProtoEnumOptionList(in []proto.Option) *starlark.List {
	values := make([]starlark.Value, len(in))
	for i, v := range in {
		values[i] = newProtoEnumOptionStruct(v)
	}
	return starlark.NewList(values)
}

func newProtoEnumOptionStruct(e proto.Option) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("ProtoEnumOption"),
		starlark.StringDict{
			"name":     starlark.String(e.Name),
			"constant": starlark.String(e.Constant.Source),
		},
	)
}

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
