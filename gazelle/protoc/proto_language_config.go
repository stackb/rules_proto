package protoc

import (
	"log"
	"path"
	"strconv"
)

type ProtoLanguageConfig struct {
	Name           string
	Load           string
	Kind           string
	Enabled        bool
	Pattern        string
	Plugins        []*ProtoPluginConfig
	Implementation ProtoLanguage
}

// MustParseDirective parses the directive string or panics.
func (s *ProtoLanguageConfig) MustParseDirective(cfg *ProtoPackageConfig, d, param, value string) {
	switch param {
	case "enable":
		enabled, err := strconv.ParseBool(value)
		if err != nil {
			log.Fatalf("bad directive %q: value %q is not a boolean", d, value)
		}
		s.Enabled = enabled
	case "match":
		if _, err := path.Match(value, ""); err == path.ErrBadPattern {
			log.Fatalf("bad directive %q: value %q is not a glob: %v", d, value, err)
		}
		s.Pattern = value
	case "plugin":
		if plugin, ok := cfg.plugins[value]; !ok {
			log.Fatalf("bad directive %q: unknown plugin %q", d, value)
		} else {
			s.Plugins = append(s.Plugins, plugin)
		}
	case "kind":
		s.Kind = value
	case "load":
		s.Load = value
	default:
		log.Fatalf("bad directive %q: unknown parameter %q", d, value)
	}
}

// Clone, well... it clones the config.
func (s *ProtoLanguageConfig) Clone() *ProtoLanguageConfig {
	clone := &ProtoLanguageConfig{
		Name:           s.Name,
		Load:           s.Load,
		Kind:           s.Kind,
		Enabled:        s.Enabled,
		Pattern:        s.Pattern,
		Plugins:        make([]*ProtoPluginConfig, len(s.Plugins)),
		Implementation: s.Implementation,
	}
	for i, p := range s.Plugins {
		clone.Plugins[i] = p.Clone()
	}
	return clone
}
