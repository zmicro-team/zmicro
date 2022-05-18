package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
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

//func Data(c *gin.Context, data any) {
//
//}

type ginKey struct{}

func NewContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginKey{}, c)
}

func FromContext(ctx context.Context) (c *gin.Context, ok bool) {
	c, ok = ctx.Value(ginKey{}).(*gin.Context)
	return
}

type Implemented struct{}

func (*Implemented) Validate(ctx context.Context, v any) error {
	return Validate(ctx, v)
}

func (*Implemented) ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	ErrorEncoder(c, err, isBadRequest)
}
