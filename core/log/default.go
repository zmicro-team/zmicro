package log

import (
	"os"
)

var defaultLogger = New(os.Stderr, InfoLevel, WithCaller(true))

var (
	GetLogger   = defaultLogger.Logger
	SetLevel    = defaultLogger.SetLevel
	Enabled     = defaultLogger.Enabled
	V           = defaultLogger.V
	WithValuer  = defaultLogger.WithValuer
	WithContext = defaultLogger.WithContext
	With        = defaultLogger.With
	Named       = defaultLogger.Named

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
)

func ResetDefault(l *Logger) {
	defaultLogger = l

	Debug = defaultLogger.Debug
	Info = defaultLogger.Info
	Warn = defaultLogger.Warn
	Error = defaultLogger.Error
	DPanic = defaultLogger.DPanic
	Panic = defaultLogger.Panic
	Fatal = defaultLogger.Fatal
	Debugf = defaultLogger.Debugf
	Infof = defaultLogger.Infof
	Warnf = defaultLogger.Warnf
	Errorf = defaultLogger.Errorf
	DPanicf = defaultLogger.DPanicf
	Panicf = defaultLogger.Panicf
	Fatalf = defaultLogger.Fatalf

	SetLevel = defaultLogger.SetLevel
}

func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Sync()
	}
	return nil
}

func Default() *Logger {
	return defaultLogger
}
