package rules_go

import (
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	langgo "github.com/bazelbuild/bazel-gazelle/language/go"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/plugin/grpc/grpcgo"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoGoLibraryRuleName = "proto_go_library"
	goLibraryRuleSuffix    = "_go_proto"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:"+ProtoGoLibraryRuleName,
		&goLibrary{
			pluginName: grpcgo.ProtocGenGoPluginName,
			kindName:   ProtoGoLibraryRuleName,
		})
}

// func hasServicesAndGrpcOption(library protoc.ProtoLibrary, plugin *protoc.PluginConfiguration) bool {
// 	// if any of the proto_library files have grpc service definitions AND the
// 	// grpc option is configured, emit a grpc_go_library rule instead.
// 	if !protoc.HasServices(library.Files()...) {
// 		return false
// 	}
// 	for option, want := range plugin.Config.Options {
// 		if option == "grpc" && want {
// 			return true
// 		}
// 	}
// 	return false
// }

// goLibrary implements LanguageRule for the '{proto|grpc}_go_library' rule from
// @rules_proto.
type goLibrary struct {
	pluginName string
	kindName   string
	// shouldProvideRule func(library protoc.ProtoLibrary, plugin *protoc.PluginConfiguration) bool
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
	plugin := pc.GetPluginConfiguration(s.pluginName)
	if plugin == nil {
		log.Fatalf("expected plugin configuration for %q to be defined", s.pluginName)
	}
	if len(plugin.Outputs) == 0 {
		return nil
	}

	outputs := make([]string, len(plugin.Outputs))
	for i, output := range plugin.Outputs {
		outputs[i] = path.Join(pc.Rel, path.Base(output))
	}
	return &goLibraryRule{
		kindName:       s.kindName,
		ruleNameSuffix: goLibraryRuleSuffix,
		outputs:        outputs,
		ruleConfig:     cfg,
		config:         pc,
		resolver:       protoc.ResolveDepsWithSuffix(goLibraryRuleSuffix),
	}
}

// goLibraryRule implements RuleProvider for 'go_library'-derived rules.
type goLibraryRule struct {
	kindName       string
	ruleNameSuffix string
	outputs        []string
	config         *protoc.ProtocConfiguration
	ruleConfig     *protoc.LanguageRuleConfig
	resolver       protoc.DepsResolver
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

// Deps computes the deps list for the rule.
func (s *goLibraryRule) Deps() []string {
	return s.ruleConfig.GetDeps()
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
	return getGoPackageOption(s.ruleConfig.Config, s.config.Rel, s.config.Library.Files())
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
	if s.resolver == nil {
		return
	}
	s.resolver(s, s.config, c, r, importsRaw, from)
}

func getGoPackageOption(c *config.Config, rel string, files []*protoc.File) string {
	for _, file := range files {
		if value, _ := protoc.GetNamedOption(file.Options(), "go_package"); value != "" {
			if strings.LastIndexByte(value, '/') == -1 {
				return langgo.InferImportPath(c, rel)
			}
			if i := strings.LastIndexByte(value, ';'); i != -1 {
				return value[:i]
			}
			return value
		}
	}
	return ""
}
