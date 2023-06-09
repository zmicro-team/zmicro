package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/zmicro-team/zmicro/core/encoding"
	"github.com/zmicro-team/zmicro/core/encoding/jsonpb"
	"github.com/zmicro-team/zmicro/core/errors"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ Carrier = (*Carry)(nil)

type ErrorTranslator interface {
	Translate(err error) error
}

type Carry struct {
	Validation *validator.Validate
	Encoding   *encoding.Encoding
	// translate error
	translate ErrorTranslator
}

func NewCarry() *Carry {
	e := encoding.New()
	err := e.Register(encoding.MIMEJSON, &Codec{
		Codec: &jsonpb.Codec{
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
	return &Carry{
		Validation: func() *validator.Validate {
			v := validator.New()
			v.SetTagName("binding")
			return v
		}(),
		Encoding: e,
	}
}
func (cy *Carry) SetEncoding(e *encoding.Encoding) *Carry {
	cy.Encoding = e
	return cy
}

func (cy *Carry) SetValidation(v *validator.Validate) *Carry {
	cy.Validation = v
	return cy
}

func (cy *Carry) SetTranslateError(e ErrorTranslator) *Carry {
	cy.translate = e
	return cy
}

func (*Carry) WithValueUri(req *http.Request, params gin.Params) *http.Request {
	return WithValueUri(req, params)
}
func (cy *Carry) Bind(c *gin.Context, v any) error {
	if cy.Encoding == nil {
		return c.ShouldBind(v)
	}
	return cy.Encoding.Bind(c.Request, v)
}
func (cy *Carry) BindQuery(c *gin.Context, v any) error {
	if cy.Encoding == nil {
		return c.ShouldBindQuery(v)
	}
	return cy.Encoding.BindQuery(c.Request, v)
}
func (cy *Carry) BindUri(c *gin.Context, v any) error {
	if cy.Encoding == nil {
		return c.ShouldBindUri(v)
	}
	return cy.Encoding.BindUri(c.Request, v)
}
func (*Carry) ErrorBadRequest(c *gin.Context, err error) {
	Error(c, errors.ErrBadRequest(err.Error()))
}
func (cy *Carry) Error(c *gin.Context, err error) {
	if cy.translate != nil {
		err = cy.translate.Translate(err)
	}
	Error(c, err)
}
func (cy *Carry) Render(c *gin.Context, v any) {
	if cy.Encoding == nil {
		JSON(c, v)
	}
	c.Writer.WriteHeader(http.StatusOK)
	err := cy.Encoding.Render(c.Writer, c.Request, v)
	if err != nil {
		c.String(http.StatusInternalServerError, "Render failed cause by %v", err)
	}
}
func (cy *Carry) Validator() *validator.Validate {
	return cy.Validation
}
func (cy *Carry) Validate(ctx context.Context, v any) error {
	return cy.Validation.StructCtx(ctx, v)
}
func (cy *Carry) StructCtx(ctx context.Context, v any) error {
	return cy.Validation.StructCtx(ctx, v)
}
func (cy *Carry) Struct(v any) error {
	return cy.Validation.Struct(v)
}
func (cy *Carry) VarCtx(ctx context.Context, v any, tag string) error {
	return cy.Validation.VarCtx(ctx, v, tag)
}
func (cy *Carry) Var(v any, tag string) error {
	return cy.Validation.Var(v, tag)
}
