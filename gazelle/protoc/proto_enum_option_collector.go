package protoc

import (
	"github.com/emicklei/proto"
)

// protoEnumOptionCollector implements the proto.Visitor interface.  It is
// applied to top-level Enums gathered during the walk phase to search deeper
// within the proto ast.
type protoEnumOptionCollector struct {
	// the collected options after visiting the nodes.  Includes both options on
	// enums as well as enum field options.
	options []proto.Option
}

func (c *protoEnumOptionCollector) VisitMessage(m *proto.Message) {}
func (c *protoEnumOptionCollector) VisitService(v *proto.Service) {}
func (c *protoEnumOptionCollector) VisitSyntax(s *proto.Syntax)   {}
func (c *protoEnumOptionCollector) VisitPackage(p *proto.Package) {}
func (c *protoEnumOptionCollector) VisitOption(o *proto.Option) {
	c.options = append(c.options, *o)
}
func (c *protoEnumOptionCollector) VisitImport(i *proto.Import)           {}
func (c *protoEnumOptionCollector) VisitNormalField(i *proto.NormalField) {}
func (c *protoEnumOptionCollector) VisitEnumField(i *proto.EnumField) {
	for _, v := range i.Elements {
		v.Accept(c)
	}
}
func (c *protoEnumOptionCollector) VisitEnum(e *proto.Enum) {
	for _, v := range e.Elements {
		v.Accept(c)
	}
}
func (c *protoEnumOptionCollector) VisitComment(e *proto.Comment)       {}
func (c *protoEnumOptionCollector) VisitOneof(o *proto.Oneof)           {}
func (c *protoEnumOptionCollector) VisitOneofField(o *proto.OneOfField) {}
func (c *protoEnumOptionCollector) VisitReserved(r *proto.Reserved)     {}
func (c *protoEnumOptionCollector) VisitRPC(r *proto.RPC)               {}
func (c *protoEnumOptionCollector) VisitMapField(f *proto.MapField)     {}
func (c *protoEnumOptionCollector) VisitGroup(g *proto.Group)           {}
func (c *protoEnumOptionCollector) VisitExtensions(e *proto.Extensions) {}
