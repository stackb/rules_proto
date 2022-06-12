package protoc

import (
	"fmt"
	"os"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/emicklei/proto"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func LoadStarlarkPluginFromFile(workDir, filename, name string, reporter func(msg string), errorReporter func(err error)) (Plugin, error) {
	filename, err := resolveStarlarkFilename(workDir, filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin file %q: %w", filename, err)
	}
	defer f.Close()

	return loadStarlarkPlugin(name, filename, f, reporter, errorReporter)
}

func loadStarlarkPlugin(name, filename string, src interface{}, reporter func(msg string), errorReporter func(err error)) (Plugin, error) {

	newErrorf := func(msg string, args ...interface{}) error {
		err := fmt.Errorf(filename+": "+msg, args...)
		errorReporter(err)
		return err
	}

	plugins := make(map[string]*starlarkstruct.Struct)
	rules := make(map[string]*starlarkstruct.Struct)
	predeclared := newPredeclared(plugins, rules)

	_, thread, err := loadStarlarkProgram(filename, src, predeclared, reporter, errorReporter)
	if err != nil {
		return nil, err
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

// starlarkPlugin is a plugin implemented in starlark that implements the protoc
// plugin interface.
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
		lbl := label.NoLabel
		labelStr := labelValue.(starlark.String).GoString()
		if labelStr != "" {
			var err error
			lbl, err = label.Parse(labelStr)
			if err != nil {
				p.errorReporter("PluginConfiguration.label parse: %v", err)
				return nil
			}
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

		return starlarkstruct.FromStringDict(
			Symbol("PluginConfiguration"),
			starlark.StringDict{
				"label":   starlark.String(labelStr),
				"outputs": outputs,
				"out":     starlark.String(out),
				"options": options,
			},
		), nil
	}
}

func newStarlarkPluginFunction(plugins map[string]*starlarkstruct.Struct) goStarlarkFunction {
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
			"plugin_config":  newLanguagePluginConfigStruct(ctx.PluginConfig),
			"package_config": newPackageConfigStruct(&ctx.PackageConfig),
			"proto_library":  newProtoLibraryStruct(ctx.ProtoLibrary),
		},
	)
}

func newLanguagePluginConfigStruct(cfg LanguagePluginConfig) *starlarkstruct.Struct {
	var labelStr string
	if cfg.Label != label.NoLabel {
		labelStr = cfg.Label.String()
	}
	return starlarkstruct.FromStringDict(
		Symbol("LanguagePluginConfig"),
		starlark.StringDict{
			"name":           starlark.String(cfg.Name),
			"implementation": starlark.String(cfg.Implementation),
			"label":          starlark.String(labelStr),
			"options":        newStringList(cfg.GetOptions()),
			"deps":           newStringList(cfg.GetDeps()),
			"enabled":        starlark.Bool(cfg.Enabled),
		},
	)
}

func newPackageConfigStruct(cfg *PackageConfig) *starlarkstruct.Struct {
	if cfg == nil {
		return starlarkstruct.FromStringDict(
			Symbol("PackageConfig"),
			starlark.StringDict{
				"config": newConfigStruct(&config.Config{}),
			},
		)
	}
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
			"rule":                newStarlarkProtoLibraryRuleStruct(p.Rule()),
		},
	)
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
			"relname":      starlark.String(f.Relname()),
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

func newStarlarkProtoLibraryRuleStruct(r *rule.Rule) *starlarkstruct.Struct {
	if r == nil {
		return starlarkstruct.FromStringDict(
			Symbol("Rule"),
			starlark.StringDict{
				"name":       starlark.String(""),
				"kind":       starlark.String(""),
				"srcs":       &starlark.List{},
				"deps":       &starlark.List{},
				"tags":       &starlark.List{},
				"visibility": &starlark.List{},
			},
		)
	}
	return starlarkstruct.FromStringDict(
		Symbol("Rule"),
		starlark.StringDict{
			"name":       starlark.String(r.Name()),
			"kind":       starlark.String(r.Kind()),
			"srcs":       newStringList(r.AttrStrings("srcs")),
			"deps":       newStringList(r.AttrStrings("deps")),
			"tags":       newStringList(r.AttrStrings("tags")),
			"visibility": newStringList(r.AttrStrings("visibility")),
		},
	)

}
