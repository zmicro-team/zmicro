package http

import (
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
