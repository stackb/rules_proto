package protobuf

import (
	"container/list"
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

// TransitiveImportMappingsKey stores a map[string]string on the library
const TransitiveImportMappingsKey = "_transitive_importmappings"

const ProtocGenGoPluginName = "golang:protobuf:protoc-gen-go"

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenGoPlugin{})
}

// ProtocGenGoPlugin implements Plugin for the the gogo_* family of plugins.
type ProtocGenGoPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenGoPlugin) Name() string {
	return ProtocGenGoPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenGoPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}
	mappings, _ := GetImportMappings(ctx.PluginConfig.GetOptions())

	// record M associations now.
	//
	// TODO(pcj): where and when is the optimal time to do this?  protoc-gen-go,
	// protoc-gen-gogo, and protoc-gen-go-grpc all use this.  Perhaps they
	// should *all* perform it, just to be sure?
	for k, v := range mappings {
		// "option" is used as the name since we cannot leave that part of the
		// label empty.
		protoc.GlobalResolver().Provide("proto", "M", k, label.New("", v, "option")) // FIXME(pcj): should this not be config.RepoName?
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/golang/protobuf", "protoc-gen-go"),
		Outputs: p.outputs(ctx.ProtoLibrary, mappings),
		Options: ctx.PluginConfig.GetOptions(),
	}
}

func (p *ProtocGenGoPlugin) ResolvePluginOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	return ResolvePluginOptionsTransitive(cfg, r, from)
}

func ResolvePluginOptionsTransitive(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	transitiveMappings := ResolveTransitiveImportMappings(r, from)

	options := make([]string, 0)

	for _, opt := range cfg.Options {
		if !strings.HasPrefix(opt, "M") {
			options = append(options, opt)
			continue
		}

		parts := strings.SplitN(opt[1:], "=", 2)
		if len(parts) != 2 {
			options = append(options, opt)
			continue
		}

		imp := parts[0]
		if _, ok := transitiveMappings[imp]; ok {
			options = append(options, opt)
			continue
		}

		// if we get here, the M option is not in the set of transitives for
		// this rule, so leave it out.
	}

	return options
}

func ResolveTransitiveImportMappings(r *rule.Rule, from label.Label) map[string]string {
	lib := r.PrivateAttr(protoc.ProtoLibraryKey)
	if lib == nil {
		return nil
	}
	library := lib.(protoc.ProtoLibrary)
	libRule := library.Rule()

	// already created?
	if transitiveMappings, ok := libRule.PrivateAttr(TransitiveImportMappingsKey).(map[string]string); ok {
		return transitiveMappings
	}

	// nope.
	transitiveMappings := make(map[string]string)
	resolver := protoc.GlobalResolver()

	seen := make(map[string]bool)
	stack := list.New()
	for _, src := range library.Srcs() {
		stack.PushBack(path.Join(from.Pkg, src))
	}
	// for every source file in the proto library, gather the list of source
	// files on which it depends, until there are no more unprocessed sources.
	// Foreach one check if there is an importmapping for it and record the
	// association.
	for {
		if stack.Len() == 0 {
			break
		}
		current := stack.Front()
		stack.Remove(current)

		protofile := current.Value.(string)
		if seen[protofile] {
			continue
		}
		seen[protofile] = true

		depends := resolver.Resolve("proto", "depends", protofile)
		for _, dep := range depends {
			stack.PushBack(path.Join(dep.Label.Pkg, dep.Label.Name))
		}

		mappings := resolver.Resolve("proto", "M", protofile)
		if len(mappings) > 0 {
			first := mappings[0]
			transitiveMappings[protofile] = path.Join(first.Label.Pkg)
		}
	}

	libRule.SetPrivateAttr(TransitiveImportMappingsKey, transitiveMappings)

	return transitiveMappings
}

func (p *ProtocGenGoPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

func (p *ProtocGenGoPlugin) outputs(lib protoc.ProtoLibrary, importMappings map[string]string) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !(f.HasMessages() || f.HasEnums()) {
			continue
		}
		srcs = append(srcs, GetGoOutputBaseName(f, importMappings)+".pb.go")
	}
	return srcs
}

func GetGoOutputBaseName(f *protoc.File, importMappings map[string]string) string {
	base := f.Name
	pkg := f.Package()
	// see https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L347
	if mapping := importMappings[path.Join(f.Dir, f.Basename)]; mapping != "" {
		base = path.Join(mapping, base)
	} else if goPackage, _, ok := protoc.GoPackageOption(f.Options()); ok {
		base = path.Join(goPackage, base)
	} else if pkg.Name != "" {
		base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
	}
	return base
}

func GetImportMappings(options []string) (map[string]string, []string) {
	// gather options that look like protoc-gen-go "importmapping" (M) options
	// (e.g Mfoo.proto=github.com/example/foo).
	mappings := make(map[string]string)
	rest := make([]string, 0)

	for _, opt := range options {
		if !strings.HasPrefix(opt, "M") {
			rest = append(rest, opt)
			continue
		}
		parts := strings.SplitN(opt[1:], "=", 2)
		if len(parts) != 2 {
			rest = append(rest, opt)
			continue
		}
		mappings[parts[0]] = parts[1]
	}

	return mappings, rest
}
