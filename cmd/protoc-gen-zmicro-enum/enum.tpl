// Code generated by protoc-gen-go-gin. DO NOT EDIT.
// versions:
//   - protoc-gen-go-enum {{.Version}}
//   - protoc             {{.ProtocVersion}}
{{- if .IsDeprecated}}
// {{.Source}} is a deprecated file.
{{- else}}
// source: {{.Source}}
{{- end}}

package {{.Package}}

{{- range $e := .Enums}}
// __{{$e.Name}}Mapping {{$e.Name}} mapping
var __{{$e.Name}}Mapping = map[int]string{
{{- range $ee := $e.Values}}
	{{$ee.Number}}: "{{$ee.Mapping}}",
{{- end}}
}
// Get{{$e.Name}}Desc get mapping description
// {{$e.Comment}}
func Get{{$e.Name}}Desc(t int) string {
	return __{{$e.Name}}Mapping[t]
}
{{- end}}

