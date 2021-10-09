package rules_go

import (
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"

	// langgo "github.com/bazelbuild/bazel-gazelle/language/go"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoGoLibraryRuleName = "proto_go_library"
	goLibraryRuleSuffix    = "_go_proto"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:"+ProtoGoLibraryRuleName,
		&goLibrary{
			kindName: ProtoGoLibraryRuleName,
		})
}

// goLibrary implements LanguageRule for the '{proto|grpc}_go_library' rule from
// @rules_proto.
type goLibrary struct {
	kindName string
}

// Name implements part of the LanguageRule interface.
func (s *goLibrary) Name() string {
	return s.kindName
}

// KindInfo implements part of the LanguageRule interface.
func (s *goLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs":       true,
			"deps":       true,
			"visibility": true,
		},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *goLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    fmt.Sprintf("@build_stack_rules_proto//rules/go:%s.bzl", s.kindName),
		Symbols: []string{s.kindName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *goLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	// collect the outputs and the deps.  Search all the PluginConfigurations.
	// If the produced .go files, include them and add their deps.
	outputs := make([]string, 0)
	deps := make([]string, 0)

	for _, pluginConfig := range pc.Plugins {
		for _, out := range pluginConfig.Outputs {
			if path.Ext(out) == ".go" {
				outputs = append(outputs, out)
				deps = append(deps, pluginConfig.Config.GetDeps()...)
			}
		}
	}

	if len(outputs) == 0 {
		return nil
	}

	for i, output := range outputs {
		outputs[i] = path.Join(pc.Rel, path.Base(output))
	}

	return &goLibraryRule{
		kindName:       s.kindName,
		ruleNameSuffix: goLibraryRuleSuffix,
		outputs:        protoc.DeduplicateAndSort(outputs),
		deps:           protoc.DeduplicateAndSort(deps),
		ruleConfig:     cfg,
		config:         pc,
	}
}

// goLibraryRule implements RuleProvider for 'go_library'-derived rules.
type goLibraryRule struct {
	kindName       string
	ruleNameSuffix string
	outputs        []string
	deps           []string
	config         *protoc.ProtocConfiguration
	ruleConfig     *protoc.LanguageRuleConfig
}

// Kind implements part of the ruleProvider interface.
func (s *goLibraryRule) Kind() string {
	return s.kindName
}

// Name implements part of the ruleProvider interface.
func (s *goLibraryRule) Name() string {
	return s.config.Library.BaseName() + s.ruleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *goLibraryRule) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.outputs {
		srcs = append(srcs, protoc.StripRel(s.config.Rel, output))
	}
	return srcs
}

// Deps implements the protoc.DepsProvider interface.
func (s *goLibraryRule) Deps() []string {
	deps := s.ruleConfig.GetDeps()
	deps = append(deps, s.deps...)
	resolvedDeps := protoc.ResolveLibraryRewrites(s.ruleConfig.GetRewrites(), s.config.Library)
	deps = append(deps, resolvedDeps...)
	return protoc.DeduplicateAndSort(deps)
}

// Visibility implements part of the ruleProvider interface.
func (s *goLibraryRule) Visibility() []string {
	visibility := make([]string, 0)
	for k, want := range s.ruleConfig.Visibility {
		if !want {
			continue
		}
		visibility = append(visibility, k)
	}
	sort.Strings(visibility)
	return visibility
}

// ImportPath computes the import path.
func (s *goLibraryRule) ImportPath() string {
	return s.getGoPackageOption()
}

// Rule implements part of the ruleProvider interface.
func (s *goLibraryRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	importpath := s.ImportPath()
	if importpath != "" {
		newRule.SetAttr("importpath", importpath)
	}

	newRule.SetAttr("srcs", s.Srcs())

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Resolve implements part of the RuleProvider interface.
func (s *goLibraryRule) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {

	// collect a list of dependencies, then partition them into 'embeds' if
	// another go library is in the same package.
	all := s.Deps()

	for _, d := range s.config.Library.Deps() {
		if strings.HasPrefix(d, "@com_google_protobuf//") {
			continue
		}
		d = strings.TrimSuffix(d, "_proto")
		all = append(all, d+goLibraryRuleSuffix)
	}

	deps := make([]string, 0)
	embeds := make([]string, 0)

	for _, dep := range all {
		deps = append(deps, dep)

		// TOOD: do we really need to populate embeds?
		// l, err := label.Parse(dep)
		// if err != nil {
		// 	log.Fatalf("resolve deps failed for for rule %s %s: label parse %q error: %v", r.Kind(), r.Name(), dep, err)
		// }

		// if l.Pkg == "" { // "this package"
		// 	embeds = append(embeds, dep)
		// } else {
		// deps = append(deps, dep)
		// }
	}

	if len(embeds) > 0 {
		r.SetAttr("embed", embeds)
	}
	if len(deps) > 0 {
		r.SetAttr("deps", deps)
	}
}

func (s *goLibraryRule) getGoPackageOption() string {
	for _, file := range s.config.Library.Files() {
		if value, _ := protoc.GetNamedOption(file.Options(), "go_package"); value != "" {
			if strings.LastIndexByte(value, '/') == -1 {
				// return langgo.InferImportPath(c, rel)
				continue // TODO: do more research here on if this is the correct approach
			}
			if i := strings.LastIndexByte(value, ';'); i != -1 {
				return value[:i]
			}
			return value
		}
	}

	// fallback
	return s.getPluginImportMappingOption()
}

func (s *goLibraryRule) getPluginImportMappingOption() string {
	// first, iterate all the plugins and gather options that look like
	// protoc-gen-go "importmapping" (M) options (e.g
	// Mfoo.proto=github.com/example/foo).
	mappings := make(map[string]string)

	tryParseMapping := func(opt string) {
		if !strings.HasPrefix(opt, "M") {
			return
		}
		parts := strings.SplitN(opt[1:], "=", 2)
		if len(parts) != 2 {
			return
		}
		mappings[parts[0]] = parts[1]
	}

	// search all plugins
	for _, plugin := range s.config.Plugins {
		for _, opt := range plugin.Options {
			tryParseMapping(opt)
		}
	}
	// and all rule options
	for _, opt := range s.ruleConfig.GetOptions() {
		tryParseMapping(opt)
	}

	// now that we've gathered all possible options; search all library files
	// (e.g. foo.proto) and see if we can find a match.
	for _, file := range s.config.Library.Files() {
		filename := path.Join(file.Dir, file.Basename)
		mapping := mappings[filename]
		if mapping != "" {
			return mapping
		}
	}

	log.Println(s.kindName, s.Name(), "unable to determine an importpath.  Try adding a go_package option to the .proto file or use an 'M' importmapping on the plugin -OR- rule (see documentation for protoc-gen-go for syntax)")

	return ""
}
