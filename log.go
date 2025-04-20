package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	// LevelKey
	// Ключ по умолчанию для записи уровня лога
	LevelKey = "level"

	// SourceKey
	// Ключ по умолчанию для записи источника (места вызова записи в коде)
	SourceKey = "source"

	// TimeKey
	// Ключ по умолчанию для записи временной метки
	TimeKey = "time"

	// MessageKey
	// Ключ по умолчанию для записи текстового сообщения
	MessageKey = "message"

	// ErrorKey
	// Ключ по умолчанию для записи текста ошибки
	ErrorKey = "error"
)

const (
	_traceValue = "TRACE"
	_debugValue = "DEBUG"
	_infoValue  = "INFO"
	_warnValue  = "WARN"
	_errorValue = "ERROR"
)

// Trace
// Записывает TRACE-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
// Если уровень Logger выше TRACE, то лог должен быть проигнорирован без обработки.
func Trace(ctx context.Context, msg string, args ...Arg) {
	defaultLoggerFor(ctx).Trace(ctx, msg, args...)
}

// Info
// Записывает INFO-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
// Если уровень Logger выше INFO, то лог должен быть проигнорирован без обработки.
func Info(ctx context.Context, msg string, args ...Arg) {
	defaultLoggerFor(ctx).Info(ctx, msg, args...)
}

// Debug
// Записывает DEBUG-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
// Если уровень Logger выше DEBUG, то лог должен быть проигнорирован без обработки.
func Debug(ctx context.Context, msg string, args ...Arg) {
	defaultLoggerFor(ctx).Debug(ctx, msg, args...)
}

// Warn
// Записывает WARN-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
// Если уровень Logger выше WARN, то лог должен быть проигнорирован без обработки.
func Warn(ctx context.Context, err error, msg string, args ...Arg) {
	defaultLoggerFor(ctx).Warn(ctx, err, msg, args...)
}

// Error
// Записывает ERROR-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
// Если уровень Logger выше ERROR, то лог должен быть проигнорирован без обработки.
func Error(ctx context.Context, err error, msg string, args ...Arg) {
	defaultLoggerFor(ctx).Error(ctx, err, msg, args...)
}

// Write
// Записывает лог с указанным уровнем, сообщением и аргументами, а также аргументами, переданными в ctx.
// Если уровень Logger выше указанного, то лог должен быть проигнорирован без обработки.
func Write(ctx context.Context, level Level, msg string, args ...Arg) {
	defaultLoggerFor(ctx).Write(ctx, level, msg, args...)
}

var _defaultLogger = NewLogger()

func SetDefault(l Logger) {
	if l == nil {
		return
	}

	_defaultLogger = l
}

func GetDefault() Logger {
	return _defaultLogger
}

func defaultLoggerFor(ctx context.Context) Logger {
	if l := GetContextLogger(ctx); l != nil {
		return l
	}

	return _defaultLogger
}

// Format
// Предписывает Logger формат записи логов
type Format uint

const (
	// FormatText
	// Позволяет записывать логи в произвольном текстовом формате. Может кардинально отличаться у разных реализаций
	FormatText Format = iota

	// FormatJSON
	// Позволяет записывать логи в JSON, формат общий для всех реализаций, но порядок атрибутов может отличаться
	FormatJSON

	// FormatLogFmt
	// Позволяет записывать логи в LogFmt, формат общий для всех реализаций, но порядок атрибутов может отличаться
	FormatLogFmt
)

func (f Format) String() string {
	switch f {
	case FormatText:
		return "text"
	case FormatJSON:
		return "json"
	case FormatLogFmt:
		return "logfmt"
	default:
		return fmt.Sprintf("Format<%d>", f)
	}
}

func (f Format) IsValid() bool {
	return f >= FormatText && f <= FormatLogFmt
}

type Options struct {
	// Writer
	// Поток вывода, в который записываются логи. По умолчанию используется io.Stdout
	Writer io.Writer

	// Format
	// Формат вывода, по умолчанию - FormatText
	Format Format

	// LevelKey
	// Ключ для записи уровня лога, по умолчанию - LevelKey
	LevelKey string

	// Level
	// Наименьший уровень логов, которые допустимо записывать, по умолчанию - LevelDebug
	Level Level

	// SourceKey
	// Ключ для записи источника, по умолчанию - SourceKey
	SourceKey string

	// AddSource
	// Флаг добавления источника к логу
	AddSource bool

	// Skip
	// Кол-во вызовов для пропуска, может использоваться для библиотечных вызовов по умолчанию - 0
	Skip int

	// TimeKey
	// Ключ для записи временной метки, по умолчанию - TimeKey
	TimeKey string

	// TimeFormat
	// Формат временной метки для записи, по умолчанию - time.RFC3339
	TimeFormat string

	// MessageKey
	// Ключ для записи текстового сообщения в лог, по умолчанию - MessageKey
	MessageKey string
}

// Option
// Опция-функция для настройки Logger
type Option func(o Options) Options

// WithFormat
// Определяет Format для Logger
func WithFormat(format Format) Option {
	if !format.IsValid() {
		return emptyOption
	}

	return func(o Options) Options {
		o.Format = format
		return o
	}
}

// WithLevel
// Определяет минимальный уровень Logger
func WithLevel(key string, level Level) Option {
	if len(key) == 0 {
		key = LevelKey
	}

	return func(o Options) Options {
		o.LevelKey = key
		o.Level = level
		return o
	}
}

// WithSource
// Включает запись источника по указанному ключу
func WithSource(key string) Option {
	if len(key) == 0 {
		key = SourceKey
	}

	return func(o Options) Options {
		o.SourceKey = key
		o.AddSource = true
		return o
	}
}

// WithSkip
// Устанавливает кол-во пропущенных вызовов при определении источника
func WithSkip(skip int) Option {
	return func(o Options) Options {
		o.Skip = skip
		return o
	}
}

// WithTime
// Устанавливает ключ и формат временной метки
func WithTime(key, format string) Option {
	if len(key) == 0 {
		key = TimeKey
	}

	if len(format) == 0 {
		format = time.RFC3339
	}

	return func(o Options) Options {
		o.TimeKey = key
		o.TimeFormat = format
		return o
	}
}

// WithMessageKey
// Устанавливает ключ для записи текстового сообщения
func WithMessageKey(key string) Option {
	return func(o Options) Options {
		o.MessageKey = key
		return o
	}
}

// WithWriter
// Устанавливает поток вывода логов
func WithWriter(w io.Writer) Option {
	return func(o Options) Options {
		o.Writer = w
		return o
	}
}

// emptyOption
// Пустая опция, возвращающая исходный Options
func emptyOption(o Options) Options {
	return o
}

// NoContext
// "Пустой" context.Context без аргументов и Logger
var NoContext = context.Background()

type optionChain []Option

func (c optionChain) apply(o Options) Options {
	for _, opt := range c {
		o = opt(o)
	}

	return o
}

// defaultOptions
// Опции, заполненные значениями по умолчанию
func defaultOptions() Options {
	return Options{
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
}

type Logger interface {
	// WithArgs
	// Создает новый Logger совмещающий в себе настройки текущего Logger и аргументы, указанные в args.
	// Если args отсутствуют, допустимо возвращать текущий Logger без изменений
	WithArgs(args ...Arg) Logger

	// WithContext
	// Создает новый Logger совмещающий в себе настройки текущего Logger и аргументы, переданные в ctx.
	// Если ctx отсутствует, допустимо возвращать текущий Logger без изменений
	WithContext(ctx context.Context) Logger

	// WithOptions
	// Создает новый Logger совмещающий в себе настройки текущего Logger и опции, указанные в opts.
	// Если opts отсутствуют, допустимо возвращать текущий Logger без изменений
	WithOptions(opts ...Option) Logger

	// Trace
	// Записывает TRACE-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
	// Если уровень Logger выше TRACE, то лог должен быть проигнорирован без обработки.
	Trace(ctx context.Context, msg string, args ...Arg)

	// Debug
	// Записывает DEBUG-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
	// Если уровень Logger выше DEBUG, то лог должен быть проигнорирован без обработки.
	Debug(ctx context.Context, msg string, args ...Arg)

	// Info
	// Записывает INFO-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
	// Если уровень Logger выше INFO, то лог должен быть проигнорирован без обработки.
	Info(ctx context.Context, msg string, args ...Arg)

	// Warn
	// Записывает WARN-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
	// Если уровень Logger выше WARN, то лог должен быть проигнорирован без обработки.
	Warn(ctx context.Context, err error, msg string, args ...Arg)

	// Error
	// Записывает ERROR-лог с указанным сообщением и аргументами, а также аргументами, переданными в ctx.
	// Если уровень Logger выше ERROR, то лог должен быть проигнорирован без обработки.
	Error(ctx context.Context, err error, msg string, args ...Arg)

	// Write
	// Записывает лог с указанным уровнем, сообщением и аргументами, а также аргументами, переданными в ctx.
	// Если уровень Logger выше указанного, то лог должен быть проигнорирован без обработки.
	Write(ctx context.Context, level Level, msg string, args ...Arg)
}
