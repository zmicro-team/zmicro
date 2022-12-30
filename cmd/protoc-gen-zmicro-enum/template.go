package main

import (
	"embed"
	"errors"
	"io"
	"text/template"
)

type EnumValue struct {
	Number     int
	Value      string
	CamelValue string
	Mapping    string
	Comment    string
}

type Enum struct {
	Name    string
	Comment string
	Values  []*EnumValue
}

//go:embed enum.tpl
var Static embed.FS

var TemplateFuncs = template.FuncMap{
	"add":            func(a, b int) int { return a + b },
	"snakecase":      func(s string) string { return SnakeCase(s, false) },
	"kebabcase":      func(s string) string { return Kebab(s, false) },
	"camelcase":      func(s string) string { return CamelCase(s) },
	"smallcamelcase": func(s string) string { return SmallCamelCase(s, false) },
}
var enumTemplate = template.Must(template.New("components").
	Funcs(TemplateFuncs).
	ParseFS(Static, "enum.tpl")).
	Lookup("enum.tpl")

type EnumFile struct {
	Version       string
	ProtocVersion string
	IsDeprecated  bool
	Source        string
	Package       string
	Enums         []*Enum
}

func (e *EnumFile) execute(t *template.Template, w io.Writer) error {
	return t.Execute(w, e)
}

func ParseTemplateFromFile(filename string) (*template.Template, error) {
	if filename == "" {
		return nil, errors.New("required template filename")
	}
	tt, err := template.New("custom").
		Funcs(TemplateFuncs).
		ParseFiles(filename)
	if err != nil {
		return nil, err
	}
	ts := tt.Templates()
	if len(ts) == 0 {
		return nil, errors.New("not found any template")
	}
	return ts[0], nil
}
