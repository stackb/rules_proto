package protobuf

import (
	"log"
	"path"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"

	"github.com/stackb/rules_proto/pkg/protoc"
)

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested in
// depth-first post-order.
//
// args contains the arguments for GenerateRules. This is passed as a struct to
// avoid breaking implementations in the future when new fields are added.
//
// A GenerateResult struct is returned. Optional fields may be added to this
// type in the future.
//
// Any non-fatal errors this function encounters should be logged using
// log.Print.
func (pl *protobufLang) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	cfg := pl.getOrCreatePackageConfig(args.Config)

	files := make(map[string]*protoc.File)
	for _, f := range args.RegularFiles {
		if !protoc.IsProtoFile(f) {
			continue
		}
		file := protoc.NewFile(args.Rel, f)
		if err := file.Parse(); err != nil {
			log.Printf("warning: unparseable proto file dir=%s, file=%s: %v", args.Dir, file.Basename, err)
			continue
		}
		files[f] = file

		// Record the list of dependencies for this proto file.  Dependents are
		// encoded as labels as a matter of practicality given the API of the
		// resolver.
		for _, imp := range file.Imports() {
			dir := path.Dir(imp.Filename)
			if dir == "." {
				dir = ""
			}
			pl.resolver.Provide(
				"proto",
				"depends",
				path.Join(file.Dir, file.Basename),
				label.New("", dir, path.Base(imp.Filename)),
			)
		}
	}

	protoLibraries := make([]protoc.ProtoLibrary, 0)
	for _, r := range args.OtherGen {
		internalLabel := label.New(args.Config.RepoName, args.Rel, r.Name())
		protoc.GlobalRuleIndex().Put(internalLabel, r)

		if r.Kind() != "proto_library" {
			continue
		}

		srcs := r.AttrStrings("srcs")
		srcLabels := make([]label.Label, len(srcs))
		for i, src := range srcs {
			srcLabel, err := label.Parse(src)
			if err != nil {
				log.Fatalf("%s %q: unparseable source label %q: %v", r.Kind(), r.Name(), src, err)
			}
			srcLabels[i] = srcLabel

			// record the label that "provides" each proto file.
			pl.resolver.Provide(
				"proto",
				"proto",
				path.Join(args.Rel, src),
				internalLabel,
			)
		}

		lib := protoc.NewOtherProtoLibrary(args.File, r, matchingFiles(files, srcLabels)...)
		protoLibraries = append(protoLibraries, lib)
	}

	pkg := protoc.NewPackage(args.Rel, cfg, protoLibraries...)
	pl.packages[args.Rel] = pkg

	rules := pkg.Rules()

	// special case if we want to override go_googleapis deps.
	if pl.overrideGoGooleapis && len(protoLibraries) > 0 {
		rules = append(rules, makeProtoOverrideRule(protoLibraries))
	}

	imports := make([]interface{}, len(rules))
	for i, r := range rules {
		imports[i] = r.PrivateAttr(config.GazelleImportsKey)
		internalLabel := label.New("", args.Rel, r.Name())
		protoc.GlobalRuleIndex().Put(internalLabel, r)
	}

	// special case if this is the root BUILD file and the user requested to
	// write the imports file.
	if args.Rel == "" && pl.importsOutFile != "" {
		if err := protoc.GlobalResolver().SaveFile(pl.importsOutFile, pl.repoName); err != nil {
			log.Printf("error saving import file: %s: %v", pl.importsOutFile, err)
		}
	}

	return language.GenerateResult{
		Gen:     rules,
		Imports: imports,
		Empty:   pkg.Empty(),
	}
}

func matchingFiles(files map[string]*protoc.File, srcs []label.Label) []*protoc.File {
	matching := make([]*protoc.File, 0)
	for _, src := range srcs {
		if file, ok := files[src.Name]; ok {
			matching = append(matching, file)
		}
	}
	return matching
}
