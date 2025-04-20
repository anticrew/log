package benchmark

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	zaplogfmt "github.com/sykesm/zap-logfmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapWriter struct {
	out io.Writer
}

func (z *zapWriter) Write(p []byte) (n int, err error) {
	return z.out.Write(p)
}

func (z *zapWriter) Sync() error {
	return nil
}

func Benchmark_Zap(b *testing.B) {
	b.Run("logfmt", func(b *testing.B) {
		l := zap.New(
			zapcore.NewCore(
				zaplogfmt.NewEncoder(zap.NewProductionEncoderConfig()),
				&zapWriter{
					out: io.Discard,
				},
				zapcore.InfoLevel,
			),
			zap.WithCaller(true),
		)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkZap(l)
		}
	})

	b.Run("text", func(b *testing.B) {
		l := zap.New(
			zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
				&zapWriter{
					out: io.Discard,
				},
				zapcore.InfoLevel,
			),
			zap.WithCaller(true),
		)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkZap(l)
		}
	})

	b.Run("json", func(b *testing.B) {
		l := zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				&zapWriter{
					out: io.Discard,
				},
				zapcore.InfoLevel,
			),
			zap.WithCaller(true),
		)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			benchmarkZap(l)
		}
	})

	b.Run("json-any", func(b *testing.B) {
		l := zap.New(
			zapcore.NewCore(
				zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
				&zapWriter{
					out: io.Discard,
				},
				zapcore.InfoLevel,
			),
			zap.WithCaller(true),
		)

		b.ResetTimer()
		b.ReportAllocs()

		for b.Loop() {
			l.Log(zap.DebugLevel-1, "debug", zap.Any("any", Entity{}))
			l.Debug("debug", zap.Any("any", Entity{}))
			l.Info("info", zap.Any("any", Entity{}))
			l.Warn("warn", zap.Any("any", Entity{}))
			l.Error("error", zap.Error(assert.AnError), zap.Any("any", Entity{}))
		}
	})
}

func benchmarkZap(l *zap.Logger) {
	l.Log(zap.DebugLevel-1, "debug")
	l.Debug("debug")
	l.Info("info")
	l.Warn("warn", zap.Error(assert.AnError))
	l.Error("error", zap.Error(assert.AnError))
}
