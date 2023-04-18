package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zmicro-team/zmicro/core/errors"
)

var disableBindValidation bool
var defaultValidator = func() *validator.Validate {
	v := validator.New()
	v.SetTagName("binding")
	return v
}()

// Deprecated: use Carrier interface
func DisableBindValidation() {
	disableBindValidation = true
}

// Deprecated: use Carrier interface
func Validator() *validator.Validate {
	return defaultValidator
}

// Deprecated: use Carrier interface
func Validate(ctx context.Context, v any) error {
	if disableBindValidation {
		return nil
	}
	return defaultValidator.StructCtx(ctx, v)
}

// Deprecated: use Carrier interface
func ErrorEncoder(c *gin.Context, err error, isBadRequest bool) {
	if isBadRequest {
		err = errors.ErrBadRequest(err.Error())
	}
	Error(c, err)
}
