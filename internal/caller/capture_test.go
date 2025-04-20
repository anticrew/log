package caller

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const MaxFrames = int(^uint(0)>>1) / 2

func loadPaths(t *testing.T) (goPath, work string) {
	winToUnix := func(s string) string {
		return strings.ReplaceAll(s, "\\", "/")
	}

	goPath = os.Getenv("GOPATH")
	if len(goPath) == 0 {
		homeDir, err := os.UserHomeDir()
		require.NoError(t, err)
		require.NotEmpty(t, homeDir)

		goPath = filepath.Join(homeDir, "go")
	}

	work, err := os.Getwd()
	require.NoError(t, err)

	return winToUnix(goPath), winToUnix(work)
}

func inside(skip int) (*Caller, func(), error) {
	return Capture(skip)
}

func Test_Capture(t *testing.T) {
	t.Parallel()

	gopath, work := loadPaths(t)

	type testCase struct {
		fn       func(skip int) (*Caller, func(), error)
		skip     int
		expected runtime.Frame
		err      error
	}

	testCases := map[string]testCase{
		"direct-skip-0": {
			fn:   Capture,
			skip: 0,
			expected: runtime.Frame{
				File:     work + "/capture_test.go",
				Line:     116, // where is Capture called
				Function: "github.com/anticrew/log/internal/caller.Test_Capture.func1.1",
			},
		},
		"direct-skip-1": {
			fn:   Capture,
			skip: 1,
			expected: runtime.Frame{
				File:     gopath + "/pkg/mod/github.com/stretchr/testify@v1.10.0/assert/assertions.go",
				Line:     1239, // where is Capture called
				Function: "github.com/stretchr/testify/assert.didPanic",
			},
		},
		"inside-skip-0": {
			fn:   inside,
			skip: 0,
			expected: runtime.Frame{
				File:     work + "/capture_test.go",
				Line:     37, // where is Capture called
				Function: "github.com/anticrew/log/internal/caller.inside",
			},
		},
		"inside-skip-1": {
			fn:   inside,
			skip: 1,
			expected: runtime.Frame{
				File:     work + "/capture_test.go",
				Line:     116, // where is Capture called
				Function: "github.com/anticrew/log/internal/caller.Test_Capture.func1.1",
			},
		},
		"inside-skip-2": {
			fn:   inside,
			skip: 2,
			expected: runtime.Frame{
				File:     gopath + "/pkg/mod/github.com/stretchr/testify@v1.10.0/assert/assertions.go",
				Line:     1239, // where is Capture called
				Function: "github.com/stretchr/testify/assert.didPanic",
			},
		},
		"err-no-frames": {
			fn:   Capture,
			skip: MaxFrames, // half of max of int
			err:  ErrNoFrames,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var (
				caller  *Caller
				dispose func()
				err     error
			)

			require.NotPanics(t, func() {
				caller, dispose, err = test.fn(test.skip)
			})

			if test.err == nil {
				require.NoError(t, err)

				assert.Equal(t, test.expected.File, caller.File)
				assert.Equal(t, test.expected.Line, caller.Line)
				assert.Equal(t, test.expected.Function, caller.Function)
			} else {
				require.ErrorIs(t, err, test.err)
			}

			if dispose != nil {
				dispose()
			}
		})
	}
}

func Benchmark_Capture(b *testing.B) {
	b.ReportAllocs()

	for b.Loop() {
		_, dispose, _ := Capture(0)
		dispose()
	}

	b.StopTimer()
}
