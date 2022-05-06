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
