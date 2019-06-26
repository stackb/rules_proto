package main

func makePhp() *Language {
	return &Language{
		Dir:   "php",
		Name:  "php",
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "php_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//php:php"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates php protobuf artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "php_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//php:php", "//php:grpc_php"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates php protobuf+gRPC artifacts",
				Attrs:            aspectProtoCompileAttrs,
				BazelCIExclusionReason: "experimental",
				Experimental:     true,
			},
		},
	}
}
