package server

import (
	"github.com/smallnest/rpcx/server"
)

type Options struct {
	Name          string
	Addr          string
	InitRpcServer InitRpcServerFunc

	// registry
	BasePath       string
	UpdateInterval int
	EtcdAddr       []string

	Tracing bool
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

type InitRpcServerFunc func(s *server.Server) error

func InitRpcServer(f InitRpcServerFunc) Option {
	return func(o *Options) {
		o.InitRpcServer = f
	}
}

func BasePath(s string) Option {
	return func(o *Options) {
		o.BasePath = s
	}
}

func UpdateInterval(i int) Option {
	return func(o *Options) {
		o.UpdateInterval = i
	}
}

func EtcdAddr(a []string) Option {
	return func(o *Options) {
		o.EtcdAddr = a
	}
}

func Tracing(b bool) Option {
	return func(o *Options) {
		o.Tracing = b
	}
}
