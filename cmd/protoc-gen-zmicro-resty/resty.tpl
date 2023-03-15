{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{$svrType}}HTTPClient interface {
{{- range .Methods}}
	{{.Comment}}
{{- if eq .Num 0}}
	{{.Name}}(context.Context, *{{.Request}}, ...http.CallOption) (*{{.Reply}}, error)
{{- else}}
	{{.Name}}_{{.Num}}(context.Context, *{{.Request}}, ...http.CallOption) (*{{.Reply}}, error)
{{- end}}
{{- end}}
}

type {{$svrType}}HTTPClientImpl struct {
	cc *http.Client
}

func New{{$svrType}}HTTPClient(c *http.Client) {{$svrType}}HTTPClient {
	return &{{$svrType}}HTTPClientImpl{
		cc: c,
	}
}

{{range .Methods}}
{{- if eq .Num 0}}
func (c *{{$svrType}}HTTPClientImpl){{.Name}}(ctx context.Context, req *{{.Request}}, opts ...http.CallOption) (*{{.Reply}}, error) {
{{- else}}
func (c *{{$svrType}}HTTPClientImpl){{.Name}}_{{.Num}}(ctx context.Context, req *{{.Request}}, opts ...http.CallOption) (*{{.Reply}}, error) {
{{- end}}
	var err error
	var resp {{.Reply}}

	settings := http.DefaultCallOption("{{.Path}}")
	for _, opt := range opts {
		opt(&settings)
	}
	{{- if .HasVars}}
	path := c.cc.EncodeURL(settings.Path, req, {{not .HasBody}})
	{{- else}}
	{{- if not .HasBody}}
	var query string

	query, err = c.cc.EncodeQuery(req)
	if err != nil {
		return nil, err
	}
	path := settings.Path
	if query != "" {
		path += "?" + query
	}
	{{- else}}
	path := settings.Path
	{{- end}}
	{{- end}}
	ctx = http.WithValueCallOption(ctx, settings)
	err = c.cc.Invoke(ctx, "{{.Method}}", path, {{if .HasBody -}}req{{.Body}}{{- else}}nil{{- end}}, &resp{{.ResponseBody}})
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
{{end}}