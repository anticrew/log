package log

import (
	"context"
	"io"
	"os"
	"time"
)

const (
	LevelKey   = "level"
	SourceKey  = "source"
	TimeKey    = "time"
	MessageKey = "message"
	ErrorKey   = "error"
)

const (
	_traceValue = "TRACE"
	_debugValue = "DEBUG"
	_infoValue  = "INFO"
	_warnValue  = "WARN"
	_errorValue = "ERROR"
)

var _defaultLogger = NewLogger()

func Trace(ctx context.Context, msg string, args ...Arg) {
	_defaultLogger.Trace(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...Arg) {
	_defaultLogger.Info(ctx, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...Arg) {
	_defaultLogger.Debug(ctx, msg, args...)
}

func Warn(ctx context.Context, err error, msg string, args ...Arg) {
	_defaultLogger.Warn(ctx, err, msg, args...)
}

func Error(ctx context.Context, err error, msg string, args ...Arg) {
	_defaultLogger.Error(ctx, err, msg, args...)
}

func Write(ctx context.Context, level Level, msg string, args ...Arg) {
	_defaultLogger.Write(ctx, level, msg, args...)
}

type Format uint

const (
	FormatText Format = iota
	FormatJSON
)

type Options struct {
	Writer io.Writer
	Format Format

	LevelKey string
	Level    Level

	SourceKey string
	AddSource bool
	Skip      int

	TimeKey    string
	TimeFormat string

	MessageKey string
}

type Option func(o Options) Options

func WithFormat(format Format) Option {
	return func(o Options) Options {
		o.Format = format
		return o
	}
}

func WithLevel(key string, level Level) Option {
	return func(o Options) Options {
		o.LevelKey = key
		o.Level = level
		return o
	}
}

func WithSource(key string) Option {
	return func(o Options) Options {
		o.SourceKey = key
		o.AddSource = true
		return o
	}
}

func WithSkip(skip int) Option {
	return func(o Options) Options {
		o.Skip = skip
		return o
	}
}

func WithTime(key, format string) Option {
	return func(o Options) Options {
		o.TimeKey = key
		o.TimeFormat = format
		return o
	}
}

func WithMessageKey(key string) Option {
	return func(o Options) Options {
		o.MessageKey = key
		return o
	}
}

func WithWriter(w io.Writer) Option {
	return func(o Options) Options {
		o.Writer = w
		return o
	}
}

var NoContext = context.Background()

type optionChain []Option

func (c optionChain) apply() Options {
	o := Options{
		Writer:     os.Stdout,
		Format:     FormatText,
		LevelKey:   LevelKey,
		Level:      LevelDebug,
		SourceKey:  SourceKey,
		AddSource:  false,
		Skip:       0,
		TimeKey:    TimeKey,
		TimeFormat: time.RFC3339,
		MessageKey: MessageKey,
	}

	for _, opt := range c {
		o = opt(o)
	}

	return o
}
