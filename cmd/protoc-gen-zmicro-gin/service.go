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

	RpcMode     string
	UseEncoding bool

	// not used when disable template
	MethodSets map[string]*methodDesc // unique because additional_bindings
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
	methodSets := make(map[string]struct{})
	// http interface defined
	if s.Deprecated {
		g.P(deprecationComment)
	}
	g.P("// ", serverInterfaceName(s.ServiceType), " ", s.Comment)
	g.P("type ", serverInterfaceName(s.ServiceType), " interface {")
	for _, m := range s.Methods {
		_, ok := methodSets[m.Name]
		if ok { // unique because additional_bindings
			continue
		}
		methodSets[m.Name] = struct{}{}
		if m.Deprecated {
			g.P(deprecationComment)
		}
		g.P(m.Comment)
		if s.RpcMode == "rpcx" {
			g.P(serverMethodNameForRpcx(g, m))
		} else {
			g.P(serverMethodName(g, m))
		}
	}
	g.P("}")
	g.P()
	// register http server handler
	if s.Deprecated {
		g.P(deprecationComment)
	}
	g.P("func Register", s.ServiceType, "HTTPServer(g *", g.QualifiedGoIdent(ginPackage.Ident("RouterGroup")), ", srv ", serverInterfaceName(s.ServiceType), ") {")
	g.P(`r := g.Group("")`)
	g.P("{")
	for _, m := range s.Methods {
		g.P("r.", m.Method, `("`, m.Path, `", `, serverHandlerMethodName(s.ServiceType, m), "(srv))")
	}
	g.P("}")
	g.P("}")
	g.P()
	// handler
	for _, m := range s.Methods {
		if m.Deprecated {
			g.P(deprecationComment)
		}
		g.P("func ", serverHandlerMethodName(s.ServiceType, m), "(srv ", s.ServiceType, "HTTPServer", ") ", g.QualifiedGoIdent(ginPackage.Ident("HandlerFunc")), " {")
		{ // gin.HandleFunc closure
			g.P("return func(c *", g.QualifiedGoIdent(ginPackage.Ident("Context")), ") {")
			g.P("carrier := ", g.QualifiedGoIdent(transportHttpPackage.Ident("FromCarrier")), "(c.Request.Context())")
			if s.UseEncoding && m.HasVars {
				g.P("c.Request = carrier.WithValueUri(c.Request, c.Params)")
			}
			{ // binding
				g.P("shouldBind := func(req *", m.Request, ") error {")
				if s.UseEncoding {
					if m.HasBody {
						g.P("if err := carrier.Bind(c, req", m.Body, "); err != nil {")
						g.P("return err")
						g.P("}")
						if m.Body != "" {
							g.P("if err := carrier.BindQuery(c, req); err != nil {")
							g.P("return err")
							g.P("}")
						}
					} else {
						if m.Method != "PATCH" {
							g.P("if err := carrier.BindQuery(c, req", m.Body, "); err != nil {")
							g.P("return err")
							g.P("}")
						}
					}
					if m.HasVars {
						g.P("if err := carrier.BindUri(c, req); err != nil {")
						g.P("return err")
						g.P("}")
					}
				} else {
					if m.HasBody {
						g.P("if err := c.ShouldBind(req", m.Body, "); err != nil {")
						g.P("return err")
						g.P("}")
						if m.Body != "" {
							g.P("if err := c.ShouldBindQuery(req); err != nil {")
							g.P("return err")
							g.P("}")
						}
					} else {
						if m.Method != "PATCH" {
							g.P("if err := c.ShouldBindQuery(req", m.Body, "); err != nil {")
							g.P("return err")
							g.P("}")
						}
					}
					if m.HasVars {
						g.P("if err := c.ShouldBindUri(req); err != nil {")
						g.P("return err")
						g.P("}")
					}
				}
				g.P("return carrier.Validate(c.Request.Context(), req)")
				g.P("}")
			}
			g.P()
			{ // done
				g.P("var err error")
				g.P("var req ", m.Request)
				if s.RpcMode == "rpcx" {
					g.P("var reply *", m.Reply, "= new(", m.Reply, ")")
				} else {
					g.P("var reply *", m.Reply)
				}
				g.P()
				g.P("if err = shouldBind(&req); err != nil {")
				if args.DisableErrorBadRequest {
					g.P("carrier.Error(c, err)")
				} else {
					g.P("carrier.ErrorBadRequest(c, err)")
				}
				g.P("return")
				g.P("}")
				if s.RpcMode == "rpcx" {
					g.P("err = srv.", m.Name, "(c.Request.Context(), &req, reply)")
				} else {
					g.P("reply, err = srv.", m.Name, "(c.Request.Context(), &req)")
				}
				g.P("if err != nil {")
				g.P("carrier.Error(c, err)")
				g.P("return")
				g.P("}")
				g.P("carrier.Render(c, reply", m.ResponseBody, ")")
			}
			g.P("}")
		}
		g.P("}")
		g.P()
	}

	return nil
}

func serverInterfaceName(serverType string) string {
	return serverType + "HTTPServer"
}

func serverMethodNameForRpcx(g *protogen.GeneratedFile, m *methodDesc) string {
	return m.Name + "(" + g.QualifiedGoIdent(contextPackage.Ident("Context")) + ", *" + m.Request + ", *" + m.Reply + ")" + "error"
}

func serverMethodName(g *protogen.GeneratedFile, m *methodDesc) string {
	return m.Name + "(" + g.QualifiedGoIdent(contextPackage.Ident("Context")) + ", *" + m.Request + ") (*" + m.Reply + ", error)"
}

func serverHandlerMethodName(serverType string, m *methodDesc) string {
	return "_" + serverType + "_" + m.Name + strconv.Itoa(m.Num) + "_HTTP_Handler"
}
