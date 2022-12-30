package main

import (
	"embed"
	"io"
	"text/template"
)

//go:embed errno.tpl
var Static embed.FS

var errnoTemplate = template.Must(template.New("components").ParseFS(Static, "errno.tpl")).
	Lookup("errno.tpl")

type errorInfo struct {
	Name       string
	Code       int
	Value      string
	CamelValue string
	Message    string
}

type errorWrapper struct {
	Errors []*errorInfo
}

func (e *errorWrapper) execute(w io.Writer) error {
	return errnoTemplate.Execute(w, e)
}
