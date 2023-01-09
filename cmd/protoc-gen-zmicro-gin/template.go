package main

import (
	"embed"
	"io"
	"text/template"
)

//go:embed http.tpl
var Static embed.FS

var ginHttpTemplate = template.Must(template.New("components").ParseFS(Static, "http.tpl")).
	Lookup("http.tpl")

type serviceDesc struct {
	ServiceType string // Greeter
	ServiceName string // helloworld.Greeter
	Metadata    string // api/v1/helloworld.proto
	Methods     []*methodDesc
	MethodSets  map[string]*methodDesc // unique because additional_bindings

	UseCustomResponse bool
	RpcMode           string
	AllowFromAPI      bool
}

type methodDesc struct {
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

func (s *serviceDesc) execute(w io.Writer) error {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	return ginHttpTemplate.Execute(w, s)
}
