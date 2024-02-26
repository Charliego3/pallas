package main

import (
	"fmt"
	"net/http"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	contextPackage = protogen.GoImportPath("context")
	httpxPackage   = protogen.GoImportPath("github.com/charliego3/pallas/httpx")
	typesPackage   = protogen.GoImportPath("github.com/charliego3/pallas/types")
)

type method struct {
	name    string
	method  string
	path    string
	handler string
	in, out string
}

// generate generates a _http.http.go file containing HTTP service definitions.
func generate(gen *protogen.Plugin, f *protogen.File) {
	if len(f.Services) == 0 || !hasHTTPMethod(f) {
		return
	}

	httpname := f.GeneratedFilenamePrefix + "_http.pb.go"
	descname := f.GeneratedFilenamePrefix + "_desc.pb.go"
	hg := addHeader(gen, f, httpname)
	dg := addHeader(gen, f, descname)
	hg.QualifiedGoIdent(contextPackage.Ident(""))
	hg.P("// This is a compile-time assertion to ensure that this generated file")
	hg.P("// is compatible with the pallas package it is being compiled against.")
	hg.P("var _ = new(", httpxPackage.Ident("CallOption"), ")")
	hg.P("var _ = new(", typesPackage.Ident("Service"), ")")
	hg.P()

	dg.P("// This is a compile-time assertion to ensure that this generated file")
	dg.P("// is compatible with the pallas package it is being compiled against.")
	dg.P("var _ = new(", typesPackage.Ident("Service"), ")")
	dg.P()
	for _, s := range f.Services {
		methods := getMethods(s)
		generateService(gen, hg, s, methods)
		generateDesc(gen, dg, s, methods)
	}
}

func checkDeprecate(s *protogen.Service, g *protogen.GeneratedFile) {
	if s.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		g.P("//")
		g.P("// Deprecated: Do not use.")
	}
}

func generateService(gen *protogen.Plugin, g *protogen.GeneratedFile, s *protogen.Service, methods []method) {
	checkDeprecate(s, g)
	g.P("type ", s.GoName, "HTTPServer interface {")
	for _, m := range methods {
		g.P("\t", m.name, "(ctx context.Context, in *", m.in, ") (*", m.out, ", error)")
	}
	g.P("}")
	g.P()

	g.P("func Register", s.GoName, "HTTPServer(s *httpx.Server, srv ", s.GoName, "HTTPServer) {")
	for _, m := range methods {
		g.P("\ts.HandleMethod(\"", m.method, "\", \"", m.path, "\", ", m.handler, "(srv.(types.Service)).(httpx.Handler))")
	}
	g.P("}")
	g.P()

	for _, m := range methods {
		g.P("func ", m.handler, "(srv types.Service) any {")
		g.P("\treturn httpx.HandlerFunc(func(c *httpx.Context) error {")
		g.P("\t\treq := new(", m.in, ")")
		g.P("\t\tif err := c.Bind(req); err != nil {")
		g.P("\t\t\treturn err")
		g.P("\t\t}")
		g.P("\t\treturn srv.(", s.GoName, "Server).", m.name, "(c.Context, req)")
		g.P("\t})")
		g.P("}")
		g.P()
	}
}

func generateDesc(gen *protogen.Plugin, g *protogen.GeneratedFile, s *protogen.Service, methods []method) {
	checkDeprecate(s, g)
	name := "Unimplemented" + s.GoName + "DescServer"
	g.P("type ", name, " struct {")
	g.P("\tUnimplemented", s.GoName, "Server")
	g.P("}")
	g.P()
	g.P("func (", name, ") Desc() types.ServiceDesc {")
	g.P("\treturn types.ServiceDesc{")
	g.P("\t\tGrpc: ", s.GoName, "_ServiceDesc,")
	g.P("\t\tHttp: types.HttpServiceDesc{")
	g.P("\t\t\tServiceName: \"\",")
	g.P("\t\t\tHandlerType: nil,")
	g.P("\t\t\tMethods: []types.HttpMethodDesc{")
	for _, m := range methods {
		g.P("\t\t\t\t{")
		g.P("\t\t\t\t\tMethod: \"", m.method, "\",")
		g.P("\t\t\t\t\tTemplate: \"", m.path, "\",")
		g.P("\t\t\t\t\tHandler: ", m.handler, ",")
		g.P("\t\t\t\t},")
	}
	g.P("\t\t\t},\n\t\t},\n\t}\n}\n")
}

func getMethods(s *protogen.Service) (requests []method) {
	for _, m := range s.Methods {
		desc := m.Desc
		if desc.IsStreamingClient() || desc.IsStreamingServer() {
			continue
		}

		if rule, ok := proto.GetExtension(
			m.Desc.Options(),
			annotations.E_Http,
		).(*annotations.HttpRule); ok {
			for _, binding := range append(rule.AdditionalBindings, rule) {
				var method method
				switch pattern := binding.Pattern.(type) {
				case *annotations.HttpRule_Get:
					method.path = pattern.Get
					method.method = http.MethodGet
				case *annotations.HttpRule_Post:
					method.path = pattern.Post
					method.method = http.MethodPost
				case *annotations.HttpRule_Put:
					method.path = pattern.Put
					method.method = http.MethodPut
				case *annotations.HttpRule_Delete:
					method.path = pattern.Delete
					method.method = http.MethodDelete
				case *annotations.HttpRule_Patch:
					method.path = pattern.Patch
					method.method = http.MethodPatch
				case *annotations.HttpRule_Custom:
					method.path = pattern.Custom.Path
					method.method = pattern.Custom.Kind
				}
				method.name = m.GoName
				method.handler = fmt.Sprintf("_%s_%s_%s_HTTP_Handler", s.GoName, m.GoName, method.method)
				method.in = string(m.Desc.Input().Name())
				method.out = string(m.Desc.Output().Name())
				requests = append(requests, method)
			}
		}
	}
	return
}

func hasHTTPMethod(f *protogen.File) bool {
	for _, serv := range f.Services {
		for _, method := range serv.Methods {
			desc := method.Desc
			if desc.IsStreamingServer() || desc.IsStreamingClient() {
				continue
			}

			rule, ok := proto.GetExtension(
				desc.Options(),
				annotations.E_Http,
			).(*annotations.HttpRule)
			if ok && rule != nil {
				return true
			}
		}
	}
	return false
}

func getProtocVersion(gen *protogen.Plugin) string {
	v := gen.Request.GetCompilerVersion()
	if v == nil {
		return "unknow"
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), v.GetSuffix())
}
