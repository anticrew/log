package log

import "context"

// argsKey
// Структура-ключ для хранения набора аргументов в контексте
type argsKey struct{}

// GetContextArgs
// Возвращает набор аргументов, содержащийся в context.Context по ключу argsKey{}. Если context.Context не указан или по ключу ничего не
// найдено, возвращает пустой набор (nil)
func GetContextArgs(ctx context.Context) []Arg {
	if ctx == nil {
		return nil
	}

	args, ok := ctx.Value(argsKey{}).([]Arg)
	if !ok {
		return nil
	}

	return args
}

// AddContextArgs
// Дополняет набор аргументов, содержащийся в context.Context по ключу argsKey{}, указанным новым набором аргументов, добавляя их в конец
// списка. Если context.Context не указан или список аргументов пуст, возвращает исходный context.Context без изменений
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

// SetContextArgs
// Устанавливает набор аргументов в context.Context по ключу argsKey{}, перезаписывая имеющийся набор, если таковой был. Если
// context.Context не указан или список аргументов пуст, возвращает исходный context.Context без изменений
func SetContextArgs(ctx context.Context, args ...Arg) context.Context {
	if ctx == nil || len(args) == 0 {
		return ctx
	}

	return context.WithValue(ctx, argsKey{}, args)
}

// GetContextArg
// Возвращает один аргумент и статус получения, содержащийся в context.Context по указанному ключу
func GetContextArg(ctx context.Context, key any) (a Arg, ok bool) {
	if key == nil {
		return a, false
	}

	a, ok = ctx.Value(argsKey{}).(Arg)
	return a, ok
}

// SetContextArg
// Возвращает context.Context, содержащий указанный аргумент по указанному ключу. Если context.Context или ключ не указан, возвращает
// исходный context.Context без изменений
func SetContextArg(ctx context.Context, key any, arg Arg) context.Context {
	if ctx == nil || key == nil {
		return ctx
	}
	return context.WithValue(ctx, key, arg)
}

// loggerKey
// Структура-ключ для хранения Logger в контексте
type loggerKey struct{}

// GetContextLogger
// Возвращает Logger, полученный из context.Context по ключу loggerKey{} или nil, если по ключу ничего не найдено
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

// SetContextLogger
// Возвращает context.Context, содержащий указанный Logger по ключу loggerKey{} или nil, если context.Context не указан. Если в указаанном
// context.Context уже содержался Logger, он станет недоступен из дочернего контекста
func SetContextLogger(ctx context.Context, l Logger) context.Context {
	if ctx == nil {
		return nil
	}

	return context.WithValue(ctx, loggerKey{}, l)
}
