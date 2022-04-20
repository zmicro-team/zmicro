package client

type Options struct {
	ServiceName string
	ServiceAddr string
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
