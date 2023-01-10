package http

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/zmicro-team/zmicro/core/encoding"
	"github.com/zmicro-team/zmicro/core/encoding/codec"
)

var globalEncoding = encoding.New()

func RegisterMarshaler(mime string, marshaler codec.Marshaler) error {
	return globalEncoding.Register(mime, marshaler)
}
func GetMarshaler(mime string) codec.Marshaler {
	return globalEncoding.Get(mime)
}
func DeleteMarshaler(mime string) error {
	return globalEncoding.Delete(mime)
}
func InboundMarshalerForRequest(req *http.Request) (string, codec.Marshaler) {
	return globalEncoding.InboundForRequest(req)
}
func OutboundMarshalerForRequest(req *http.Request) codec.Marshaler {
	return globalEncoding.OutboundForRequest(req)
}
func Bind(c *gin.Context, v any) error {
	return globalEncoding.Bind(c.Request, v)
}
func BindQuery(c *gin.Context, v any) error {
	return globalEncoding.BindQuery(c.Request, v)
}
func BindUri(c *gin.Context, v any) error {
	return globalEncoding.BindUri(c.Request, v)
}
func Render(c *gin.Context, statusCode int, v any) {
	c.Writer.WriteHeader(statusCode)
	err := globalEncoding.Render(c.Writer, c.Request, v)
	if err != nil {
		c.String(500, "Render failed cause by %v", err)
	}
}

func RequestWithUri(req *http.Request, params gin.Params) *http.Request {
	vars := make(url.Values, len(params))
	for _, p := range params {
		vars.Set(p.Key, p.Value)
	}
	return encoding.RequestWithUri(req, vars)
}
