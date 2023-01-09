{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
{{$useCustomResp := .UseCustomResponse}}
{{$rpcMode := .RpcMode}}
{{$allowFromAPI := .AllowFromAPI}}
{{$useEncoding := .UseEncoding}}
type {{$svrType}}HTTPServer interface {
{{- range .MethodSets}}
	{{.Comment}}
{{- if eq $rpcMode "rpcx"}}
	{{.Name}}(context.Context, *{{.Request}}, *{{.Reply}}) error
{{- else}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
{{- end}}
{{- if not $useEncoding}}
	// Validate the request.
    Validate(context.Context, any) error
{{- end}}
	// ErrorEncoder encode error response.
	ErrorEncoder(c *gin.Context, err error, isBadRequest bool)
{{- if $useCustomResp}}
	// ResponseEncoder encode response.
	ResponseEncoder(c *gin.Context, v any)
{{- end}}
}

type Unimplemented{{$svrType}}HTTPServer struct {}

{{- range .MethodSets}}
{{- if eq $rpcMode "rpcx"}}
func (*Unimplemented{{$svrType}}HTTPServer) {{.Name}}(context.Context, *{{.Request}}, *{{.Reply}}) error {
	return errors.New("method {{.Name}} not implemented")
}
{{- else}}
func (*Unimplemented{{$svrType}}HTTPServer) {{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error) {
	return nil, errors.New("method {{.Name}} not implemented")
}
{{- end}}
{{- end}}
{{- if not $useEncoding}}
func (*Unimplemented{{$svrType}}HTTPServer) Validate(context.Context, any) error { return nil }
{{- end}}
func (*Unimplemented{{$svrType}}HTTPServer) ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	var code = 500
	if isBadRequest {
		code = 400
	}
	c.String(code, err.Error())
}
{{- if $useCustomResp}}
func (*Unimplemented{{$svrType}}HTTPServer) ResponseEncoder(c *gin.Context, v any) {
	c.JSON(200, v)
}
{{- end}}

func Register{{$svrType}}HTTPServer(g *gin.RouterGroup, srv {{$svrType}}HTTPServer) {
	r := g.Group("")
	{{- range .Methods}}
	r.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPServer) gin.HandlerFunc {
	return func(c *gin.Context) {
		{{- if and $useEncoding .HasVars}}
		c.Request = http.RequestWithUri(c.Request, c.Params)
		{{- end}}
		shouldBind := func(req *{{.Request}}) error {
		    {{- if $useEncoding}}
		    {{- if .HasBody}}
			if err := http.Bind(c, req{{.Body}}); err != nil {
				return err
			}
			{{- if not (eq .Body "")}}
			if err := http.BindQuery(c, req); err != nil {
				return err
			}
			{{- end}}
			{{- else}}
			{{- if not (eq .Method "PATCH")}}
			if err := http.BindQuery(c, req{{.Body}}); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			{{- if .HasVars}}
			if err := http.BindUri(c, req); err != nil {
				return err
			}
			{{- end}}
			return http.Validate(c.Request.Context(), req)
		    {{- else}}
			{{- if .HasBody}}
			if err := c.ShouldBind(req{{.Body}}); err != nil {
				return err
			}
			{{- if not (eq .Body "")}}
			if err := c.ShouldBindQuery(req); err != nil {
				return err
			}
			{{- end}}
			{{- else}}
			{{- if not (eq .Method "PATCH")}}
			if err := c.ShouldBindQuery(req{{.Body}}); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			{{- if .HasVars}}
			if err := c.ShouldBindUri(req); err != nil {
				return err
			}
			{{- end}}
			return srv.Validate(c.Request.Context(), req)
			{{- end}}
		}

		var err error
		var req {{.Request}}
		var rsp *{{.Reply}} {{- if eq $rpcMode "rpcx"}}= new({{.Reply}}){{- end}}

		if err = shouldBind(&req); err != nil {
			srv.ErrorEncoder(c, err, true)
			return
		}
		{{- if eq $rpcMode "rpcx"}}
		err = srv.{{.Name}}(c.Request.Context(), &req, rsp)
		{{- else}}
		rsp, err = srv.{{.Name}}(c.Request.Context(), &req)
		{{- end}}
		if err != nil {
			srv.ErrorEncoder(c, err, false)
			return
		}
		{{- if $useCustomResp}}
		srv.ResponseEncoder(c, rsp{{.ResponseBody}})
		{{- else}}
	    {{- if $useEncoding}}
	    http.Render(c, 200, rsp{{.ResponseBody}})
        {{- else}}
		c.JSON(200, rsp{{.ResponseBody}})
		{{- end}}
		{{- end}}
	}
}
{{end}}

{{- if $allowFromAPI}}
type From{{$svrType}}HTTPServer interface {
{{- range .MethodSets}}
	{{.Comment}}
{{- if eq $rpcMode "rpcx"}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- else}}
	{{.Name}}(context.Context, *{{.Request}}, *{{.Reply}}) error
{{- end}}
{{- end}}
{{- if not $useEncoding}}
    Validate(context.Context, any) error
{{- end}}
	ErrorEncoder(c *gin.Context, err error, isBadRequest bool)
{{- if $useCustomResp}}
	ResponseEncoder(c *gin.Context, v any)
{{- end}}
}

type From{{$svrType}} struct {
	From{{$svrType}}HTTPServer
}

func NewFrom{{$svrType}}HTTPServer(from From{{$svrType}}HTTPServer) {{$svrType}}HTTPServer {
	return &From{{$svrType}}{from}
}

{{- range .MethodSets}}
{{- if eq $rpcMode "rpcx"}}
func (f *From{{$svrType}}) {{.Name}}(ctx context.Context, req *{{.Request}}, rsp *{{.Reply}}) error {
	result, err := f.From{{$svrType}}HTTPServer.{{.Name}}(ctx, req)
	if err != nil {
		return err
	}
	if result == nil {
		*rsp = {{.Reply}}{}
	} else {
		*rsp = *result
	}
	return nil
}
{{- else}}
func (f *From{{$svrType}}) {{.Name}}(ctx context.Context, req *{{.Request}}) (*{{.Reply}}, error) {
	var err error
	var rsp {{.Reply}}

	err = f.From{{$svrType}}HTTPServer.{{.Name}}(ctx, req, &rsp)
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}
{{- end}}
{{- end}}
{{- end}}