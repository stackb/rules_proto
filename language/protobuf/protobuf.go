package protobuf

import (
	"github.com/bazelbuild/bazel-gazelle/language"

	"github.com/stackb/rules_proto/v4/pkg/language/protobuf"

	_ "github.com/stackb/rules_proto/v4/pkg/plugin/bufbuild"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/builtin"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/gogo/protobuf"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/golang/protobuf"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/grpc/grpc"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/grpc/grpcgo"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/grpc/grpcjava"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/grpc/grpcnode"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/grpc/grpcweb"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/grpcecosystem/grpcgateway"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/prost"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/prost_serde"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/tonic"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/scalapb/scalapb"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/scalapb/zio_grpc"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/stackb/grpc_js"
	_ "github.com/stackb/rules_proto/v4/pkg/plugin/stephenh/ts-proto"
	_ "github.com/stackb/rules_proto/v4/pkg/rule/rules_cc"
	_ "github.com/stackb/rules_proto/v4/pkg/rule/rules_closure"
	_ "github.com/stackb/rules_proto/v4/pkg/rule/rules_go"
	_ "github.com/stackb/rules_proto/v4/pkg/rule/rules_java"
	_ "github.com/stackb/rules_proto/v4/pkg/rule/rules_nodejs"
	_ "github.com/stackb/rules_proto/v4/pkg/rule/rules_python"
	_ "github.com/stackb/rules_proto/v4/pkg/rule/rules_rust"
	_ "github.com/stackb/rules_proto/v4/pkg/rule/rules_scala"
)

// NewLanguage is called by Gazelle to install this language extension in a
// binary.  This package "language/protobuf" is separate from
// "pkg/language/protobuf" because this one bundles all the plugin and rule
// implementations in the repo whereas the other is the "pure" language
// implementation, with no pre-population of the registries.
func NewLanguage() language.Language {
	return protobuf.NewProtobufLang("protobuf")
}
