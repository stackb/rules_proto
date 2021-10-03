package protoc

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v3"
)

// YConfig is used to configure a combined set of plugins, rules, and languages
// in a single YAML file.  This is the format of the -proto_language_config_file
// flag.
type YConfig struct {
	Plugin   []*YPlugin   `yaml:"plugins"`
	Rule     []*YRule     `yaml:"rules"`
	Language []*YLanguage `yaml:"languages"`
}

// YPlugin represents a LanguagePluginConfig in YAML.
type YPlugin struct {
	Name           string   `yaml:"name"`
	Implementation string   `yaml:"implementation"`
	Enabled        bool     `yaml:"enabled"`
	Option         []string `yaml:"options"`
	Label          string   `yaml:"label"`
}

// YRule represents a LanguageRuleConfig in YAML.
type YRule struct {
	Name           string   `yaml:"name"`
	Implementation string   `yaml:"implementation"`
	Enabled        bool     `yaml:"enabled"`
	Deps           []string `yaml:"deps"`
	Visibility     []string `yaml:"visibility"`
}

// YLanguage represents a LanguageConfig in YAML.
type YLanguage struct {
	Name           string   `yaml:"name"`
	Implementation string   `yaml:"implementation"`
	Enabled        bool     `yaml:"enabled"`
	Plugin         []string `yaml:"plugins"`
	Rule           []string `yaml:"rules"`
}

// ParseYConfigFile parses the given filename and returns a YConfig pointer or
// error if a file read or parse error occurs.
func ParseYConfigFile(filename string) (*YConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("yaml read error %s: %w", filename, err)
	}
	var config YConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("yaml parse error %s: %w", filename, err)
	}
	return &config, nil
}
