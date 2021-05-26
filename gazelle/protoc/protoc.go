package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/language"

	extension "github.com/stackb/rules_proto/pkg/gazelle/protoc"

	_ "github.com/stackb/rules_proto/pkg/builtin"
	_ "github.com/stackb/rules_proto/pkg/gogo/protobuf"
	_ "github.com/stackb/rules_proto/pkg/grpc/grpcgo"
	_ "github.com/stackb/rules_proto/pkg/grpc/grpcjava"
	_ "github.com/stackb/rules_proto/pkg/grpcecosystem/grpcgateway"
)

// NewLanguage is called by Gazelle to install this language extension in a
// binary.
func NewLanguage() language.Language {
	return extension.NewProtoc("protoc")
}
