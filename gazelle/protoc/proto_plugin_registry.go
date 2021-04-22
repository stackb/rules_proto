package protoc

import "fmt"

var (
	// the plugin registry
	plugins = make(map[string]ProtoPlugin)
	// ErrUnknownPlugin is the error returned when a plugin is not known.
	ErrUnknownPlugin = fmt.Errorf("Unknown ProtoPlugin")
)

// LookupProtoPlugin returns the implementation under the given name.
func LookupProtoPlugin(name string) (ProtoPlugin, error) {
	lang, ok := plugins[name]
	if !ok {
		return nil, ErrUnknownPlugin
	}
	return lang, nil
}

// MustRegisterProtoPlugin installs a ProtoPlugin implementation under the
// given name in the global plugin registry.  Panic will occur of the same
// plugin is registered multiple times.
func MustRegisterProtoPlugin(name string, lang ProtoPlugin) {
	_, ok := plugins[name]
	if ok {
		panic("duplicate proto_plugin registration: " + name)
	}
	plugins[name] = lang
}
