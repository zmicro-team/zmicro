package main

import (
	"strconv"

	"google.golang.org/protobuf/compiler/protogen"
)

type serviceDesc struct {
	Deprecated  bool   // 是否弃用
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/v1/helloworld.proto
	Comment     string // 注释
	Methods     []*methodDesc
}

type methodDesc struct {
	Deprecated bool // 是否弃用
	// method
	Name    string // 方法名
	Num     int    // 方法号
	Request string // 请求结构
	Reply   string // 回复结构
	Comment string // 方法注释
	// http_rule
	Path         string // 路径
	Method       string // 方法
	HasVars      bool   // 是否有url参数
	HasBody      bool   // 是否有消息体
	Body         string // 请求消息体
	ResponseBody string // 回复消息体
}

func executeServiceDesc(g *protogen.GeneratedFile, s *serviceDesc) error {
	// http interface defined
	if s.Deprecated {
		g.P(deprecationComment)
	}
	g.P("// ", clientInterfaceName(s.ServiceType), " ", s.Comment)
	g.P("type ", clientInterfaceName(s.ServiceType), " interface {")
	for _, m := range s.Methods {
		if m.Deprecated {
			g.P(deprecationComment)
		}
		g.P(m.Comment)
		g.P(clientMethodName(g, m, true))
	}
	g.P("}")
	g.P()

	// http client implement.
	g.P("type ", clientImplStructName(s.ServiceType), " struct {")
	g.P("cc *", g.QualifiedGoIdent(transportHttpPackage.Ident("Client")))
	g.P("}")
	g.P()
	// http client factory method.
	if s.Deprecated {
		g.P(deprecationComment)
	}
	g.P("// ", clientNewFunctionName(s.ServiceType), " ", s.Comment)
	g.P("func ", clientNewFunctionName(s.ServiceType), "(c *", g.QualifiedGoIdent(transportHttpPackage.Ident("Client")), ") ", clientInterfaceName(s.ServiceType), " {")
	g.P("return &", clientImplStructName(s.ServiceType), " {")
	g.P("cc: c,")
	g.P("}")
	g.P("}")
	g.P()

	// http client implement methods.
	for _, m := range s.Methods {
		if m.Deprecated {
			g.P(deprecationComment)
		}
		g.P(m.Comment)
		g.P("func (c *", clientImplStructName(s.ServiceType), ")", clientMethodName(g, m, false), " {")
		g.P("var err error")
		g.P("var resp ", m.Reply)
		g.P()
		g.P(`settings := c.cc.CallSetting("`, m.Path, `", opts...)`)

		if m.HasVars {
			g.P("path := c.cc.EncodeURL(settings.Path, req, ", strconv.FormatBool(!m.HasBody), ")")
		} else {
			if m.HasBody {
				g.P("path := settings.Path")
			} else {
				g.P("var query string")
				g.P()
				g.P("query, err = c.cc.EncodeQuery(req)")
				g.P("if err != nil {")
				g.P("return nil, err")
				g.P("}")
				g.P("path := settings.Path")
				g.P(`if query != "" {`)
				g.P(`path += "?" + query`)
				g.P("}")
			}
		}

		reqValue := "nil"
		if m.HasBody {
			reqValue = "req" + m.Body
		}
		if args.UseInvoke2 {
			g.P(`err = c.cc.Invoke2(ctx, "`, m.Method, `", path, `, reqValue, ", &resp", m.ResponseBody, ", settings)")
		} else {
			g.P("ctx = ", g.QualifiedGoIdent(transportHttpPackage.Ident("WithValueCallOption")), "(ctx, settings)")
			g.P(`err = c.cc.Invoke(ctx, "`, m.Method, `", path, `, reqValue, ", &resp", m.ResponseBody, ")")
		}
		g.P("if err != nil {")
		g.P("return nil, err")
		g.P("}")
		g.P("return &resp, nil")
		g.P("}")
		g.P()
	}

	return nil
}

func clientInterfaceName(serviceType string) string {
	return serviceType + "HTTPClient"
}

func clientImplStructName(serviceType string) string {
	return serviceType + "HTTPClientImpl"
}

func clientNewFunctionName(serviceType string) string {
	return "New" + serviceType + "HTTPClient"
}

func clientMethodName(g *protogen.GeneratedFile, m *methodDesc, isDeclaration bool) string {
	ctxParam := ""
	reqParam := ""
	optsParam := ""
	if !isDeclaration {
		ctxParam = "ctx"
		reqParam = "req"
		optsParam = "opts"
	}

	num := ""
	if m.Num != 0 { // unique because additional_bindings
		num = "_" + strconv.Itoa(m.Num)
	}

	return m.Name + num + "(" + ctxParam + " " + g.QualifiedGoIdent(contextPackage.Ident("Context")) +
		", " + reqParam + " *" + m.Request + ", " +
		optsParam + " ..." + g.QualifiedGoIdent(transportHttpPackage.Ident("CallOption")) +
		") (*" + m.Reply + ", error)"
}
