package protoc

// PyLanguageName is the name of the proto_descriptor_set implementation.
const PyLanguageName = "proto_descriptor_set"

func init() {
	MustRegisterProtoLanguage(PyLanguageName, &PyLanguage{})
}

// PyLanguage implements a ProtoLanguage for plugins that generate *.py files.
type PyLanguage struct{}

// GenerateRules implements the ProtoLanguage interface.
func (s *PyLanguage) GenerateRules(
	rel string,
	cfg *ProtoPackageConfig,
	libs []ProtoLibrary,
) []RuleProvider {
	rules := make([]RuleProvider, 0)

	for _, lib := range libs {
		rules = append(rules, NewProtoDescriptorSet(lib))
	}

	return rules
}
