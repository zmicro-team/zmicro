package config

type Options struct {
	Type      string
	Path      string
	Callbacks []func(IConfig)
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

func Callbacks(f ...func(IConfig)) Option {
	return func(o *Options) {
		o.Callbacks = f
	}
}
