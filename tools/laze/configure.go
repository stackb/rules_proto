package laze

import (
	"flag"
	"log"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// Configurer satisfies the config.Configurer interface. It's the
// language-specific configuration extension.
type Configurer struct{}

// RegisterFlags registers command-line flags used by the extension. This method
// is called once with the root configuration when Gazelle starts. RegisterFlags
// may set an initial values in Config.Exts. When flags are set, they should
// modify these values.
func (py *Configurer) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	// TODO(pcj): implement.
}

// CheckFlags validates the configuration after command line flags are parsed.
// This is called once with the root configuration when Gazelle starts.
// CheckFlags may set default values in flags or make implied changes.
func (py *Configurer) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	// TODO(pcj): implement.
	return nil
}

// KnownDirectives returns a list of directive keys that this Configurer can
// interpret. Gazelle prints errors for directives that are not recoginized by
// any Configurer.
func (py *Configurer) KnownDirectives() []string {
	// TODO(pcj): implement.
	return make([]string, 0)
}

// Configure modifies the configuration using directives and other information
// extracted from a build file. Configure is called in each directory.
//
// c is the configuration for the current directory. It starts out as a copy of
// the configuration for the parent directory.
//
// rel is the slash-separated relative path from the repository root to the
// current directory. It is "" for the root directory itself.
//
// f is the build file for the current directory or nil if there is no existing
// build file.
func (py *Configurer) Configure(c *config.Config, rel string, f *rule.File) {
	// TODO(pcj): implement.
	log.Println("Configure: %+v", c)
}
