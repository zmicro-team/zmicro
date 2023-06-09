package main

import (
	"embed"
	"io"
	"text/template"
)

//go:embed http.tpl
var Static embed.FS

// Deprecated
var ginHttpTemplate = template.Must(template.New("components").ParseFS(Static, "http.tpl")).
	Lookup("http.tpl")

// Deprecated
func (s *serviceDesc) execute(w io.Writer) error {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	return ginHttpTemplate.Execute(w, s)
}
