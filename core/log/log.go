package log

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Valuer is returns a log value.
type Valuer func(ctx context.Context) Field

type Logger struct {
	l           *zap.Logger
	lv          *zap.AtomicLevel
	development bool
	addCaller   bool
	callSkip    int
	fn          []Valuer
	ctx         context.Context
}

func NewTee(writers []io.Writer, level Level, opts ...Option) *Logger {
	logger := &Logger{callSkip: 1, ctx: context.Background()}
	lv := zap.NewAtomicLevelAt(level)
	logger.lv = &lv
	for _, opt := range opts {
		opt.apply(logger)
	}

	cfg := zap.NewProductionConfig()
	if logger.development {
		cfg = zap.NewDevelopmentConfig()
	}
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}
	var enc zapcore.Encoder
	if logger.development {
		enc = zapcore.NewConsoleEncoder(cfg.EncoderConfig)
	} else {
		enc = zapcore.NewJSONEncoder(cfg.EncoderConfig)
	}

	var cores []zapcore.Core
	for _, w := range writers {
		core := zapcore.NewCore(
			enc,
			zapcore.AddSync(w),
			lv,
		)
		cores = append(cores, core)
	}

	options := []zap.Option{zap.WithCaller(logger.addCaller)}
	if logger.development {
		options = append(options, zap.Development())
	}
	options = append(options, zap.AddCallerSkip(logger.callSkip))

	logger.l = zap.New(
		zapcore.NewTee(cores...),
		options...,
	)

	return logger
}

func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	return NewTee([]io.Writer{writer}, level, opts...)
}

func (l *Logger) Logger() *zap.Logger {
	return l.l
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func (l *Logger) SetLevel(lv Level) {
	l.lv.SetLevel(lv)
}

// Enabled returns true if the given level is at or above this level.
func (l *Logger) Enabled(lvl zapcore.Level) bool {
	return l.lv.Enabled(lvl)
}

// V returns true if the given level is at or above this level.
// same as Enabled
func (l *Logger) V(lvl int) bool {
	return l.lv.Enabled(zapcore.Level(lvl))
}

// WithValuer with Valuer function.
func (l *Logger) WithValuer(fs ...Valuer) *Logger {
	fn := make([]Valuer, 0, len(fs)+len(l.fn))
	fn = append(fn, l.fn...)
	fn = append(fn, fs...)
	return &Logger{
		l.l,
		l.lv,
		l.development,
		l.addCaller,
		l.callSkip,
		fn,
		l.ctx,
	}
}

// WithContext return log with inject context.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	return &Logger{
		l.l,
		l.lv,
		l.development,
		l.addCaller,
		l.callSkip,
		l.fn,
		ctx,
	}
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{
		l.l.With(fields...),
		l.lv,
		l.development,
		l.addCaller,
		l.callSkip,
		l.fn,
		l.ctx,
	}
}

// Named adds a sub-scope to the logger's name. See Log.Named for details.
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		l.l.Named(name),
		l.lv,
		l.development,
		l.addCaller,
		l.callSkip,
		l.fn,
		l.ctx,
	}
}

func (l *Logger) Debug(v any, fields ...Field) {
	if !l.lv.Enabled(DebugLevel) {
		return
	}
	l.l.Debug(fmt.Sprint(v), injectFields(l.ctx, l.fn, fields...)...)
}

func (l *Logger) Info(v any, fields ...Field) {
	if !l.lv.Enabled(InfoLevel) {
		return
	}
	l.l.Info(fmt.Sprint(v), injectFields(l.ctx, l.fn, fields...)...)
}

func (l *Logger) Warn(v any, fields ...Field) {
	if !l.lv.Enabled(WarnLevel) {
		return
	}
	l.l.Warn(fmt.Sprint(v), injectFields(l.ctx, l.fn, fields...)...)
}

func (l *Logger) Error(v any, fields ...Field) {
	if !l.lv.Enabled(ErrorLevel) {
		return
	}
	l.l.Error(fmt.Sprint(v), injectFields(l.ctx, l.fn, fields...)...)
}

func (l *Logger) DPanic(v any, fields ...Field) {
	if !l.lv.Enabled(DPanicLevel) {
		return
	}
	l.l.DPanic(fmt.Sprint(v), injectFields(l.ctx, l.fn, fields...)...)
}

func (l *Logger) Panic(v any, fields ...Field) {
	if !l.lv.Enabled(PanicLevel) {
		return
	}
	l.l.Panic(fmt.Sprint(v), injectFields(l.ctx, l.fn, fields...)...)
}

func (l *Logger) Fatal(v any, fields ...Field) {
	if !l.lv.Enabled(FatalLevel) {
		return
	}
	l.l.Fatal(fmt.Sprint(v), injectFields(l.ctx, l.fn, fields...)...)
}

func (l *Logger) Debugf(template string, args ...any) {
	if !l.lv.Enabled(DebugLevel) {
		return
	}
	l.l.With(injectFields(l.ctx, l.fn)...).Sugar().Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...any) {
	if !l.lv.Enabled(InfoLevel) {
		return
	}
	l.l.With(injectFields(l.ctx, l.fn)...).Sugar().Infof(template, args...)
}

func (l *Logger) Warnf(template string, args ...any) {
	if !l.lv.Enabled(WarnLevel) {
		return
	}
	l.l.With(injectFields(l.ctx, l.fn)...).Sugar().Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...any) {
	if !l.lv.Enabled(ErrorLevel) {
		return
	}
	l.l.With(injectFields(l.ctx, l.fn)...).Sugar().Errorf(template, args...)
}

func (l *Logger) DPanicf(template string, args ...any) {
	if !l.lv.Enabled(DPanicLevel) {
		return
	}
	l.l.With(injectFields(l.ctx, l.fn)...).Sugar().DPanicf(template, args...)
}

func (l *Logger) Panicf(template string, args ...any) {
	if !l.lv.Enabled(PanicLevel) {
		return
	}
	l.l.With(injectFields(l.ctx, l.fn)...).Sugar().Panicf(template, args...)
}

func (l *Logger) Fatalf(template string, args ...any) {
	if !l.lv.Enabled(FatalLevel) {
		return
	}
	l.l.With(injectFields(l.ctx, l.fn)...).Sugar().Fatalf(template, args...)
}

func injectFields(ctx context.Context, vs []Valuer, fd ...Field) []Field {
	var fields []Field

	switch {
	case len(vs) == 0 && len(fd) == 0:
		// do nothing
	case len(vs) > 0 && len(fd) > 0:
		fields = make([]Field, 0, len(vs)+len(fd))
		for _, f := range vs {
			fields = append(fields, f(ctx))
		}
		for _, v := range fd {
			fields = append(fields, v)
		}
	case len(vs) > 0:
		fields = make([]Field, 0, len(vs))
		for _, f := range vs {
			fields = append(fields, f(ctx))
		}
	default: // len(fd) > 0
		fields = fd
	}
	return fields
}
