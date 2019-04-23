package main

func makePhp() *Language {
	return &Language{
		Dir:  "php",
		Name: "php",
		Rules: []*Rule{
			&Rule{
				Name:           "php_proto_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//php:php"},
				Usage:          usageTemplate,
				Example:        protoCompileExampleTemplate,
				Doc:            "Generates php protobuf artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
			&Rule{
				Name:           "php_grpc_compile",
				Implementation: compileRuleTemplate,
				Plugins:        []string{"//php:php", "//php:grpc_php"},
				Usage:          grpcUsageTemplate,
				Example:        grpcCompileExampleTemplate,
				Doc:            "Generates php protobuf+gRPC artifacts",
				Attrs:          append(protoCompileAttrs, []*Attr{}...),
			},
		},
	}
}
