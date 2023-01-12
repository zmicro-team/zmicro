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
	// Validate the request.
    Validate(context.Context, any) error
	// ErrorEncoder encode error response.
	ErrorEncoder(c *gin.Context, err error, isBadRequest bool)
{{- if $useEncoding}}
    // Bind checks the Method and Content-Type to select codec.Marshaler automatically,
    // Depending on the "Content-Type" header different bind are used.
    Bind(c *gin.Context, v any) error
    // BindQuery binds the passed struct pointer using the query codec.Marshaler.
    BindQuery(c *gin.Context, v any) error
    // BindUri binds the passed struct pointer using the uri codec.Marshaler.
    // NOTE: before use this, you should set uri params in the request context with RequestWithUri.
    BindUri(c *gin.Context, v any) error
    // RequestWithUri sets the URL params for the given request.
    RequestWithUri(req *http.Request, params gin.Params) *http.Request
    // Render encode response.
    Render(c *gin.Context, v any)
{{- else}}
{{- if $useCustomResp}}
	// Render encode response.
	Render(c *gin.Context, v any)
{{- end}}
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
func (*Unimplemented{{$svrType}}HTTPServer) Validate(context.Context, any) error { return nil }
func (*Unimplemented{{$svrType}}HTTPServer) ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	var code = 500
	if isBadRequest {
		code = 400
	}
	c.String(code, err.Error())
}
{{- if $useEncoding}}
func (*Unimplemented{{$svrType}}HTTPServer) Bind(c *gin.Context, v any) error {
    return c.ShouldBind(v)
}
func (*Unimplemented{{$svrType}}HTTPServer) BindQuery(c *gin.Context, v any) error {
    return c.ShouldBindQuery(v)
}
func (*Unimplemented{{$svrType}}HTTPServer) BindUri(c *gin.Context, v any) error {
    return c.ShouldBindUri(v)
}
func (*Unimplemented{{$svrType}}HTTPServer) RequestWithUri(req *http.Request, _ gin.Params) *http.Request {
    return req
}
func (*Unimplemented{{$svrType}}HTTPServer) Render(c *gin.Context, v any) {
    c.JSON(200, v)
}
{{- else}}
{{- if $useCustomResp}}
func (*Unimplemented{{$svrType}}HTTPServer) Render(c *gin.Context, v any) {
	c.JSON(200, v)
}
{{- end}}
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
		c.Request = srv.RequestWithUri(c.Request, c.Params)
		{{- end}}
		shouldBind := func(req *{{.Request}}) error {
		    {{- if $useEncoding}}
		    {{- if .HasBody}}
			if err := srv.Bind(c, req{{.Body}}); err != nil {
				return err
			}
			{{- if not (eq .Body "")}}
			if err := srv.BindQuery(c, req); err != nil {
				return err
			}
			{{- end}}
			{{- else}}
			{{- if not (eq .Method "PATCH")}}
			if err := srv.BindQuery(c, req{{.Body}}); err != nil {
				return err
			}
			{{- end}}
			{{- end}}
			{{- if .HasVars}}
			if err := srv.BindUri(c, req); err != nil {
				return err
			}
			{{- end}}
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
			{{- end}}
			return srv.Validate(c.Request.Context(), req)
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
		{{- if or $useEncoding $useCustomResp}}
		srv.Render(c, rsp{{.ResponseBody}})
        {{- else}}
        c.JSON(200, rsp{{.ResponseBody}})
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
    Validate(context.Context, any) error
	ErrorEncoder(c *gin.Context, err error, isBadRequest bool)
{{- if $useEncoding}}
    // Bind checks the Method and Content-Type to select codec.Marshaler automatically,
    // Depending on the "Content-Type" header different bind are used.
    Bind(c *gin.Context, v any) error
    // BindQuery binds the passed struct pointer using the query codec.Marshaler.
    BindQuery(c *gin.Context, v any) error
    // BindUri binds the passed struct pointer using the uri codec.Marshaler.
    // NOTE: before use this, you should set uri params in the request context with RequestWithUri.
    BindUri(c *gin.Context, v any) error
    // RequestWithUri sets the URL params for the given request.
    RequestWithUri(req *http.Request, params gin.Params) *http.Request
    // Render encode response.
    Render(c *gin.Context, v any)
{{- else}}
{{- if $useCustomResp}}
	// Render encode response.
	Render(c *gin.Context, v any)
{{- end}}
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