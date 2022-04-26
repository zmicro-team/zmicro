package file

import (
	"io"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultFilename = "logs/app.log"
)

type Options struct {
	Filename   string
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Compress   bool
}

type Option func(*Options)

func NewWriter(opts ...Option) io.Writer {
	options := Options{
		Filename: defaultFilename,
	}

	for _, opt := range opts {
		opt(&options)
	}
	return &lumberjack.Logger{
		Filename:   options.Filename,
		MaxSize:    options.MaxSize,
		MaxAge:     options.MaxAge,
		MaxBackups: options.MaxBackups,
		Compress:   options.Compress,
	}
}

func Filename(n string) Option {
	return func(o *Options) {
		o.Filename = n
	}
}

func MaxSize(n int) Option {
	return func(o *Options) {
		o.MaxSize = n
	}
}

func MaxAge(n int) Option {
	return func(o *Options) {
		o.MaxAge = n
	}
}

func MaxBackups(sz int) Option {
	return func(o *Options) {
		o.MaxBackups = sz
	}
}

func Compress() Option {
	return func(o *Options) {
		o.Compress = true
	}
}
