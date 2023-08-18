package main

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
)

var _ = executeClientDesc

func executeClientDesc(g *protogen.GeneratedFile, s *serviceDesc) (err error) {
	methodSets := make(map[string]struct{})

	// client interface defined
	if s.Deprecated {
		g.P(deprecationComment)
	}
	g.P("type ", clientInterfaceName(s.ServiceType), " interface {")
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
			g.P(clientMethodNameForRpcx(g, m))
		} else {
			g.P(clientMethodName(g, m))
		}
	}
	g.P("}")
	g.P()
	// client struct defined
	g.P(clientStruct(s.ServiceType))
	g.P()
	g.P(clientStructOptions(s.ServiceType))
	g.P()
	g.P(newClientStruct(s.ServiceType))
	g.P()
	methodSets = make(map[string]struct{})
	for _, m := range s.Methods {
		if _, ok := methodSets[m.Name]; ok {
			continue
		}
		methodSets[m.Name] = struct{}{}
		if m.Deprecated {
			g.P(deprecationComment)
		}
		if s.RpcMode == "rpcx" {
			g.P(clientMethodForRpcx(g, m, s.ServiceType))
		} else {
			g.P(clientMethod(g, m, s.ServiceType))
		}
		g.P()
	}
	g.P(registerClient(s))
	return
}

func clientInterfaceName(serverType string) string {
	return "I" + serverType + "Client"
}

func clientStructName(serverType string) string {
	return serverType + "Client"
}

func clientStructOtpionsName(serverType string) string {
	return serverType + "ClientOptions"
	// return "ClientOptions"
}

func clientStructOptions(serverType string) string {
	result := ""
	typeName := clientStructOtpionsName(serverType)
	result = fmt.Sprintf(`type %s struct{
		EnableValidation bool
		Validate func(context.Context, any) error
	}`, typeName)
	return result
}

func clientStruct(serverType string) string {
	result := ""
	itypeName := clientInterfaceName(serverType)
	typeName := clientStructName(serverType)
	result = fmt.Sprintf(`type %s struct{
		cc %s
		options %s
	}`, typeName, itypeName,
		clientStructOtpionsName(serverType))
	return result
}

func newClientStruct(serverType string) string {
	result := ""
	itypeName := clientInterfaceName(serverType)
	typeName := clientStructName(serverType)
	optionsName := clientStructOtpionsName(serverType)
	result = fmt.Sprintf(`func New%s (cc %s, options %s) %s {
	return &%s{cc: cc, options: options}
	}`, typeName, itypeName, optionsName, itypeName, typeName)
	return result
}

func clientMethodName(g *protogen.GeneratedFile, m *methodDesc) string {
	return m.Name + "(" + g.QualifiedGoIdent(contextPackage.Ident("Context")) + ", *" + m.Request + ") (*" + m.Reply + ", error)"
}

func clientMethodNameForRpcx(g *protogen.GeneratedFile, m *methodDesc) string {
	return m.Name + "(" + g.QualifiedGoIdent(contextPackage.Ident("Context")) + ", *" + m.Request + ", *" + m.Reply + ")" + "error"
}

func clientMethod(g *protogen.GeneratedFile, m *methodDesc, serverType string) string {
	result := strings.Builder{}
	result.WriteString("func (")
	result.WriteString("cli *" + clientStructName(serverType))
	result.WriteString(") ")
	result.WriteString(m.Name)
	result.WriteString("(")
	result.WriteString("ctx ")
	result.WriteString(g.QualifiedGoIdent(contextPackage.Ident("Context")))
	result.WriteString(", ")
	result.WriteString(fmt.Sprintf("req *%s ", m.Request))
	result.WriteString(") ")
	result.WriteString("(reply *")
	result.WriteString(m.Reply)
	result.WriteString(", ")
	result.WriteString("err error ")
	result.WriteString(") ")
	result.WriteString(" {" + "\n")
	result.WriteString(" if cli.options.EnableValidation {\n ")
	result.WriteString("if err = cli.options.Validate(ctx, req); err != nil {\n return nil, err \n}\n")
	result.WriteString("}\n")
	result.WriteString("return cli.cc.")
	result.WriteString(m.Name)
	result.WriteString("(ctx, req)")
	result.WriteString("}")
	return result.String()
}

func clientMethodForRpcx(g *protogen.GeneratedFile, m *methodDesc, serverType string) string {
	result := strings.Builder{}
	result.WriteString("func (")
	result.WriteString("cli *" + clientStructName(serverType))
	result.WriteString(") ")
	result.WriteString(m.Name)
	result.WriteString("(")
	result.WriteString("ctx ")
	result.WriteString(g.QualifiedGoIdent(contextPackage.Ident("Context")))
	result.WriteString(", ")
	result.WriteString(fmt.Sprintf("req *%s, ", m.Request))
	result.WriteString(fmt.Sprintf("reply *%s ", m.Reply))
	result.WriteString(") ")
	result.WriteString("(err error) ")
	result.WriteString(" {" + "\n")
	result.WriteString(" if cli.options.EnableValidation {\n ")
	result.WriteString("if err = cli.options.Validate(ctx, req); err != nil {\n return err \n}\n")
	result.WriteString("}\n")
	result.WriteString("return cli.cc.")
	result.WriteString(m.Name)
	result.WriteString("(ctx, req, reply)")
	result.WriteString("}")
	return result.String()
}

func registerClient(s *serviceDesc) string {
	result := strings.Builder{}

	result.WriteString("var Inner" + clientStructName(s.ServiceType) + " func () " + clientInterfaceName(s.ServiceType) + " ")
	result.WriteString("\n")
	result.WriteString("func Register" + clientStructName(s.ServiceType) + "(cc " + clientInterfaceName(s.ServiceType) + ", options " + clientStructOtpionsName(s.ServiceType) + ") {")
	result.WriteString("\n")
	result.WriteString("if " + "Inner" + clientStructName(s.ServiceType) + "!= nil {")
	result.WriteString("\n")
	result.WriteString("panic(\"client already registered\"+ \" " + clientInterfaceName(s.ServiceType) + "\")")
	result.WriteString("\n")
	result.WriteString("}")
	result.WriteString("\n")
	result.WriteString("Inner" + clientStructName(s.ServiceType) + " = func() " + clientInterfaceName(s.ServiceType) + " {")
	result.WriteString("\n")
	result.WriteString("return New" + clientStructName(s.ServiceType) + "(cc, options)")
	result.WriteString("\n")
	result.WriteString("}")
	result.WriteString("\n")
	result.WriteString("}")
	return result.String()
}
