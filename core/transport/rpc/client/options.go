package client

type Options struct {
	ServiceName string
	ServiceAddr string

	// registry
	BasePath string
	EtcdAddr []string

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

func WithServiceName(n string) Option {
	return func(opts *Options) {
		opts.ServiceName = n
	}
}

func WithServiceAddr(addr string) Option {
	return func(opts *Options) {
		opts.ServiceAddr = addr
	}
}

func BasePath(s string) Option {
	return func(o *Options) {
		o.BasePath = s
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
