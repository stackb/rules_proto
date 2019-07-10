package main

func makePhp() *Language {
	return &Language{
		Dir:   "php",
		Name:  "php",
		DisplayName: "PHP",
		Notes: mustTemplate("Rules for generating PHP protobuf and gRPC `.php` files and libraries using standard Protocol Buffers and gRPC"),
		Flags: commonLangFlags,
		Rules: []*Rule{
			&Rule{
				Name:             "php_proto_compile",
				Kind:             "proto",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//php:php"},
				WorkspaceExample: protoWorkspaceTemplate,
				BuildExample:     protoCompileExampleTemplate,
				Doc:              "Generates PHP protobuf `.php` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
			&Rule{
				Name:             "php_grpc_compile",
				Kind:             "grpc",
				Implementation:   aspectRuleTemplate,
				Plugins:          []string{"//php:php", "//php:grpc_php"},
				WorkspaceExample: grpcWorkspaceTemplate,
				BuildExample:     grpcCompileExampleTemplate,
				Doc:              "Generates PHP protobuf+gRPC `.php` artifacts",
				Attrs:            aspectProtoCompileAttrs,
			},
		},
	}
}
