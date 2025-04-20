package benchmark

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Benchmark_Slog(b *testing.B) {
	b.Run("logfmt/text", func(b *testing.B) {
		//nolint: sloglint // benchmark requires text handler with io.Discard
		l := slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}))

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkSlog(b, l)
		}
	})

	b.Run("json", func(b *testing.B) {
		//nolint: sloglint // benchmark requires text handler with io.Discard
		l := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}))

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkSlog(b, l)
		}
	})

	b.Run("json-any", func(b *testing.B) {
		//nolint: sloglint // benchmark requires text handler with io.Discard
		l := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}))

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			l.Log(b.Context(), slog.LevelDebug-4, "trace", slog.Any("any", Entity{}))
			l.Debug("debug", slog.Any("any", Entity{}))
			l.Info("info", slog.Any("any", Entity{}))
			l.Warn("warn", slog.Any("error", assert.AnError), slog.Any("any", Entity{}))
			l.Error("error", slog.Any("error", assert.AnError), slog.Any("any", Entity{}))
		}
	})
}

func benchmarkSlog(b *testing.B, l *slog.Logger) {
	l.Log(b.Context(), slog.LevelDebug-4, "trace")
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn", slog.Any("error", assert.AnError))
	l.Error("error", slog.Any("error", assert.AnError))
}
