package protoc

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/emicklei/proto"
)

// NewFile takes the package directory and base name of the file (e.g.
// 'foo.proto') and constructs File
func NewFile(dir, basename string) *File {
	return &File{
		Dir:      dir,
		Basename: basename,
		Name:     strings.TrimSuffix(basename, filepath.Ext(basename)),
	}
}

// File represents a proto file that is discovered in a package.
type File struct {
	Dir      string // e.g. "path/to/package/"
	Basename string // e.g. "foo.proto"
	Name     string // e.g. "foo"

	pkg         proto.Package
	imports     []proto.Import
	options     []proto.Option
	services    []proto.Service
	messages    []proto.Message
	enums       []proto.Enum
	enumOptions []proto.Option
}

// Relname returns the relative path of the proto file.
func (f *File) Relname() string {
	if f.Dir == "" {
		return f.Basename
	}
	return filepath.Join(f.Dir, f.Basename)
}

// Package returns the defined package or the empty value.
func (f *File) Package() proto.Package {
	return f.pkg
}

// Imports returns the list of Imports defined in the proto file.
func (f *File) Imports() []proto.Import {
	return f.imports
}

// Options returns the list of top-level options defined in the proto file.
func (f *File) Options() []proto.Option {
	return f.options
}

// Services returns the list of Services defined in the proto file.
func (f *File) Services() []proto.Service {
	return f.services
}

// Messages returns the list of Messages defined in the proto file.
func (f *File) Messages() []proto.Message {
	return f.messages
}

// Enums returns the list of Enums defined in the proto file.
func (f *File) Enums() []proto.Enum {
	return f.enums
}

// EnumOptions returns the list of EnumOptions defined in the proto file.
func (f *File) EnumOptions() []proto.Option {
	return f.enumOptions
}

// HasEnums returns true if the proto file has at least one enum.
func (f *File) HasEnums() bool {
	return len(f.enums) > 0
}

// HasMessages returns true if the proto file has at least one message.
func (f *File) HasMessages() bool {
	return len(f.messages) > 0
}

// HasServices returns true if the proto file has at least one service.
func (f *File) HasServices() bool {
	return len(f.services) > 0
}

// HasEnumOption returns true if the proto file has at least one enum or enum
// field annotated with the given named field extension.
func (f *File) HasEnumOption(name string) bool {
	for _, option := range f.enumOptions {
		if option.Name == name {
			return true
		}
	}
	return false
}

// Parse reads the proto file and parses the source.
func (f *File) Parse() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not parse: %v", err)
	}

	if bwd, ok := os.LookupEnv("BUILD_WORKSPACE_DIRECTORY"); ok {
		wd = bwd
	}

	filename := filepath.Join(wd, f.Dir, f.Basename)
	reader, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open %s: %w (cwd=%s)", filename, err, wd)
	}
	defer reader.Close()

	return f.ParseReader(reader)
}

// ParseReader parses the reader and walks statements in the file.
func (f *File) ParseReader(in io.Reader) error {
	parser := proto.NewParser(in)
	definition, err := parser.Parse()
	if err != nil {
		return fmt.Errorf("could not parse %s/%s: %w", f.Dir, f.Basename, err)
	}

	proto.Walk(definition,
		proto.WithPackage(f.handlePackage),
		proto.WithOption(f.handleOption),
		proto.WithImport(f.handleImport),
		proto.WithService(f.handleService),
		proto.WithMessage(f.handleMessage),
		proto.WithEnum(f.handleEnum))

	// NOTE: f.options only holds top-level options.  To introspect the enum and
	// enum field options we need to do extra work.
	collector := &protoEnumOptionCollector{}
	for _, enum := range f.enums {
		collector.VisitEnum(&enum)
	}
	f.enumOptions = collector.options

	return nil
}

func (f *File) handlePackage(p *proto.Package) {
	f.pkg = *p
}

func (f *File) handleOption(o *proto.Option) {
	f.options = append(f.options, *o)
}

func (f *File) handleImport(i *proto.Import) {
	f.imports = append(f.imports, *i)
}

func (f *File) handleEnum(i *proto.Enum) {
	f.enums = append(f.enums, *i)
}

func (f *File) handleService(s *proto.Service) {
	f.services = append(f.services, *s)
}

func (f *File) handleMessage(m *proto.Message) {
	f.messages = append(f.messages, *m)
}

// PackageFileNameWithExtensions returns a function that computes the name of a
// predicted generated file having the given extension(s).  If the proto package
// is defined, the output file will be in the corresponding directory.
func PackageFileNameWithExtensions(exts ...string) func(f *File) []string {
	return func(f *File) []string {
		outs := make([]string, len(exts))
		name := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			name = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), name)
		}
		for i, ext := range exts {
			outs[i] = name + ext
		}
		return outs
	}
}

// RelativeFileNameWithExtensions returns a function that computes the name of a
// predicted generated file having the given extension(s) relative to the given
// dir.
func RelativeFileNameWithExtensions(reldir string, exts ...string) func(f *File) []string {
	return func(f *File) []string {
		outs := make([]string, len(exts))
		name := f.Name
		if reldir != "" {
			name = path.Join(reldir, name)
		}
		for i, ext := range exts {
			outs[i] = name + ext
		}
		return outs
	}
}

// ImportPrefixRelativeFileNameWithExtensions returns a function that computes
// the name of a predicted generated file. In this case, first
// RelativeFileNameWithExtensions is applied, then stripImportPrefix is removed
// from the predicted filename.
func ImportPrefixRelativeFileNameWithExtensions(stripImportPrefix, reldir string, exts ...string) func(f *File) []string {
	// if the stripImportPrefix is defined and "absolute" (starting with a
	// slash), this means it is relative to the repository root.
	// https://github.com/bazelbuild/bazel/issues/3867#issuecomment-441971525
	prefix := stripImportPrefix
	if strings.HasPrefix(prefix, "/") {
		prefix = prefix[1:]
	}
	relfunc := RelativeFileNameWithExtensions(reldir, exts...)
	return func(f *File) []string {
		outs := relfunc(f)
		for i, out := range outs {
			if strings.HasPrefix(out, prefix) {
				outs[i] = strings.TrimPrefix(out[len(prefix):], "/")
			}
		}
		return outs
	}
}

// HasMessagesOrEnums checks if any of the given files has a message or an enum.
func HasMessagesOrEnums(files ...*File) bool {
	for _, f := range files {
		if HasMessageOrEnum(f) {
			return true
		}
	}
	return false
}

// HasServices checks if any of the given files has a service.
func HasServices(files ...*File) bool {
	for _, f := range files {
		if HasService(f) {
			return true
		}
	}
	return false
}

// HasMessageOrEnum is a file predicate function checks if any of the given file
// has a message or an enum.
func HasMessageOrEnum(file *File) bool {
	return file.HasMessages() || file.HasEnums()
}

// Always is a file predicate function that always returns true.
func Always(file *File) bool {
	return true
}

// HasService is a file predicate function that tests if any of the given file
// has a message or an enum.
func HasService(file *File) bool {
	return file.HasServices()
}

// FlatMapFiles is a utility function intended for use in computing a list of
// output files for a given proto_library. The given apply function is executed
// foreach file that passes the filter function, and flattens the strings into a
// single list.
func FlatMapFiles(apply func(file *File) []string, filter func(file *File) bool, files ...*File) []string {
	values := make([]string, 0)
	for _, f := range files {
		if !filter(f) {
			continue
		}
		values = append(values, apply(f)...)
	}
	return values
}

// GoPackagePath replaces dots with forward slashes.
func GoPackagePath(pkg string) string {
	return strings.ReplaceAll(pkg, ".", "/")
}

// IsProtoFile returns true if the file extension looks like it should contain
// protobuf definitions.
func IsProtoFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".proto" || ext == ".protodevel"
}

// GoPackageOption is a utility function to seek for the go_package option and
// split it.  If present the return values will be populated with the importpath
// and alias (e.g. github.com/foo/bar/v1;bar -> "github.com/foo/bar/v1", "bar").
// If the option was not found the bool return argument is false.
func GoPackageOption(options []proto.Option) (string, string, bool) {
	for _, opt := range options {
		if opt.Name != "go_package" {
			continue
		}
		parts := strings.SplitN(opt.Constant.Source, ";", 2)
		switch len(parts) {
		case 0:
			return "", "", true
		case 1:
			return parts[0], "", true
		case 2:
			return parts[0], parts[1], true
		default:
			return parts[0], strings.Join(parts[1:], ";"), true
		}
	}

	return "", "", false
}

// GetNamedOption returns the value of an option.  If the option is not found,
// the bool return value is false.
func GetNamedOption(options []proto.Option, name string) (string, bool) {
	for _, opt := range options {
		if opt.Name != name {
			continue
		}
		return opt.Constant.Source, true
	}
	return "", false
}

// ToPascalCase converts a string to PascalCase.
//
// Splits on '-', '_', ' ', '\t', '\n', '\r'.
// Uppercase letters will stay uppercase,
func ToPascalCase(s string) string {
	output := ""
	var previous rune
	for i, c := range strings.TrimSpace(s) {
		if !isDelimiter(c) {
			if i == 0 || isDelimiter(previous) || unicode.IsUpper(c) {
				output += string(unicode.ToUpper(c))
			} else {
				output += string(unicode.ToLower(c))
			}
		}
		previous = c
	}
	return output
}

func isDelimiter(r rune) bool {
	return r == '.' || r == '-' || r == '_' || r == ' ' || r == '\t' || r == '\n' || r == '\r'
}
