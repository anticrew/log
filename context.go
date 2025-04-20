package log

import "context"

type argsKey struct{}

var _contextArgKeys []string

func GetContextArgs(ctx context.Context) []Arg {
	if ctx == nil {
		return nil
	}

	args, ok := ctx.Value(argsKey{}).([]Arg)
	if ok {
		return args
	}

	for _, key := range _contextArgKeys {
		var arg Arg
		arg, ok = GetContextArg(ctx, key)
		if ok {
			args = append(args, arg)
		}
	}

	return args
}

func GetContextArg(ctx context.Context, key any) (a Arg, ok bool) {
	if key == nil {
		return a, false
	}

	a, ok = ctx.Value(argsKey{}).(Arg)
	return a, ok
}

func AddContextArgs(ctx context.Context, args ...Arg) context.Context {
	if ctx == nil || len(args) == 0 {
		return ctx
	}

	existingArgs := GetContextArgs(ctx)
	if len(existingArgs) > 0 {
		newArgs := make([]Arg, 0, len(existingArgs)+len(args))
		newArgs = append(newArgs, existingArgs...)
		newArgs = append(newArgs, args...)

		args = newArgs
	}

	return SetContextArgs(ctx, args...)
}

func SetContextArgs(ctx context.Context, args ...Arg) context.Context {
	if ctx == nil || len(args) == 0 {
		return ctx
	}

	return context.WithValue(ctx, argsKey{}, args)
}

func SetContextArg(ctx context.Context, key any, arg Arg) context.Context {
	if ctx == nil || key == nil {
		return ctx
	}
	return context.WithValue(ctx, key, arg)
}

type loggerKey struct{}

func GetContextLogger(ctx context.Context) Logger {
	if ctx == nil {
		return nil
	}

	l, ok := ctx.Value(loggerKey{}).(Logger)
	if !ok {
		return nil
	}

	return l
}

func SetContextLogger(ctx context.Context, l Logger) context.Context {
	if ctx == nil {
		return nil
	}

	return context.WithValue(ctx, loggerKey{}, l)
}
