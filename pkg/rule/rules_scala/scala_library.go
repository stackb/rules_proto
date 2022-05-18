package rules_scala

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bmatcuk/doublestar"
	"github.com/emicklei/proto"

	"github.com/stackb/rules_proto/pkg/plugin/akka/akka_grpc"
	"github.com/stackb/rules_proto/pkg/plugin/scalapb/scalapb"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	GrpcscalaLibraryRuleName        = "grpc_scala_library"
	ProtoscalaLibraryRuleName       = "proto_scala_library"
	scalaLibraryRuleSuffix          = "_scala_library"
	scalaPbPluginOptionsPrivateKey  = "_scalapb_plugin"
	akkaGrpcPluginOptionsPrivateKey = "_akka_grpc_plugin"
	scalapbOptionsName              = "(scalapb.options)"
	scalapbFieldTypeName            = "(scalapb.field).type"
	scalaLangName                   = "scala"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:"+ProtoscalaLibraryRuleName,
		&scalaLibrary{
			kindName: ProtoscalaLibraryRuleName,
			shouldProvideRule: func(library protoc.ProtoLibrary, plugin *protoc.PluginConfiguration) bool {
				return !hasServicesAndGrpcOption(library, plugin)
			},
		})
	protoc.Rules().MustRegisterRule("stackb:rules_proto:"+GrpcscalaLibraryRuleName,
		&scalaLibrary{
			kindName:          GrpcscalaLibraryRuleName,
			shouldProvideRule: hasServicesAndGrpcOption,
		})
}

func hasServicesAndGrpcOption(library protoc.ProtoLibrary, plugin *protoc.PluginConfiguration) bool {
	// if any of the proto_library files have grpc service definitions AND the
	// grpc option is configured, emit a grpc_scala_library rule instead.
	if !protoc.HasServices(library.Files()...) {
		return false
	}
	for option, want := range plugin.Config.Options {
		if option == "grpc" && want {
			return true
		}
	}
	return false
}

// scalaLibrary implements LanguageRule for the 'proto_scala_library' rule from
// @rules_proto.
type scalaLibrary struct {
	kindName          string
	shouldProvideRule func(library protoc.ProtoLibrary, plugin *protoc.PluginConfiguration) bool
}

// Name implements part of the LanguageRule interface.
func (s *scalaLibrary) Name() string {
	return s.kindName
}

// KindInfo implements part of the LanguageRule interface.
func (s *scalaLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MergeableAttrs: map[string]bool{
			"srcs": true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *scalaLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    fmt.Sprintf("@build_stack_rules_proto//rules/scala:%s.bzl", s.kindName),
		Symbols: []string{s.kindName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *scalaLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	options := parseScalaLibraryOptions(s.kindName, cfg.GetOptions())

	//
	// output preparation
	//
	plugin := pc.GetPluginConfiguration(scalapb.ScalaPBPluginName)
	if plugin == nil {
		log.Fatalf("expected plugin configuration for %q to be defined", scalapb.ScalaPBPluginName)
	}
	if len(plugin.Outputs) == 0 {
		return nil
	}
	if !s.shouldProvideRule(pc.Library, plugin) {
		return nil
	}

	outputs := append(plugin.Outputs, pc.GetPluginOutputs(akka_grpc.AkkaGrpcPluginName)...)
	outputs = options.filterOutputs(outputs)
	if len(outputs) == 0 {
		return nil
	}

	return &scalaLibraryRule{
		kindName:       s.kindName,
		ruleNameSuffix: scalaLibraryRuleSuffix,
		options:        options,
		outputs:        outputs,
		ruleConfig:     cfg,
		config:         pc,
	}
}

// scalaLibraryRule implements RuleProvider for 'scala_library'-derived rules.
type scalaLibraryRule struct {
	kindName       string
	ruleNameSuffix string
	outputs        []string
	config         *protoc.ProtocConfiguration
	ruleConfig     *protoc.LanguageRuleConfig
	options        *scalaLibraryOptions
}

// Kind implements part of the ruleProvider interface.
func (s *scalaLibraryRule) Kind() string {
	return s.kindName
}

// Name implements part of the ruleProvider interface.
func (s *scalaLibraryRule) Name() string {
	return s.config.Library.BaseName() + s.ruleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *scalaLibraryRule) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.outputs {
		if strings.HasSuffix(output, ".srcjar") {
			srcs = append(srcs, protoc.StripRel(s.config.Rel, output))
		}
	}
	return srcs
}

// Deps computes the deps list for the rule.
func (s *scalaLibraryRule) Deps() []string {
	deps := s.ruleConfig.GetDeps()

	for _, pluginConfig := range s.config.Plugins {
		deps = append(deps, pluginConfig.Config.GetDeps()...)
	}

	return protoc.DeduplicateAndSort(deps)
}

// Visibility provides visibility labels.
func (s *scalaLibraryRule) Visibility() []string {
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

// Rule implements part of the ruleProvider interface.
func (s *scalaLibraryRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())

	deps := s.Deps()
	if len(deps) > 0 {
		newRule.SetAttr("deps", deps)
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	// add any imports from proto options.  Example:
	// option (scalapb.options) = {
	// 	import: "com.foo.Bar"
	// };
	scalaImports := getScalapbImports(s.config.Library.Files())
	if len(scalaImports) > 0 {
		newRule.SetPrivateAttr(config.GazelleImportsKey, scalaImports)
	}

	// set the override language such that deps of 'proto_scala_library' and
	// 'grpc_scala_library' can resolve together (matches the value used by
	// "Imports").
	newRule.SetPrivateAttr(protoc.ResolverImpLangPrivateKey, "scala")

	// add the scalapb plugin options as a private attr so we can inspect them
	// during the .Imports() phase.  For example, akka 'server_power_apis'
	// generates additional classes.
	scalaPbPlugin := s.config.GetPluginConfiguration(scalapb.ScalaPBPluginName)
	if scalaPbPlugin != nil {
		newRule.SetPrivateAttr(scalaPbPluginOptionsPrivateKey, scalaPbPlugin.Options)
	}
	akkaGrpcPlugin := s.config.GetPluginConfiguration(akka_grpc.AkkaGrpcPluginName)
	if akkaGrpcPlugin != nil {
		newRule.SetPrivateAttr(akkaGrpcPluginOptionsPrivateKey, akkaGrpcPlugin.Options)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *scalaLibraryRule) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	// 1. provide generated scala class names for message and services for
	// 'scala scala' deps.  This will allow a scala extension to resolve proto
	// deps when they import scala proto class names.
	pluginOptions := make(map[string]bool)
	if scalaPbPluginOptions, ok := r.PrivateAttr(scalaPbPluginOptionsPrivateKey).([]string); ok {
		for _, opt := range scalaPbPluginOptions {
			pluginOptions[opt] = true
		}
	}
	if akkaGrpcPluginOptions, ok := r.PrivateAttr(akkaGrpcPluginOptionsPrivateKey).([]string); ok {
		for _, opt := range akkaGrpcPluginOptions {
			pluginOptions[opt] = true
		}
	}
	from := label.New("", file.Pkg, r.Name())
	provideScalaImports(s.config.Library.Files(), protoc.GlobalResolver(), from, pluginOptions)

	// 2. create import specs for 'protobuf scala'.  This allows
	// proto_scala_library and grpc_scala_library to resolve deps.
	return protoc.ProtoLibraryImportSpecsForKind("scala", s.config.Library)
}

// Resolve implements part of the RuleProvider interface.
func (s *scalaLibraryRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	imports = s.options.filterImports(imports)

	resolveFn := protoc.ResolveDepsAttr("deps", true)
	resolveFn(c, ix, r, imports, from)

	if unresolvedDeps, ok := r.PrivateAttr(protoc.UnresolvedDepsPrivateKey).(map[string]error); ok {
		resolveScalaDeps(c, ix, r, unresolvedDeps, from)

		for imp, err := range unresolvedDeps {
			if err == nil {
				continue
			}
			log.Printf("%[1]v (%[2]s): warning: failed to resolve %[3]q: %v", from, r.Kind(), imp, err)
		}
	}
}

// resolveScalaDeps attempts to resolve labels for the given deps under the
// "scala" language.  Only unresolved deps of type ErrNoLabel are considered.
// Typically these unresolved dependencies arise from (scalapb.options) imports.
func resolveScalaDeps(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, unresolvedDeps map[string]error, from label.Label) {
	lang := "scala"

	resolvedDeps := make([]string, 0)
	for imp, err := range unresolvedDeps {
		if err != protoc.ErrNoLabel {
			continue
		}
		importSpec := resolve.ImportSpec{Lang: lang, Imp: imp}
		if l, ok := resolve.FindRuleWithOverride(c, importSpec, lang); ok {
			resolvedDeps = append(resolvedDeps, l.String())
			unresolvedDeps[imp] = nil
			continue
		}
		result := ix.FindRulesByImportWithConfig(c, importSpec, lang)
		if len(result) == 0 {
			continue
		}
		if len(result) > 1 {
			log.Println(from, "multiple rules matched for scala import %q: %v", imp, result)
			continue
		}
		resolvedDeps = append(resolvedDeps, result[0].Label.String())
		unresolvedDeps[imp] = nil
	}
	if len(resolvedDeps) > 0 {
		r.SetAttr("deps", protoc.DeduplicateAndSort(append(r.AttrStrings("deps"), resolvedDeps...)))
	}
}

func getScalapbImports(files []*protoc.File) []string {
	imps := make([]string, 0)

	for _, file := range files {
		for _, option := range file.Options() {
			if option.Name != scalapbOptionsName {
				continue
			}
			for _, constant := range option.AggregatedConstants {
				switch constant.Name {
				case "import":
					if constant.Source != "" {
						imps = append(imps, constant.Source)
					}
				}
			}
		}
		for _, msg := range file.Messages() {
			for _, child := range msg.Elements {
				if field, ok := child.(*proto.NormalField); ok {
					for _, option := range field.Options {
						if option.Name != scalapbFieldTypeName {
							continue
						}
						if option.Constant.Source != "" {
							imps = append(imps, option.Constant.Source)
						}
					}
				}
			}
		}
	}

	return protoc.DeduplicateAndSort(imps)
}

// javaPackageOption is a utility function to seek for the java_package option.
func javaPackageOption(options []proto.Option) (string, bool) {
	for _, opt := range options {
		if opt.Name != "java_package" {
			continue
		}
		return opt.Constant.Source, true
	}

	return "", false
}

func provideScalaImports(files []*protoc.File, resolver protoc.ImportResolver, from label.Label, options map[string]bool) {
	lang := "scala"

	for _, file := range files {
		pkgName := file.Package().Name
		if javaPackageName, ok := javaPackageOption(file.Options()); ok {
			pkgName = javaPackageName
		}
		if pkgName != "" {
			resolver.Provide(lang, lang, pkgName, from)
		}
		for _, e := range file.Enums() {
			name := e.Name
			if pkgName != "" {
				name = pkgName + "." + name
			}
			resolver.Provide(lang, lang, name, from)
			for _, value := range e.Elements {
				if field, ok := value.(*proto.EnumField); ok {
					fieldName := name + "." + field.Name
					resolver.Provide(lang, lang, fieldName, from)
				}
			}
		}
		for _, m := range file.Messages() {
			name := m.Name
			if pkgName != "" {
				name = pkgName + "." + name
			}
			resolver.Provide(lang, lang, name, from)
			resolver.Provide(lang, lang, name+"Proto", from)
		}
		for _, s := range file.Services() {
			name := s.Name
			if pkgName != "" {
				name = pkgName + "." + name
			}
			resolver.Provide(lang, lang, name, from)
			resolver.Provide(lang, lang, name+"Grpc", from)
			resolver.Provide(lang, lang, name+"Proto", from)
			resolver.Provide(lang, lang, name+"Client", from)
			resolver.Provide(lang, lang, name+"Handler", from)
			resolver.Provide(lang, lang, name+"Server", from)
			// TOOD: if this is configured on the proto_plugin, we won't know
			// about the plugin option.  Advertise them anyway.
			// if options["server_power_apis"] {
			resolver.Provide(lang, lang, name+"PowerApi", from)
			resolver.Provide(lang, lang, name+"PowerApiHandler", from)
			resolver.Provide(lang, lang, name+"ClientPowerApi", from)
			// }
		}
	}
}

// scalaLibraryOptions represents the parsed flag configuration for a scalaLibrary
type scalaLibraryOptions struct {
	noResolve map[string]bool
	noOutput  []string
}

func parseScalaLibraryOptions(kindName string, args []string) *scalaLibraryOptions {
	flags := flag.NewFlagSet(kindName, flag.ExitOnError)

	var noresolveFlagValue string
	flags.StringVar(&noresolveFlagValue, "noresolve", "", "--noresolve=<path>.proto suppresses deps resolution of <path>.proto")

	var nooutputFlagValue string
	flags.StringVar(&nooutputFlagValue, "nooutput", "", "--nooutput=<file>.proto suppresses rule output for <file>.proto.  If after removing all matching files, no outputs remain, the rule will not be emitted.")

	if err := flags.Parse(args); err != nil {
		log.Fatalf("failed to parse flags for %q: %v", kindName, err)
	}

	config := &scalaLibraryOptions{
		noResolve: make(map[string]bool),
		noOutput:  make([]string, 0),
	}
	for _, value := range strings.Split(noresolveFlagValue, ",") {
		config.noResolve[value] = true
	}
	config.noOutput = strings.Split(nooutputFlagValue, ",")

	return config
}

func (o *scalaLibraryOptions) filterOutputs(in []string) (out []string) {
next:
	for _, value := range in {
		for _, pattern := range o.noOutput {
			match, err := doublestar.PathMatch(pattern, value)
			if err != nil {
				log.Fatalf("bad --nooutput pattern %q: %v", pattern, err)
			}
			if match {
				continue next
			}
		}
		out = append(out, value)
	}
	return
}

func (o *scalaLibraryOptions) filterImports(in []string) (out []string) {
	for _, value := range in {
		if o.noResolve[value] {
			continue
		}
		out = append(out, value)
	}
	return
}
