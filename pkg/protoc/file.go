package protoc

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

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

// File represents a proto file we discover in a package.
type File struct {
	Dir      string // e.g. "rosetta/rosetta/common/"
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

// GetPackage returns the defined package or the empty value.
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

// Parse the source and walk the statements in the file.
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

	return f.parseReader(reader)
}

func (f *File) parseReader(in io.Reader) error {
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

// GoPackagePath replaces dots with forward slashes.
func GoPackagePath(pkg string) string {
	return strings.ReplaceAll(pkg, ".", "/")
}

// IsProtoFile returns true if the file extension looks like it should contain protobuf definitions.
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
