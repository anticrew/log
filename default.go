package log

import (
	"sync/atomic"
)

var (
	def atomic.Pointer[Logger]
)

func SetDefault(l Logger) {
	l = l.WithSkip(1)
	def.Store(&l)
}

func Default() Logger {
	l := def.Load()
	if l == nil {
		return nil
	}

	return *l
}

func Trace(msg string, attrs ...Attr) {
	Default().Trace(msg, attrs...)
}

func Debug(msg string, attrs ...Attr) {
	Default().Debug(msg, attrs...)
}

func Info(msg string, attrs ...Attr) {
	Default().Info(msg, attrs...)
}

func Warn(msg string, attrs ...Attr) {
	Default().Warn(msg, attrs...)
}

func Error(err error, message string, attrs ...Attr) {
	Default().Error(err, message, attrs...)
}

func Fatal(err error, message string, attrs ...Attr) {
	Default().Fatal(err, message, attrs...)
}

func Write(level Level, msg string, attrs ...Attr) {
	Default().Write(level, msg, attrs...)
}
