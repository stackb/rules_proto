package protoc

// ProtoLanguage implementations are capable of providing a set of rules based
// on a proto package config and a set of proto_library instances.
type ProtoLanguage interface {
	// GenerateRules takes the config and a list of *ProtoLibrary and returns a
	// preliminary list of RuleProviders. Implementation should note that the
	// preliminary list may be filtered via other gazelle directives.
	// rel is the relative directory path to the build file (e.g. GenerateArgs.Rel)
	GenerateRules(rel string, c *ProtoPackageConfig, p *ProtoLanguageConfig, libs []ProtoLibrary) []RuleProvider
}
