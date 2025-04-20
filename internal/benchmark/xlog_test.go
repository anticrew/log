package benchmark

import (
	"io"
	"testing"

	"github.com/anticrew/log"
	"github.com/stretchr/testify/assert"
)

func Benchmark_Xlog(b *testing.B) {
	b.Run("logfmt", func(b *testing.B) {
		l := log.NewLogger(log.WithFormat(log.FormatLogFmt), log.WithWriter(io.Discard),
			log.WithLevel("", log.LevelInfo), log.WithSource(log.SourceKey))

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkXlog(l)
		}
	})

	b.Run("text", func(b *testing.B) {
		l := log.NewLogger(log.WithFormat(log.FormatText), log.WithWriter(io.Discard),
			log.WithLevel("", log.LevelInfo), log.WithSource(log.SourceKey))

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkXlog(l)
		}
	})

	b.Run("json", func(b *testing.B) {
		l := log.NewLogger(log.WithFormat(log.FormatJSON), log.WithWriter(io.Discard),
			log.WithLevel("", log.LevelInfo), log.WithSource(log.SourceKey))

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkXlog(l)
		}
	})

	b.Run("json-any", func(b *testing.B) {
		l := log.NewLogger(log.WithFormat(log.FormatJSON), log.WithWriter(io.Discard),
			log.WithLevel("", log.LevelInfo), log.WithSource(log.SourceKey))

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			l.Trace(log.NoContext, "trace", log.Any("any", Entity{}))
			l.Debug(log.NoContext, "debug", log.Any("any", Entity{}))
			l.Info(log.NoContext, "info", log.Any("any", Entity{}))
			l.Warn(log.NoContext, assert.AnError, "warn", log.Any("any", Entity{}))
			l.Error(log.NoContext, assert.AnError, "error", log.Any("any", Entity{}))
		}
	})
}

func benchmarkXlog(l log.Logger) {
	l.Trace(log.NoContext, "trace")
	l.Debug(log.NoContext, "debug")
	l.Info(log.NoContext, "info")
	l.Warn(log.NoContext, assert.AnError, "warn")
	l.Error(log.NoContext, assert.AnError, "error")
}
