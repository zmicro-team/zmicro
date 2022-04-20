package zmicro

import "github.com/iobrother/zmicro/core/transport/rpc/server"

type Options struct {
	InitRpcServer server.InitRpcServerFunc
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

func WithInitRpcServer(f server.InitRpcServerFunc) Option {
	return func(o *Options) {
		o.InitRpcServer = f
	}
}
