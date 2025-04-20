//go:build anticrew_log_zap

package log

import (
	"context"
	"io"
	"time"

	zaplogfmt "github.com/sykesm/zap-logfmt"

	"github.com/anticrew/go-x/pool"
	"github.com/anticrew/go-x/xio"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

type Arg = zap.Field

func Err(err error) Arg {
	s := "nil"
	if err != nil {
		s = err.Error()
	}

	return zap.String(ErrorKey, s)
}

func String(key, value string) Arg {
	return zap.String(key, value)
}

func Uint(key string, value uint) Arg {
	return zap.Uint(key, value)
}

func Uint64(key string, value uint64) Arg {
	return zap.Uint64(key, value)
}

func Int(key string, value int) Arg {
	return zap.Int(key, value)
}

func Int64(key string, value int64) Arg {
	return zap.Int64(key, value)
}

func Float32(key string, value float32) Arg {
	return zap.Float32(key, value)
}

func Float64(key string, value float64) Arg {
	return zap.Float64(key, value)
}

func Bool(key string, value bool) Arg {
	return zap.Bool(key, value)
}

func Time(key string, value time.Time) Arg {
	return zap.Time(key, value)
}

func Duration(key string, value time.Duration) Arg {
	return zap.Duration(key, value)
}

func Any(key string, value any) Arg {
	return zap.Any(key, value)
}

const (
	LevelTrace = zap.DebugLevel - 1
	LevelDebug = zap.DebugLevel
	LevelInfo  = zap.InfoLevel
	LevelWarn  = zap.WarnLevel
	LevelError = zap.ErrorLevel
)

type logger struct {
	log *zap.Logger
	opt Options
}

func NewLogger(options ...Option) Logger {
	return createFromOptions(defaultOptions(), options)
}

func createFromOptions(opt Options, options []Option) Logger {
	opt = optionChain(options).apply(opt)

	cfg := zap.NewProductionEncoderConfig()
	cfg.LevelKey = opt.LevelKey

	levels := newLevelsConfig()
	cfg.EncodeLevel = levels.encode

	cfg.CallerKey = opt.SourceKey
	if len(opt.SourceKey) > 0 {
		cfg.EncodeCaller = func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			buf := xio.NewBuffer()
			defer buf.Dispose()

			buf.WriteString(caller.TrimmedPath()).
				WriteByte(' ').
				WriteString(caller.Function)

			encoder.AppendByteString(buf.Bytes())
		}
	}

	if len(opt.TimeKey) > 0 {
		cfg.TimeKey = opt.TimeKey
	}
	if len(opt.TimeFormat) > 0 {
		cfg.EncodeTime = func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(t.Format(opt.TimeFormat))
		}
	}

	if len(opt.MessageKey) > 0 {
		cfg.MessageKey = opt.MessageKey
	}

	var encoder zapcore.Encoder
	switch opt.Format {
	case FormatText:
		encoder = zapcore.NewConsoleEncoder(cfg)
	case FormatJSON:
		encoder = zapcore.NewJSONEncoder(cfg)
	case FormatLogFmt:
		encoder = zaplogfmt.NewEncoder(cfg)
	}

	writer, ok := opt.Writer.(zapcore.WriteSyncer)
	if !ok {
		writer = &zapWriter{
			out: opt.Writer,
		}
	}

	l := zap.New(
		zapcore.NewCore(
			encoder,
			&zapWriter{
				out: writer,
			},
			opt.Level,
		),
		zap.WithCaller(opt.AddSource),
		zap.AddCallerSkip(opt.Skip+2), // +2 to skip logAttrs and level-dependent function
	)

	return &logger{
		log: l,
		opt: opt,
	}
}

func (l *logger) WithArgs(args ...Arg) Logger {
	return &logger{
		log: l.log.With(args...),
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

var _argsPool = pool.NewPool[[]Arg](func() []Arg {
	return make([]Arg, 0, 16)
})

func (l *logger) logAttrs(ctx context.Context, level Level, err error, msg string, args []Arg) {
	newArgs := _argsPool.Get()
	defer _argsPool.Put(newArgs)

	newArgs = append(newArgs, GetContextArgs(ctx)...)
	newArgs = append(newArgs, args...)

	if err != nil {
		newArgs = append(newArgs, Err(err))
	}

	l.log.Log(level, msg, newArgs...)
}

type levelsConfig struct {
	values map[Level]string
}

func newLevelsConfig() *levelsConfig {
	return &levelsConfig{
		values: map[Level]string{
			LevelTrace: _traceValue,
			LevelDebug: _debugValue,
			LevelInfo:  _infoValue,
			LevelWarn:  _warnValue,
			LevelError: _errorValue,
		},
	}
}

func (c *levelsConfig) encode(level Level, encoder zapcore.PrimitiveArrayEncoder) {
	if value, ok := c.values[level]; ok {
		encoder.AppendString(value)
	} else {
		encoder.AppendString(level.String())
	}
}

type zapWriter struct {
	out io.Writer
}

func (z *zapWriter) Write(p []byte) (n int, err error) {
	return z.out.Write(p)
}

func (z *zapWriter) Sync() error {
	return nil
}
