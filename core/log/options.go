package log

type Option interface {
	apply(logger *Logger)
}

type optionFunc func(logger *Logger)

func (f optionFunc) apply(log *Logger) {
	f(log)
}

func WithCaller(enabled bool) Option {
	return optionFunc(func(l *Logger) {
		l.addCaller = enabled
	})
}

func WithCallerSkip(skip int) Option {
	return optionFunc(func(l *Logger) {
		l.callSkip = skip
	})
}

func Development() Option {
	return optionFunc(func(l *Logger) {
		l.development = true
	})
}

func WithDefaultValuer(vs ...Valuer) Option {
	return optionFunc(func(l *Logger) {
		if len(vs) > 0 {
			fn := make([]Valuer, 0, len(vs)+len(l.fn))
			fn = append(fn, l.fn...)
			fn = append(fn, vs...)
			l.fn = fn
		}
	})
}
