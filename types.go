package log

import "context"

type Logger interface {
	With(args ...Arg) Logger
	WithContext(ctx context.Context) Logger

	Trace(ctx context.Context, msg string, args ...Arg)
	Debug(ctx context.Context, msg string, args ...Arg)
	Info(ctx context.Context, msg string, args ...Arg)
	Warn(ctx context.Context, err error, msg string, args ...Arg)
	Error(ctx context.Context, err error, msg string, args ...Arg)
	Write(ctx context.Context, level Level, msg string, args ...Arg)
}
