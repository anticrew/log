package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/anticrew/log"
	"github.com/anticrew/log/internal/buffer"

	"github.com/stretchr/testify/assert"
)

func TestJsonMarshaler_Marshal(t *testing.T) {
	t.Parallel()

	type testCase struct {
		record   *log.Record
		expected string
		anyError bool
		err      error
	}

	dateUTC := time.Date(2025, 5, 1, 10, 0, 0, 0, time.UTC)

	testCases := map[string]testCase{
		"trace": {
			record: &log.Record{
				Time:  dateUTC,
				Level: log.LevelTrace,
				Attrs: log.NewAttrs(),
			},
			expected: `{"time":"2025-05-01T10:00:00Z","level":"TRACE","message":""}`,
		},
		"trace-message": {
			record: &log.Record{
				Time:    dateUTC,
				Level:   log.LevelTrace,
				Message: "hello world",
				Attrs:   log.NewAttrs(),
			},
			expected: `{"time":"2025-05-01T10:00:00Z","level":"TRACE","message":"hello world"}`,
		},
		"debug-attr": {
			record: &log.Record{
				Time:  dateUTC,
				Level: log.LevelDebug,
				Attrs: log.NewAttrs().Append(
					log.String("string", "value"),
					log.Int("int", -42),
					log.Uint64("uint", 42),
					log.Float64("float", 3.14159),
					log.Bool("bool", true),
					log.Duration("duration", time.Second),
					log.Time("when", dateUTC),
				),
			},
			expected: `{"time":"2025-05-01T10:00:00Z","level":"DEBUG","message":"","string":"value","int":-42,` +
				`"uint":42,"float":3.14159,"bool":true,"duration":"1s","when":"2025-05-01T10:00:00Z"}`,
		},
		"debug-any": {
			record: &log.Record{
				Time:  dateUTC,
				Level: log.LevelDebug,
				Attrs: log.NewAttrs().Append(
					log.Any("any", struct {
						Name string
						Age  int
					}{
						Name: "Author",
						Age:  42,
					}),
				),
			},
			expected: `{"time":"2025-05-01T10:00:00Z","level":"DEBUG","message":"","any":{"Name":"Author","Age":42}}`,
		},
		"debug-json-fail": {
			record: &log.Record{
				Time:  dateUTC,
				Level: log.LevelDebug,
				Attrs: log.NewAttrs().Append(
					log.Any("chan", make(chan int)),
				),
			},
			anyError: true,
		},
		"debug-attr-duplicate": {
			record: &log.Record{
				Time:  dateUTC,
				Level: log.LevelDebug,
				Attrs: log.NewAttrs().Append(
					log.String("string", "value2"),
				),
			},
			expected: `{"time":"2025-05-01T10:00:00Z","level":"DEBUG","message":"","string":"value2"}`,
		},
		"debug-key-duplicate": {
			record: &log.Record{
				Time:  dateUTC,
				Level: log.LevelDebug,
				Attrs: log.NewAttrs().Append(
					log.String("level", "9.0"),
				),
			},
			err: ErrKeyExists,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			defer test.record.Dispose()

			out := buffer.New()
			defer out.Dispose()

			m := NewJsonMarshaler(out)
			defer m.Dispose()

			err := m.Marshal(test.record)
			switch {
			case test.anyError:
				require.Error(t, err)
			case test.err != nil:
				require.ErrorIs(t, err, test.err)
			default:
				require.NoError(t, err)
			}

			if len(test.expected) == 0 {
				assert.Empty(t, out.String())
			} else {
				assert.JSONEq(t, test.expected, out.String())
			}
		})
	}
}

func TestJsonMarshaler_MarshalKeys(t *testing.T) {
	// Этот тест меняет статические настройки ключей для симуляции неадекватного пользователя.
	// Поэтому его НЕЛЬЗЯ запускать в параллельном режиме (он сломает другие тесты)!

	t.Run("level-as-time", func(t *testing.T) {
		oldLevelKey := LevelKey
		LevelKey = TimeKey
		defer func() {
			LevelKey = oldLevelKey
		}()

		out := buffer.New()
		defer out.Dispose()

		m := NewJsonMarshaler(out)
		defer m.Dispose()

		r := log.NewRecord(log.LevelDebug, "", log.NewAttrs())
		r.Dispose()

		err := m.Marshal(r)
		require.ErrorIs(t, err, ErrKeyExists)
	})

	t.Run("message-as-level", func(t *testing.T) {
		oldMessageKey := MessageKey
		MessageKey = LevelKey
		defer func() {
			MessageKey = oldMessageKey
		}()

		out := buffer.New()
		defer out.Dispose()

		m := NewJsonMarshaler(out)
		defer m.Dispose()

		r := log.NewRecord(log.LevelDebug, "", log.NewAttrs())
		r.Dispose()

		err := m.Marshal(r)
		require.ErrorIs(t, err, ErrKeyExists)
	})
}
