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
	pluginDeps := make([]string, 0)

	for _, pluginConfig := range pc.Plugins {
		for _, out := range pluginConfig.Outputs {
			if path.Ext(out) == ".go" {
				outputs = append(outputs, out)
				pluginDeps = append(pluginDeps, pluginConfig.Config.GetDeps()...)
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
		deps:           protoc.DeduplicateAndSort(pluginDeps),
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

// deps returns all known "configured" dependencies:
// 1. Those given by the plugin implementations that contributed outputs (their 'deps' directive).
// 2. Those given by 'deps' directive on the rule config.
// 3. Those given by resolving "rewrite" specs against the proto file imports.
func (s *goLibraryRule) configDeps() []string {
	deps := s.deps
	deps = append(deps, s.ruleConfig.GetDeps()...)
	resolvedDeps := protoc.ResolveLibraryRewrites(s.ruleConfig.GetRewrites(), s.config.Library)
	deps = append(deps, resolvedDeps...)
	return deps
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

// importPath computes the import path.
func (s *goLibraryRule) importPath() string {
	// Try 'M' options first
	if imp := s.getPluginImportMappingOption(); imp != "" {
		return imp
	}
	// Fallback to the 'go_package' option in an imported file
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

	log.Printf("warning: unknown 'importpath' for %s rule //%s:%s.  Try adding the 'go_package' option to the .proto file or use an 'M' importmap option",
		s.kindName,
		s.config.Rel,
		s.Name(),
	)

	// fallback
	return ""
}

// Rule implements part of the ruleProvider interface.
func (s *goLibraryRule) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	importpath := s.importPath()
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
	// another go library has the same importpath.
	all := s.configDeps()

	// log.Println("proto_go_library rewrites", s.ruleConfig.GetRewrites())
	// log.Println("proto_go_library resolve deps", all)
	// populated a set to check for duplicates
	// seen := make(map[string]bool)

	// for _, v := range all {
	// 	seen[v] = true
	// }

	// for _, d := range s.config.Library.Deps() {
	// 	if strings.HasPrefix(d, "@com_google_protobuf//") {
	// 		continue
	// 	}
	// 	d = strings.TrimSuffix(d, "_proto")
	// 	name := d + goLibraryRuleSuffix
	// 	l := label.New(c.RepoName, s.config.Rel, name)
	// 	if !seen[l.String()] {
	// 		all = append(all, l.String())
	// 	}
	// }

	// deps := make([]string, 0)
	// embeds := make([]string, 0)

	// for _, dep := range all {
	// 	deps = append(deps, dep)

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
	// }

	// if len(embeds) > 0 {
	// 	r.SetAttr("embed", embeds)
	// }
	if len(all) > 0 {
		r.SetAttr("deps", protoc.DeduplicateAndSort(all))
	}
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

	return ""
}
