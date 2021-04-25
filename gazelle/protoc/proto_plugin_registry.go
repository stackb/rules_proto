package protoc

import "errors"

// ErrUnknownPlugin is the error returned when a plugin is not known.
var ErrUnknownPlugin = errors.New("unknown ProtoPlugin")

// ProtoPluginRegistry represents a library of plugin implementations.
type ProtoPluginRegistry interface {
	// PluginNames returns a sorted list of plugin names.
	PluginNames() []string
	// LookupProtoPlugin returns the implementation under the given name.
	LookupProtoPlugin(name string) (ProtoPlugin, error)
	// MustRegisterProtoPlugin installs a ProtoPlugin implementation under the given
	// name in the global plugin registry.  Panic will occur of the same plugin is
	// registered multiple times.
	MustRegisterProtoPlugin(name string, plugin ProtoPlugin)
}

// PluginNames for the .
func PluginNames() []string {
	return protoc.PluginNames()
}

// LookupProtoPlugin returns the implementation under the given name.
func LookupProtoPlugin(name string) (ProtoPlugin, error) {
	return protoc.LookupProtoPlugin(name)
}

// MustRegisterProtoPlugin installs a ProtoPlugin implementation under the given
// name in the global plugin registry.  Panic will occur of the same plugin is
// registered multiple times.
func MustRegisterProtoPlugin(name string, plugin ProtoPlugin) {
	protoc.MustRegisterProtoPlugin(name, plugin)
}
