//go:build anticrew_log_slog

package log

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/anticrew/log/internal/caller"

	"github.com/anticrew/go-x/pool"
)

type Level = slog.Level

const (
	LevelTrace = slog.LevelDebug - 4
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
)

type Arg = slog.Attr

func Err(err error) Arg {
	s := "nil"
	if err != nil {
		s = err.Error()
	}

	return slog.String(ErrorKey, s)
}

func String(key, value string) Arg {
	return slog.String(key, value)
}

func Uint(key string, value uint) Arg {
	return slog.Uint64(key, uint64(value))
}

func Uint64(key string, value uint64) Arg {
	return slog.Uint64(key, value)
}

func Int(key string, value int) Arg {
	return slog.Int(key, value)
}

func Int64(key string, value int64) Arg {
	return slog.Int64(key, value)
}

func Float32(key string, value float32) Arg {
	return slog.Float64(key, float64(value))
}

func Float64(key string, value float64) Arg {
	return slog.Float64(key, value)
}

func Bool(key string, value bool) Arg {
	return slog.Bool(key, value)
}

func Time(key string, value time.Time) Arg {
	return slog.Time(key, value)
}

func Duration(key string, value time.Duration) Arg {
	return slog.Duration(key, value)
}

func Any(key string, value any) Arg {
	return slog.Any(key, value)
}

type logger struct {
	log *slog.Logger

	levels *levelsConfig

	opt Options
}

func NewLogger(options ...Option) Logger {
	return createFromOptions(defaultOptions(), options)
}

func createFromOptions(opt Options, options []Option) Logger {
	opt = optionChain(options).apply(opt)

	config := newLevelsConfig(opt.LevelKey, opt.Level)
	handlerOpt := &slog.HandlerOptions{
		AddSource: false, // always false, we handle source manually
		Level:     opt.Level,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			switch a.Key {
			case slog.LevelKey:
				return config.replaceValue(a)

			case slog.TimeKey:
				return replaceTime(opt, a)

			case slog.MessageKey:
				return replaceMessage(opt, a)

			default:
				return a
			}
		},
	}

	var handler slog.Handler
	switch opt.Format {
	case FormatText, FormatLogFmt:
		handler = slog.NewTextHandler(opt.Writer, handlerOpt)
	case FormatJSON:
		handler = slog.NewJSONHandler(opt.Writer, handlerOpt)
	}

	l := slog.New(handler)

	return &logger{
		log:    l,
		levels: config,
		opt:    opt,
	}
}

var _anyPool = pool.NewPool(func() []any {
	return make([]any, 0, 16)
}, pool.WithReset(func(args []any) []any {
	return args[:0]
}))

func (l *logger) WithArgs(args ...Arg) Logger {
	a := _anyPool.Get()
	defer _anyPool.Put(a)

	for _, arg := range args {
		a = append(a, arg)
	}

	return &logger{
		log:    l.log.With(a...),
		levels: l.levels,
		opt:    l.opt,
	}
}

func (l *logger) WithContext(ctx context.Context) Logger {
	return l.WithArgs(GetContextArgs(ctx)...)
}

func (l *logger) WithOptions(options ...Option) Logger {
	return createFromOptions(l.opt, options)
}

func (l *logger) Trace(ctx context.Context, msg string, args ...Arg) {
	l.logAttrs(ctx, LevelTrace, nil, msg, args)
}

func (l *logger) Debug(ctx context.Context, msg string, args ...Arg) {
	l.logAttrs(ctx, LevelDebug, nil, msg, args)
}

func (l *logger) Info(ctx context.Context, msg string, args ...Arg) {
	l.logAttrs(ctx, LevelInfo, nil, msg, args)
}

func (l *logger) Warn(ctx context.Context, err error, msg string, args ...Arg) {
	l.logAttrs(ctx, LevelWarn, err, msg, args)
}

func (l *logger) Error(ctx context.Context, err error, msg string, args ...Arg) {
	l.logAttrs(ctx, LevelError, err, msg, args)
}

func (l *logger) Write(ctx context.Context, level Level, msg string, args ...Arg) {
	l.logAttrs(ctx, level, nil, msg, args)
}

var _argsPool = pool.NewPool(func() []Arg {
	return make([]Arg, 0, 16)
}, pool.WithReset(func(args []Arg) []Arg {
	return args[:0]
}))

func (l *logger) logAttrs(ctx context.Context, level Level, err error, msg string, args []Arg) {
	if l.levels.enabled < level {
		return
	}

	newArgs := _argsPool.Get()
	defer _argsPool.Put(newArgs)

	if ctx != NoContext {
		newArgs = append(newArgs, GetContextArgs(ctx)...)
	}

	newArgs = append(newArgs, args...)

	if err != nil {
		newArgs = append(newArgs, Err(err))
	}

	if l.opt.AddSource {
		newArgs = append(newArgs, l.getSourceArg(2))
	}

	l.log.LogAttrs(ctx, level, msg, newArgs...)
}

func (l *logger) getSourceArg(skip int) Arg {
	src, err := caller.Take(l.opt.Skip + skip + 1)
	if err == nil {
		return slog.String(l.opt.SourceKey, src)
	}

	return slog.String(l.opt.SourceKey, fmt.Sprintf("(error = %v)", err))
}

type levelsConfig struct {
	key     string
	enabled Level
	values  map[Level]slog.Value
}

func newLevelsConfig(key string, enabled Level) *levelsConfig {
	return &levelsConfig{
		key:     key,
		enabled: enabled,
		values: map[Level]slog.Value{
			LevelTrace: slog.StringValue(_traceValue),
			LevelDebug: slog.StringValue(_debugValue),
			LevelInfo:  slog.StringValue(_infoValue),
			LevelWarn:  slog.StringValue(_warnValue),
			LevelError: slog.StringValue(_errorValue),
		},
	}
}

func (c *levelsConfig) replaceValue(a slog.Attr) slog.Attr {
	level, ok := a.Value.Any().(slog.Level)
	if !ok {
		return a
	}

	a.Key = c.key

	var value slog.Value
	if value, ok = c.values[level]; ok {
		a.Value = value
	}

	return a
}

func replaceTime(opt Options, a slog.Attr) slog.Attr {
	if len(opt.TimeKey) > 0 {
		a.Key = opt.TimeKey
	}
	if len(opt.TimeFormat) > 0 {
		a.Value = slog.StringValue(a.Value.Time().Format(opt.TimeFormat))
	}

	return a
}

func replaceMessage(opt Options, a slog.Attr) slog.Attr {
	if len(opt.MessageKey) == 0 {
		return a
	}

	a.Key = opt.MessageKey
	return a
}
