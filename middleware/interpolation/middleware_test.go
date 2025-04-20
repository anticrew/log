package interpolation

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/anticrew/log"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware_Handle(t *testing.T) {
	t.Parallel()

	type testCase struct {
		src      log.Record
		expected log.Record
		err      error
	}

	testCases := map[string]testCase{
		"attr-one": {
			src: log.Record{
				Message: "Hello {{ name }}",
				Attrs:   log.NewAttrs().Append(log.String("name", "world")),
			},
			expected: log.Record{
				Message: "Hello world",
				Attrs:   log.NewAttrs().Append(log.String("name", "world")),
			},
		},
		"err": {
			src: log.Record{
				Message: "Hello {{ name ",
				Attrs:   log.NewAttrs().Append(log.String("name", "world")),
			},
			expected: log.Record{
				Message: "Hello {{ name ",
				Attrs:   log.NewAttrs().Append(log.String("name", "world")),
			},
			err: ErrUnclosedKey,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			mw := New()
			actual, err := mw.Handle(&test.src)

			if test.err == nil {
				require.NoError(t, err)
				require.NotNil(t, actual)

				assert.Equal(t, test.expected.Level, actual.Level)
				assert.Equal(t, test.expected.Message, actual.Message)

				expectedAttrs := copyAttrs(test.expected.Attrs)
				actualAttrs := copyAttrs(actual.Attrs)

				assert.Len(t, actualAttrs, len(expectedAttrs))

				for i, attr := range actualAttrs {
					assert.Equal(t, expectedAttrs[i].Key, attr.Key)
					assert.Equal(t, expectedAttrs[i].Value.String(), attr.Value.String())
				}
			} else {
				require.ErrorIs(t, err, test.err)
			}
		})
	}
}

func copyAttrs(a *log.Attrs) []log.Attr {
	result := make([]log.Attr, 0, 64)

	a.Range(func(a log.Attr) bool {
		result = append(result, a)

		return false
	})

	return result
}
