package server

import (
	"github.com/smallnest/rpcx/server"
)

type Options struct {
	InitRpcServer InitRpcServerFunc
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{}

	for _, o := range opts {
		o(&options)
	}

	return options
}

type InitRpcServerFunc func(s *server.Server) error

func WithInitRpcServer(f InitRpcServerFunc) Option {
	return func(o *Options) {
		o.InitRpcServer = f
	}
}
