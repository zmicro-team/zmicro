package http

import "github.com/gin-gonic/gin"

type Options struct {
	InitHttpServer InitHttpServerFunc
	Mode           string
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
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
