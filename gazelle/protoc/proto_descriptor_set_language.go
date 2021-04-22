package protoc

// ProtoDescriptorSetLanguageName is the name of the proto_descriptor_set implementation.
const ProtoDescriptorSetLanguageName = "proto_descriptor_set"

func init() {
	MustRegisterProtoLanguage(ProtoDescriptorSetLanguageName, &ProtoDescriptorSetLanguage{})
}

// ProtoDescriptorSetLanguage implements a ProtoLanguage for
// proto_descriptor_set targets.
type ProtoDescriptorSetLanguage struct{}

// GenerateRules implements the ProtoLanguage interface.
func (s *ProtoDescriptorSetLanguage) GenerateRules(
	rel string,
	c *ProtoPackageConfig,
	p *ProtoLanguageConfig,
	libs []ProtoLibrary,
) []RuleProvider {
	rules := make([]RuleProvider, 0)

	for _, lib := range libs {
		rules = append(rules, NewProtoDescriptorSet(lib))
	}

	return rules
}
