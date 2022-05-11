package protoc

import (
	"fmt"
	"os"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

func isStarlarkLanguageRule(filename string) bool {
	return strings.HasSuffix(filename, ".star")
}

func loadStarlarkLanguageRuleFromFile(name, filename string, reporter func(msg string), errorReporter func(err error)) (LanguageRule, error) {
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
func (p *starlarkLanguageRule) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{}
}

// KindInfo returns the gazelle KindInfo.
func (p *starlarkLanguageRule) KindInfo() rule.KindInfo {
	return rule.KindInfo{}
}

// ProvideRule takes the given configration and compilation and emits a
// RuleProvider.  If the state of the ProtocConfiguration is such that the
// rule should not be emitted, implementation should return nil.
func (p *starlarkLanguageRule) ProvideRule(rc *LanguageRuleConfig, pc *ProtocConfiguration) RuleProvider {

	provideRule, err := p.rule.Attr("provide_rule")
	if err != nil {
		p.errorReporter("rule %q has no provide_rule function", p.name)
		return nil
	}

	thread := new(starlark.Thread)
	thread.Print = p.reporter
	value, err := starlark.Call(thread, provideRule, starlark.Tuple{
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
		result = &starlarkRuleProvider{
			provider:      value,
			reporter:      p.reporter,
			errorReporter: p.errorReporter,
		}
	default:
		p.errorReporter("rule %q provide_rule returned invalid type: %T", p.name, value)
		return nil
	}

	return result
}

// starlarkRuleProvider implements RuleProvider via a starlark struct.
type starlarkRuleProvider struct {
	provider      *starlarkstruct.Struct
	reporter      func(thread *starlark.Thread, msg string)
	errorReporter func(msg string, args ...interface{}) error
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
	name, err := s.provider.Attr("name")
	if err != nil {
		s.errorReporter("provider %T has no name", s)
		return ""
	}
	return name.(starlark.String).GoString()

}

// Rule implements part of the RuleProvider interface.
func (s *starlarkRuleProvider) Rule(othergen ...*rule.Rule) *rule.Rule {
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *starlarkRuleProvider) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {

}

// Imports implements part of the RuleProvider interface.
func (s *starlarkRuleProvider) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	return nil
}

func newLanguageRuleConfigStruct(rc *LanguageRuleConfig) *starlarkstruct.Struct {
	if rc == nil {
		return starlarkstruct.FromStringDict(
			Symbol("LanguageRuleConfig"),
			starlark.StringDict{
				"config":         newConfigStruct(&config.Config{}),
				"deps":           &starlark.Dict{},
				"enabled":        starlark.Bool(false),
				"implementation": starlark.String(""),
				"name":           starlark.String(""),
				"options":        &starlark.Dict{},
				"visibility":     &starlark.Dict{},
			},
		)
	}
	return starlarkstruct.FromStringDict(
		Symbol("LanguageRuleConfig"),
		starlark.StringDict{
			"config":         newConfigStruct(rc.Config),
			"deps":           newStringBoolDict(rc.Deps),
			"enabled":        starlark.Bool(rc.Enabled),
			"implementation": starlark.String(rc.Implementation),
			"name":           starlark.String(rc.Name),
			"options":        newStringBoolDict(rc.Options),
			"visibility":     newStringBoolDict(rc.Visibility),
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
