package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/zmicro-team/zmicro/core/encoding"
	"github.com/zmicro-team/zmicro/core/encoding/jsonpb"
	"github.com/zmicro-team/zmicro/core/errors"
)

func Error(c *gin.Context, err error) {
	if err == nil {
		c.Status(http.StatusOK)
		return
	}

	e := errors.FromError(err)
	code := int(e.Code)
	if e.Code >= 1000 {
		code = 599
	}
	c.JSON(code, e)
	c.Abort()
}

func JSON(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

type Implemented struct {
	Encoding *encoding.Encoding
}

func NewDefaultImplemented() *Implemented {
	e := encoding.New()
	err := e.Register(encoding.MIMEJSON, &Codec{
		Marshaler: &jsonpb.Codec{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:  true,
				UseEnumNumbers: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	})
	if err != nil {
		panic(err)
	}
	return &Implemented{
		Encoding: e,
	}
}

func (i *Implemented) Validate(ctx context.Context, v any) error {
	return Validate(ctx, v)
}

func (*Implemented) ValidateMap(ctx context.Context, data map[string]interface{}, rules map[string]interface{}) map[string]interface{} {
	return Validator().ValidateMapCtx(ctx, data, rules)
}

func (*Implemented) ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	ErrorEncoder(c, err, isBadRequest)
}

func (i *Implemented) Bind(c *gin.Context, v any) error {
	if i.Encoding == nil {
		return c.ShouldBind(v)
	}
	return i.Encoding.Bind(c.Request, v)
}
func (i *Implemented) BindQuery(c *gin.Context, v any) error {
	if i.Encoding == nil {
		return c.ShouldBindQuery(v)
	}
	return i.Encoding.BindQuery(c.Request, v)
}
func (i *Implemented) BindUri(c *gin.Context, v any) error {
	if i.Encoding == nil {
		return c.ShouldBindUri(v)
	}
	return i.Encoding.BindUri(c.Request, v)
}
func (i *Implemented) RequestWithUri(req *http.Request, params gin.Params) *http.Request {
	return RequestWithUri(req, params)
}
func (i *Implemented) Render(c *gin.Context, v any) {
	if i.Encoding == nil {
		c.JSON(http.StatusOK, v)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
	err := i.Encoding.Render(c.Writer, c.Request, v)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render failed cause by %v", err)
	}
}
