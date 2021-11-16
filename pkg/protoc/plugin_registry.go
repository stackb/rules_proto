package protoc

import "errors"

// ErrUnknownPlugin is the error returned when a plugin is not known.
var ErrUnknownPlugin = errors.New("unknown plugin")

// PluginRegistry represents a library of plugin implementations.
type PluginRegistry interface {
	// PluginNames returns a sorted list of plugin names.
	PluginNames() []string
	// LookupPlugin returns the implementation under the given name.  If the
	// plugin is not found, ErrUnknownPlugin is returned.
	LookupPlugin(name string) (Plugin, error)
	// MustRegisterPlugin installs a Plugin implementation under the given name
	// in the global plugin registry.  Panic will occur if the same plugin is
	// registered multiple times.
	MustRegisterPlugin(plugin Plugin) PluginRegistry
}
