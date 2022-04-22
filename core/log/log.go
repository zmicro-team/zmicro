package log

import (
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	DebugLevel  Level = zap.DebugLevel
	InfoLevel   Level = zap.InfoLevel
	WarnLevel   Level = zap.WarnLevel
	ErrorLevel  Level = zap.ErrorLevel
	DPanicLevel Level = zap.DPanicLevel
	PanicLevel  Level = zap.PanicLevel
	FatalLevel  Level = zap.FatalLevel
)

type Field = zap.Field

func (l *Logger) Debug(msg string, fields ...Field) {
	l.l.Debug(msg, fields...)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.l.Info(msg, fields...)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.l.Warn(msg, fields...)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.l.Error(msg, fields...)
}

func (l *Logger) DPanic(msg string, fields ...Field) {
	l.l.DPanic(msg, fields...)
}

func (l *Logger) Panic(msg string, fields ...Field) {
	l.l.Panic(msg, fields...)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.l.Fatal(msg, fields...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {
	l.l.Sugar().Debugf(template, args...)
}

func (l *Logger) Infof(template string, args ...interface{}) {
	l.l.Sugar().Infof(template, args...)
}

func (l *Logger) Warnf(template string, args ...interface{}) {
	l.l.Sugar().Warnf(template, args...)
}

func (l *Logger) Errorf(template string, args ...interface{}) {
	l.l.Sugar().Errorf(template, args...)
}

func (l *Logger) DPanicf(template string, args ...interface{}) {
	l.l.Sugar().DPanicf(template, args...)
}

func (l *Logger) Panicf(template string, args ...interface{}) {
	l.l.Sugar().Panicf(template, args...)
}

func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.l.Sugar().Fatalf(template, args...)
}

var (
	Skip        = zap.Skip
	Binary      = zap.Binary
	Bool        = zap.Bool
	Boolp       = zap.Boolp
	ByteString  = zap.ByteString
	Complex128  = zap.Complex128
	Complex128p = zap.Complex128p
	Complex64   = zap.Complex64
	Complex64p  = zap.Complex64p
	Float64     = zap.Float64
	Float64p    = zap.Float64p
	Float32     = zap.Float32
	Float32p    = zap.Float32p
	Int         = zap.Int
	Intp        = zap.Intp
	Int64       = zap.Int64
	Int64p      = zap.Int64p
	Int32       = zap.Int32
	Int32p      = zap.Int32p
	Int16       = zap.Int16
	Int16p      = zap.Int16p
	Int8        = zap.Int8
	Int8p       = zap.Int8p
	String      = zap.String
	Stringp     = zap.Stringp
	Uint        = zap.Uint
	Uintp       = zap.Uintp
	Uint64      = zap.Uint64
	Uint64p     = zap.Uint64p
	Uint32      = zap.Uint32
	Uint32p     = zap.Uint32p
	Uint16      = zap.Uint16
	Uint16p     = zap.Uint16p
	Uint8       = zap.Uint8
	Uint8p      = zap.Uint8p
	Uintptr     = zap.Uintptr
	Uintptrp    = zap.Uintptrp
	Reflect     = zap.Reflect
	Namespace   = zap.Namespace
	Stringer    = zap.Stringer
	Time        = zap.Time
	Timep       = zap.Timep
	Stack       = zap.Stack
	StackSkip   = zap.StackSkip
	Duration    = zap.Duration
	Durationp   = zap.Durationp
	Any         = zap.Any

	Debug   = defaultLogger.Debug
	Info    = defaultLogger.Info
	Warn    = defaultLogger.Warn
	Error   = defaultLogger.Error
	DPanic  = defaultLogger.DPanic
	Panic   = defaultLogger.Panic
	Fatal   = defaultLogger.Fatal
	Debugf  = defaultLogger.Debugf
	Infof   = defaultLogger.Infof
	Warnf   = defaultLogger.Warnf
	Errorf  = defaultLogger.Errorf
	DPanicf = defaultLogger.DPanicf
	Panicf  = defaultLogger.Panicf
	Fatalf  = defaultLogger.Fatalf

	SetLevel = defaultLogger.SetLevel
)

func ResetDefault(l *Logger) {
	defaultLogger = l
	
	Info = defaultLogger.Info
	Warn = defaultLogger.Warn
	Error = defaultLogger.Error
	DPanic = defaultLogger.DPanic
	Panic = defaultLogger.Panic
	Fatal = defaultLogger.Fatal
	Debug = defaultLogger.Debug
	Infof = defaultLogger.Infof
	Warnf = defaultLogger.Warnf
	Errorf = defaultLogger.Errorf
	DPanicf = defaultLogger.DPanicf
	Panicf = defaultLogger.Panicf
	Fatalf = defaultLogger.Fatalf
	Debugf = defaultLogger.Debugf

	SetLevel = defaultLogger.SetLevel
}

type Logger struct {
	l  *zap.Logger
	lv *zap.AtomicLevel
}

var defaultLogger = New(os.Stderr, InfoLevel, WithCaller(true))

func Default() *Logger {
	return defaultLogger
}

type Option = zap.Option

var (
	WithCaller    = zap.WithCaller
	AddStacktrace = zap.AddStacktrace
)

func NewTee(writers []io.Writer, level Level, opts ...Option) *Logger {
	var cores []zapcore.Core
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}

	logger := &Logger{}
	lv := zap.NewAtomicLevelAt(level)
	logger.lv = &lv

	for _, w := range writers {
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg.EncoderConfig),
			zapcore.AddSync(w),
			lv,
		)
		cores = append(cores, core)
	}

	logger.l = zap.New(zapcore.NewTee(cores...), opts...)

	return logger
}

func New(writer io.Writer, level Level, opts ...Option) *Logger {
	if writer == nil {
		panic("the writer is nil")
	}
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02T15:04:05.000Z0700"))
	}

	logger := &Logger{}
	lv := zap.NewAtomicLevelAt(level)
	logger.lv = &lv
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(writer),
		lv,
	)
	logger.l = zap.New(core, opts...)

	return logger
}

func (l *Logger) Sync() error {
	return l.l.Sync()
}

func (l *Logger) SetLevel(lv Level) {
	l.lv.SetLevel(lv)
}

func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Sync()
	}
	return nil
}
