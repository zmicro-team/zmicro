package http

import "github.com/gin-gonic/gin"

type Options struct {
	InitHttpServer InitHttpServerFunc
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

func WithInitHttpServer(f InitHttpServerFunc) Option {
	return func(o *Options) {
		o.InitHttpServer = f
	}
}
