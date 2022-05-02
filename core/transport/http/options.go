package http

import "github.com/gin-gonic/gin"

type Options struct {
	Name           string
	Addr           string
	InitHttpServer InitHttpServerFunc
	Mode           string
	Tracing        bool
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func Name(s string) Option {
	return func(o *Options) {
		o.Name = s
	}
}

func Addr(s string) Option {
	return func(o *Options) {
		o.Addr = s
	}
}

type InitHttpServerFunc func(r *gin.Engine) error

func InitHttpServer(f InitHttpServerFunc) Option {
	return func(o *Options) {
		o.InitHttpServer = f
	}
}

func Mode(s string) Option {
	return func(o *Options) {
		o.Mode = s
	}
}

func Tracing(b bool) Option {
	return func(o *Options) {
		o.Tracing = b
	}
}
