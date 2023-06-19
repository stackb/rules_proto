package protoc

import (
	"fmt"
	"log"
	"os"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func LoadStarlarkLanguageRuleFromFile(workDir, filename, name string, reporter func(msg string), errorReporter func(err error)) (LanguageRule, error) {
	filename, err := resolveStarlarkFilename(workDir, filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open rule file %q: %w", filename, err)
	}
	defer f.Close()

	return loadStarlarkLanguageRule(name, filename, f, reporter, errorReporter)
}

func loadStarlarkLanguageRule(name, filename string, src interface{}, reporter func(msg string), errorReporter func(err error)) (LanguageRule, error) {
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

	if rule, ok := rules[name]; !ok {
		return nil, newErrorf("rule %q was never declared", name)
	} else {
		return &starlarkLanguageRule{
			name:          name,
			rule:          rule,
			reporter:      thread.Print,
			errorReporter: newErrorf,
		}, nil
	}
}

func newStarlarkLanguageRuleFunction(rules map[string]*starlarkstruct.Struct) goStarlarkFunction {
	return func(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
		var name string
		var provideRule starlark.Callable
		var loadInfo, kindInfo starlark.Callable

		if err := starlark.UnpackArgs("Rule", args, kwargs,
			"name", &name,
			"load_info", &loadInfo,
			"kind_info", &kindInfo,
			"provide_rule", &provideRule,
		); err != nil {
			return nil, err
		}

		rule := starlarkstruct.FromStringDict(
			Symbol("Rule"),
			starlark.StringDict{
				"name":         starlark.String(name),
				"load_info":    loadInfo,
				"kind_info":    kindInfo,
				"provide_rule": provideRule,
			},
		)

		rules[name] = rule
		return rule, nil
	}
}

// starlarkLanguageRule is a rule implemented in starlark that implements the protoc
// LanguageRule interface.
type starlarkLanguageRule struct {
	name          string
	reporter      func(thread *starlark.Thread, msg string)
	errorReporter func(msg string, args ...interface{}) error
	rule          *starlarkstruct.Struct
}

func (p *starlarkLanguageRule) Name() string {
	return p.name
}

// LoadInfo returns the gazelle LoadInfo.
func (p *starlarkLanguageRule) LoadInfo() (result rule.LoadInfo) {
	callable, err := p.rule.Attr("load_info")
	if err != nil {
		log.Fatalf("LoadInfo() called on rule %q with no load_info function: %v", p.name, p.rule)
		return result
	}

	thread := new(starlark.Thread)
	thread.Print = p.reporter
	value, err := starlark.Call(thread, callable, starlark.Tuple{}, []starlark.Tuple{})
	if err != nil {
		log.Fatalf("rule %q load_info failed: %v", p.name, err)
		return result
	}

	switch value := value.(type) {
	case *starlarkstruct.Struct:
		result.Name = structAttrString(value, "name", p.errorReporter)
		result.Symbols = structAttrStringSlice(value, "symbols", p.errorReporter)
		result.After = structAttrStringSlice(value, "after", p.errorReporter)
	default:
		p.errorReporter("rule %q provide_rule returned invalid type: %T", p.name, value)
	}
	return
}

// KindInfo returns the gazelle KindInfo.
func (p *starlarkLanguageRule) KindInfo() (result rule.KindInfo) {
	callable, err := p.rule.Attr("kind_info")
	if err != nil {
		p.errorReporter("rule %q has no kind_info function", p.name)
		return result
	}

	thread := new(starlark.Thread)
	thread.Print = p.reporter
	value, err := starlark.Call(thread, callable, starlark.Tuple{}, []starlark.Tuple{})
	if err != nil {
		p.errorReporter("rule %q kind_info failed: %w", p.name, err)
		return result
	}

	switch value := value.(type) {
	case *starlarkstruct.Struct:
		result.MatchAny = structAttrBool(value, "match_any", p.errorReporter)
		result.MatchAttrs = structAttrStringSlice(value, "match_attrs", p.errorReporter)
		result.NonEmptyAttrs = structAttrMapStringBool(value, "non_empty_attrs", p.errorReporter)
		result.SubstituteAttrs = structAttrMapStringBool(value, "substitute_attrs", p.errorReporter)
		result.MergeableAttrs = structAttrMapStringBool(value, "mergeable_attrs", p.errorReporter)
		result.ResolveAttrs = structAttrMapStringBool(value, "resolve_attrs", p.errorReporter)
	default:
		p.errorReporter("rule %q provide_rule returned invalid type: %T", p.name, value)
	}

	return
}

// ProvideRule takes the given configration and compilation and emits a
// RuleProvider.  If the state of the ProtocConfiguration is such that the
// rule should not be emitted, implementation should return nil.
func (p *starlarkLanguageRule) ProvideRule(rc *LanguageRuleConfig, pc *ProtocConfiguration) RuleProvider {

	callable, err := p.rule.Attr("provide_rule")
	if err != nil {
		p.errorReporter("rule %q has no provide_rule function", p.name)
		return nil
	}

	thread := new(starlark.Thread)
	thread.Print = p.reporter
	value, err := starlark.Call(thread, callable, starlark.Tuple{
		newLanguageRuleConfigStruct(rc),
		newProtocConfigurationStruct(pc),
	}, []starlark.Tuple{})
	if err != nil {
		p.errorReporter("rule %q provide_rule failed: %w", p.name, err)
		return nil
	}

	var result RuleProvider
	switch value := value.(type) {
	case *starlarkstruct.Struct:
		var experimentalResolveDepsAttr string
		if attr, err := value.Attr("experimental_resolve_attr"); err == nil {
			if str, ok := attr.(starlark.String); ok {
				experimentalResolveDepsAttr = str.GoString()
			}
		}
		result = &starlarkRuleProvider{
			name:                        p.name,
			provider:                    value,
			reporter:                    p.reporter,
			errorReporter:               p.errorReporter,
			experimentalResolveDepsAttr: experimentalResolveDepsAttr,
		}
	default:
		p.errorReporter("rule %q provide_rule returned invalid type: %T", p.name, value)
		return nil
	}

	return result
}

// starlarkRuleProvider implements RuleProvider via a starlark struct.
type starlarkRuleProvider struct {
	name                        string
	provider                    *starlarkstruct.Struct
	reporter                    func(thread *starlark.Thread, msg string)
	errorReporter               func(msg string, args ...interface{}) error
	experimentalResolveDepsAttr string
}

// Kind implements part of the RuleProvider interface.
func (s *starlarkRuleProvider) Kind() string {
	kind, err := s.provider.Attr("kind")
	if err != nil {
		s.errorReporter("provider %q has no kind", s.Name())
		return ""
	}
	return kind.(starlark.String).GoString()
}

// Name implements part of the RuleProvider interface.
func (s *starlarkRuleProvider) Name() string {
	return structAttrString(s.provider, "name", s.errorReporter)
}

// Rule implements part of the RuleProvider interface.
func (s *starlarkRuleProvider) Rule(othergen ...*rule.Rule) *rule.Rule {
	callable, err := s.provider.Attr("rule")
	if err != nil {
		s.errorReporter("rule %q has no rule() function", s.name)
		return nil
	}

	thread := new(starlark.Thread)
	s.reporter(thread, "Invoking rule "+s.name)
	thread.Print = s.reporter
	value, err := starlark.Call(thread, callable, starlark.Tuple{}, []starlark.Tuple{})
	if err != nil {
		s.errorReporter("provider %q rule() failed: %w", s.name, err)
		return nil
	}

	switch value := value.(type) {
	case *starlarkstruct.Struct:
		rKind := structAttrString(value, "kind", s.errorReporter)
		if rKind == "" {
			s.errorReporter("rule %q has no kind", s.name)
			return nil
		}
		rName := structAttrString(value, "name", s.errorReporter)
		if rName == "" {
			s.errorReporter("rule %q has no name", s.name)
			return nil
		}
		r := rule.NewRule(rKind, rName)
		s.reporter(thread, "Created rule "+rKind+" "+rName)

		attrs, err := value.Attr("attrs")
		if err != nil {
			s.errorReporter("provider %q rule() returned invalid type: %T", s.name, value)
			return nil
		}
		switch attrs := attrs.(type) {
		case *starlark.Dict:
			for _, attr := range attrs.Keys() {
				attrName, ok := attr.(starlark.String)
				if !ok {
					s.errorReporter("%q rule attr key is invalid type (want string, got %T)", s.name, attr)
					continue
				}
				if attrValue, ok, err := attrs.Get(attrName); ok && err == nil {
					switch t := attrValue.(type) {
					case *starlark.Bool:
						r.SetAttr(attrName.GoString(), bool(*t))
					case *starlark.Int:
						intValue, _ := t.Int64()
						r.SetAttr(attrName.GoString(), intValue)
					case starlark.String:
						r.SetAttr(attrName.GoString(), t.GoString())
					case *starlark.List:
						s.reporter(thread, fmt.Sprintf("!!! %q rule attr %q is a list", s.name, attrName.GoString()))
						r.SetAttr(attrName.GoString(), stringSlice(t, s.errorReporter))
					default:
						s.errorReporter("%q rule attr value is invalid type (want bool, int, string, list, got %T)", s.name, t)
					}
				}
			}
		default:
			s.errorReporter("%q rule.attrs returned invalid type: %T", s.name, value)
		}
		return r
	default:
		s.errorReporter("provider %q rule() returned invalid type: %T", s.name, value)
		return nil
	}
}

// Resolve implements part of the RuleProvider interface.
func (s *starlarkRuleProvider) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	if s.experimentalResolveDepsAttr != "" {
		if r.Attr(s.experimentalResolveDepsAttr) != nil {
			ResolveDepsAttr(s.experimentalResolveDepsAttr, false)(c, ix, r, imports, from)
		}
	}
}

// Imports implements part of the RuleProvider interface.
func (s *starlarkRuleProvider) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	if s.experimentalResolveDepsAttr != "" {
		if lib, ok := r.PrivateAttr(ProtoLibraryKey).(ProtoLibrary); ok {
			return ProtoLibraryImportSpecsForKind(r.Kind(), lib)
		}
	}
	return nil
}

func newLanguageRuleConfigStruct(rc *LanguageRuleConfig) *starlarkstruct.Struct {
	if rc == nil {
		return starlarkstruct.FromStringDict(
			Symbol("LanguageRuleConfig"),
			starlark.StringDict{
				"attrs":          &starlark.Dict{},
				"config":         newConfigStruct(&config.Config{}),
				"deps":           &starlark.List{},
				"enabled":        starlark.Bool(false),
				"implementation": starlark.String(""),
				"name":           starlark.String(""),
				"options":        &starlark.List{},
				"visibility":     &starlark.List{},
			},
		)
	}
	return starlarkstruct.FromStringDict(
		Symbol("LanguageRuleConfig"),
		starlark.StringDict{
			"attrs":          newStringListDict(rc.Attrs),
			"config":         newConfigStruct(rc.Config),
			"deps":           newStringList(rc.GetDeps()),
			"enabled":        starlark.Bool(rc.Enabled),
			"implementation": starlark.String(rc.Implementation),
			"name":           starlark.String(rc.Name),
			"options":        newStringList(rc.GetOptions()),
			"visibility":     newStringList(rc.GetVisibility()),
		},
	)
}

func newProtocConfigurationStruct(pc *ProtocConfiguration) *starlarkstruct.Struct {
	if pc == nil {
		return starlarkstruct.FromStringDict(
			Symbol("ProtocConfiguration"),
			starlark.StringDict{
				"package_config":  newPackageConfigStruct(nil),
				"language_config": newLanguageConfigStruct(nil),
				"library":         starlark.None,
				"rel":             starlark.String(""),
				"prefix":          starlark.String(""),
				"outputs":         newStringList([]string{}),
				"imports":         newStringList([]string{}),
				"mappings":        &starlark.Dict{},
				"plugins":         &starlark.List{},
			},
		)
	}
	return starlarkstruct.FromStringDict(
		Symbol("ProtocConfiguration"),
		starlark.StringDict{
			"package_config":  newPackageConfigStruct(pc.PackageConfig),
			"language_config": newLanguageConfigStruct(pc.LanguageConfig),
			"proto_library":   newProtoLibraryStruct(pc.Library),
			"rel":             starlark.String(pc.Rel),
			"prefix":          starlark.String(pc.Prefix),
			"outputs":         newStringList(pc.Outputs),
			"imports":         newStringList(pc.Imports),
			"mappings":        newStringStringDict(pc.Mappings),
			"plugins":         newPluginConfigurationList(pc.Plugins),
		},
	)
}

func newLanguageConfigStruct(lc *LanguageConfig) *starlarkstruct.Struct {
	if lc == nil {
		return starlarkstruct.FromStringDict(
			Symbol("LanguageConfig"),
			starlark.StringDict{
				"name":    starlark.String(""),
				"protoc":  starlark.String(""),
				"enabled": starlark.Bool(false),
				"plugins": &starlark.Dict{},
				"rules":   &starlark.Dict{},
			},
		)
	}
	return starlarkstruct.FromStringDict(
		Symbol("LanguageConfig"),
		starlark.StringDict{
			"name":    starlark.String(lc.Name),
			"protoc":  starlark.String(lc.Protoc),
			"enabled": starlark.Bool(lc.Enabled),
			"plugins": newStringBoolDict(lc.Plugins),
			"rules":   newStringBoolDict(lc.Rules),
		},
	)
}

func newPluginConfigurationList(plugins []*PluginConfiguration) *starlark.List {
	list := make([]starlark.Value, len(plugins))
	for i, p := range plugins {
		list[i] = newPluginConfigurationStruct(*p)
	}
	return starlark.NewList(list)
}

func newPluginConfigurationStruct(p PluginConfiguration) *starlarkstruct.Struct {
	return starlarkstruct.FromStringDict(
		Symbol("PluginConfiguration"),
		starlark.StringDict{
			"config":   newLanguagePluginConfigStruct(*p.Config),
			"label":    starlark.String(p.Label.String()),
			"out":      starlark.String(p.Out),
			"options":  newStringList(p.Options),
			"outputs":  newStringList(p.Outputs),
			"mappings": newStringStringDict(p.Mappings),
		},
	)
}
