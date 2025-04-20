package caller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Take(t *testing.T) {
	t.Parallel()

	type testCase struct {
		skip     int
		expected string
		err      error
	}

	testCases := map[string]testCase{
		"skip-0": {
			skip:     0,
			expected: "caller/take_test.go:48 github.com/anticrew/log/internal/caller.Test_Take.func1.1",
		},
		"skip-1": {
			skip:     1,
			expected: "assert/assertions.go:1239 github.com/stretchr/testify/assert.didPanic",
		},
		"skip-2": {
			skip:     2,
			expected: "assert/assertions.go:1310 github.com/stretchr/testify/assert.NotPanics",
		},
		"err-no-frames": {
			skip: MaxFrames, // half of max of int
			err:  ErrNoFrames,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				actual string
				err    error
			)

			require.NotPanics(t, func() {
				actual, err = Take(test.skip)
			})

			if test.err == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, test.err)
			}

			assert.Equal(t, test.expected, actual)
		})
	}
}

func Benchmark_Take(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		s, err := Take(0)
		_ = s
		_ = err
	}
}
