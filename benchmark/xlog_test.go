package benchmark

import (
	"github.com/anticrew/log"
	"github.com/anticrew/log/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Benchmark_Xlog(b *testing.B) {
	b.Run("logfmt", func(b *testing.B) {
		l := log.New(engine.NewEngine(&nopWriter{}, log.PrintModeLogFmt), log.LevelInfo)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkXlog(l)
		}
	})

	b.Run("json", func(b *testing.B) {
		l := log.New(engine.NewEngine(&nopWriter{}, log.PrintModeJson), log.LevelInfo)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkXlog(l)
		}
	})

	b.Run("json-any", func(b *testing.B) {
		l := log.New(engine.NewEngine(&nopWriter{}, log.PrintModeJson), log.LevelInfo)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			l.Trace("trace", log.Any("any", Entity{}))
			l.Debug("debug", log.Any("any", Entity{}))
			l.Info("info", log.Any("any", Entity{}))
			l.Warn("warn", log.Any("error", assert.AnError), log.Any("any", Entity{}))
			l.Error(assert.AnError, "error", log.Any("any", Entity{}))
		}
	})

	b.Run("pretty", func(b *testing.B) {
		l := log.New(engine.NewEngine(&nopWriter{}, log.PrintModePretty), log.LevelInfo)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkXlog(l)
		}
	})

}

func benchmarkXlog(l log.Logger) {
	l.Trace("trace")
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn", log.Err(assert.AnError))
	l.Error(assert.AnError, "error")
}
