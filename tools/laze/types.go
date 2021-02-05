package laze

import "strings"

// pyModule represents a parsed Python module.
type pyModule struct {
	// The name of the Python module.
	name string
	// The path to the file containing the module. If the module is local, this
	// field is set. If not set, the module is either builtin or comes from a
	// third-party package.
	// +optional
	filepath string
	// The list of dependencies of the module, i.e. the imports of the module.
	// +optional
	dependencies []pyModule
}

// requirement represents a requirement from a requirements.txt file.
type requirement struct {
	// The raw line for the requirement extracted from the requirements.txt
	// file.
	raw string
	// The name parsed from the raw line.
	name string
	// The version parsed from the raw line.
	version string
	// The hashes parsed from the raw line.
	hashes []string
}

// parse parses the name, version and hashes from the raw line of the
// requirement.
func (r *requirement) parse() {
	split := strings.Split(r.raw, "==")
	r.raw = ""
	r.name = split[0]
	if len(split) == 1 {
		return
	}
	split = strings.Split(split[1], " ")
	r.version = split[0]
	if len(split) == 1 {
		return
	}
	args := split[1:]
	hashes := make([]string, 0, len(args))
	for _, arg := range args {
		// Halt processing if it hits a comment.
		if strings.HasPrefix(arg, "#") {
			break
		}
		if strings.HasPrefix(arg, "--hash") {
			hash := strings.TrimPrefix(arg, "--hash")
			hash = strings.TrimPrefix(hash, "=")
			hashes = append(hashes, hash)
		}
	}
	r.hashes = hashes
}
