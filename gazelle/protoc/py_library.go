package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/buildtools/build"
)

const rosettaPipDepsRequirementsLabel = "@rosetta_pip_deps//:requirements.bzl"
const rosettaPipRequirementToken = "requirement"

// we may want this capability in the future, so leaving it under a compile-time
// flag.
const addPipRequirements = false

// PyLibrary implements a ruleProvider for the "py_library" rule. This is a thin
// wrapper around a py_library that contains generated py proto sources, but
// does not interfere with the python gazelle extension.  There is one rosetta
// py_library per proto file.
type PyLibrary struct {
	Rel   string // relative path from WORKSPACE root
	Lib   ProtoLibrary
	Rules []*ProtoRule
}

// Kind implements part of the ruleProvider interface
func (s *PyLibrary) Kind() string {
	return "py_library"
}

// KindInfo implements part of the ruleProvider interface
func (s *PyLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		NonEmptyAttrs:  map[string]bool{"srcs": true, "deps": true},
		MergeableAttrs: map[string]bool{"srcs": true, "deps": true},
	}
}

// Rule implements part of the ruleProvider interface
func (s *PyLibrary) Rule() *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())
	newRule.SetAttr("srcs", s.Srcs())
	newRule.SetAttr("imports", s.getImports())
	if s.Visibility() != nil {
		newRule.SetAttr("visibility", s.Visibility())
	}

	return newRule
}

// Deps computes rule dependencies based on the proto_library rule.  Well-known
// types are excluded.
func (s *PyLibrary) Deps() []build.Expr {
	deps := make([]build.Expr, 0)

	for _, dep := range s.Lib.Deps() {
		if strings.HasPrefix(dep, "@com_google_protobuf") {
			continue
		}
		prefix := strings.TrimSuffix(dep, "_proto")
		deps = append(deps, stringExpr(prefix+"_py_library"))
	}

	if addPipRequirements {
		var needsProtobufRequirement, needsGrpcioRequirement bool
		for _, src := range s.Srcs() {
			if strings.HasSuffix(src, "_pb2.py") {
				needsProtobufRequirement = true
			} else if strings.HasSuffix(src, "_pb2_grpc.py") {
				needsGrpcioRequirement = true
			}
		}

		if needsProtobufRequirement {
			deps = append(deps, newRequirement("protobuf"))
		}
		if needsGrpcioRequirement {
			deps = append(deps, newRequirement("grpcio"))
		}
	}

	return deps
}

// Imports implements part of the ruleProvider interface
func (s *PyLibrary) Imports() []string {
	imports := []string{s.Kind()}

	return imports
}

// Name implements part of the ruleProvider interface
func (s *PyLibrary) Name() string {
	return s.Lib.BaseName() + "_py_library"
}

// Srcs returns the source files
func (s *PyLibrary) Srcs() []string {
	srcs := make([]string, 0)
	for _, rule := range s.Rules {
		srcs = append(srcs, rule.GeneratedSrcs()...)
	}
	return srcs
}

// Visibility returns the visibility constraints
func (s *PyLibrary) Visibility() []string {
	// defaulting to no visibility for now
	module := ""
	parts := strings.Split(s.Rel, "/")
	if len(parts) > 0 {
		module = parts[0]
	}
	if module == "" {
		return nil
	}
	return []string{"//" + module + ":__subpackages__"}
}

func (s *PyLibrary) getImports() []string {
	parts := strings.Split(s.Rel, "/")
	for i := range parts {
		parts[i] = ".."
	}
	return []string{path.Join(parts...)}
}

// VisitFile implements the FileVisitor interface
func (s *PyLibrary) VisitFile(file *rule.File) *rule.File {
	if !addPipRequirements {
		return file
	}

	var load *rule.Load

	for _, dep := range s.Deps() {
		if callExpr, ok := dep.(*build.CallExpr); ok {
			if name, ok := callExpr.X.(*build.LiteralExpr); ok {
				if name.Token == rosettaPipRequirementToken {
					load = rule.NewLoad(rosettaPipDepsRequirementsLabel)
					load.Add(rosettaPipRequirementToken)
					break
				}
			}
		}
	}

	if load != nil && file != nil {
		var hasLoad bool
		for _, l := range file.Loads {
			if l.Name() == rosettaPipDepsRequirementsLabel {
				hasLoad = true
			}
		}
		if !hasLoad {
			load.Insert(file, 0)
		}
	}

	return file
}

// Resolve implements part of the RuleProvider interface.
func (s *PyLibrary) Resolve(c *config.Config, r *rule.Rule, importsRaw interface{}, from label.Label) {
	deps := s.Deps()
	if len(deps) > 0 {
		r.SetAttr("deps", deps)
	}
}

func newRequirement(name string) *build.CallExpr {
	return &build.CallExpr{
		X: &build.LiteralExpr{Token: rosettaPipRequirementToken},
		List: []build.Expr{
			&build.StringExpr{Value: name},
		},
	}
}

func stringExpr(s string) *build.StringExpr {
	return &build.StringExpr{Value: s}
}
