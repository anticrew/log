package log

import (
	"errors"
	"os"

	"github.com/anticrew/log/internal/caller"
)

type Logger interface {
	Skip(s int)
	Engine() Engine
	Trace(msg string, attrs ...Attr)
	Debug(msg string, attrs ...Attr)
	Info(msg string, attrs ...Attr)
	Warn(msg string, attrs ...Attr)
	Error(err error, message string, attrs ...Attr)
	Fatal(err error, message string, attrs ...Attr)
	Write(level Level, message string, attrs ...Attr)
	With(attrs ...Attr) Logger
	WithSkip(s int) Logger
}

type logger struct {
	leveler Leveler

	engine      Engine
	middlewares []Middleware
	attrs       *Attrs

	skipCount int
}

func New(engine Engine, leveler Leveler, middlewares ...Middleware) Logger {
	return newLogger(leveler, engine, NewAttrs(), middlewares)
}

func newLogger(leveler Leveler, engine Engine, attrs *Attrs, middlewares []Middleware) Logger {
	return &logger{
		leveler:     leveler,
		engine:      engine,
		middlewares: middlewares,
		attrs:       attrs,
	}
}

func (l *logger) Skip(skip int) {
	l.skipCount = skip
}

func (l *logger) WithSkip(skip int) Logger {
	newL := l.clone(nil)
	newL.Skip(l.skipCount + skip)
	return newL
}

func (l *logger) Engine() Engine {
	return l.engine
}

func (l *logger) Trace(msg string, attrs ...Attr) {
	l.write(LevelTrace, msg, nil, attrs...)
}

func (l *logger) Debug(msg string, attrs ...Attr) {
	l.write(LevelDebug, msg, nil, attrs...)
}

func (l *logger) Info(msg string, attrs ...Attr) {
	l.write(LevelInfo, msg, nil, attrs...)
}

func (l *logger) Warn(msg string, attrs ...Attr) {
	l.write(LevelWarn, msg, nil, attrs...)
}

func (l *logger) Error(err error, message string, attrs ...Attr) {
	l.write(LevelError, message, err, attrs...)
}

func (l *logger) Fatal(err error, message string, attrs ...Attr) {
	l.write(LevelFatal, message, err, attrs...)

	//nolint: revive // ожидаемое поведение
	os.Exit(1)
}

func (l *logger) Write(level Level, message string, attrs ...Attr) {
	l.write(level, message, nil, attrs...)
}

func (l *logger) write(level Level, message string, err error, attrs ...Attr) {
	if l.leveler.Level() > level {
		return
	}

	r := NewRecord(level, message, l.attrs.Clone())
	defer r.Dispose()

	r.Attrs.Append(attrs...)
	if err != nil {
		r.Attrs.Append(Err(err))
	}

	// +1 to skip xlog.(l *logger).write call, +1 to skip xlog.(l*logger).<level> call
	callerText, callerErr := caller.Take(l.skipCount + 2)
	if callerErr != nil {
		err = errors.Join(err, callerErr)
	}

	r.Attrs.Append(String(CallerKey, callerText))

	for i := len(l.middlewares); i > 0; i-- {
		mwR, mwErr := l.middlewares[i-1].Handle(r)
		err = errors.Join(err, mwErr)

		r = mwR
	}

	err = errors.Join(err, l.engine.Write(r))
	if err != nil {
		r.Level = LevelError
		r.Message = err.Error()

		_ = l.engine.Write(r)
	}
}

func (l *logger) With(attrs ...Attr) Logger {
	return l.clone(attrs)
}

func (l *logger) clone(attrs []Attr) Logger {
	return newLogger(l.leveler, l.engine, l.attrs.Clone().Append(attrs...), l.middlewares)
}
