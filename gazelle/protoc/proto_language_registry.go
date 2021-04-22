package protoc

import "fmt"

var (
	// the language registry
	languages = make(map[string]ProtoLanguage)
	// ErrUnknownLanguage is the error returned when a language is not known.
	ErrUnknownLanguage = fmt.Errorf("Unknown ProtoLanguage")
)

// LookupProtoLanguage returns the implementation under the given name.
func LookupProtoLanguage(name string) (ProtoLanguage, error) {
	lang, ok := languages[name]
	if !ok {
		return nil, ErrUnknownLanguage
	}
	return lang, nil
}

// MustLookupProtoLanguage returns the implementation under the given name or
// panics.
func MustLookupProtoLanguage(name string) ProtoLanguage {
	lang, err := LookupProtoLanguage(name)
	if err != nil {
		panic(fmt.Sprintf("invalid or unknown proto_language %q: %v", name, err))
	}
	return lang
}

// MustRegisterProtoLanguage installs a ProtoLanguage implementation under the
// given name in the global language registry.  Panic will occur of the same
// language is registered multiple times.
func MustRegisterProtoLanguage(name string, lang ProtoLanguage) {
	_, ok := languages[name]
	if ok {
		panic("duplicate proto_language registration: " + name)
	}
	languages[name] = lang
}
