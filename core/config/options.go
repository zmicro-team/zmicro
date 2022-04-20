package config

import "context"

type Options struct {
	Type    string
	Path    string
	Context context.Context
}

type Option func(o *Options)

func Type(t string) Option {
	return func(o *Options) {
		o.Type = t
	}
}

func Path(p string) Option {
	return func(o *Options) {
		o.Path = p
	}
}

func SetOption(k, v interface{}) Option {
	return func(opts *Options) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, k, v)
	}
}
