package engine

import (
	"testing"

	"github.com/anticrew/log/internal/buffer"

	"github.com/stretchr/testify/assert"
)

func Test_isDirty(t *testing.T) {
	t.Parallel()

	type testCase struct {
		buf   *buffer.Buffer
		count int
		dirty bool
	}

	testCases := map[string]testCase{
		"empty-0": {
			buf:   buffer.New(),
			count: 0,
			dirty: false,
		},
		"content-0": {
			buf:   buffer.New().WriteString("content"),
			count: 0,
			dirty: true,
		},
		"empty-1": {
			buf:   buffer.New(),
			count: 1,
			dirty: false,
		},
		"content-1": {
			buf:   buffer.New().WriteString("content"),
			count: 1,
			dirty: true,
		},
		"a-1": {
			buf:   buffer.New().WriteString("a"),
			count: 1,
			dirty: false,
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.dirty, isDirty(test.buf, test.count))

			test.buf.Dispose()
		})
	}
}
