package zmicro

import (
	"github.com/iobrother/zmicro/core/transport/http"
	"github.com/iobrother/zmicro/core/transport/rpc/server"
)

type Options struct {
	InitRpcServer  server.InitRpcServerFunc
	InitHttpServer http.InitHttpServerFunc
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func InitRpcServer(f server.InitRpcServerFunc) Option {
	return func(o *Options) {
		o.InitRpcServer = f
	}
}

func InitHttpServer(f http.InitHttpServerFunc) Option {
	return func(o *Options) {
		o.InitHttpServer = f
	}
}
