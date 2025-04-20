package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClear(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		var (
			src, dst map[string]int
		)

		assert.NotPanics(t, func() {
			dst = Clear(src)
		})

		assert.Nil(t, dst)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		var (
			src = make(map[string]int)
			dst map[string]int
		)

		assert.NotPanics(t, func() {
			dst = Clear(src)
		})

		assert.NotNil(t, dst)
		assert.Empty(t, dst)
	})

	t.Run("data", func(t *testing.T) {
		t.Parallel()

		var (
			src = map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			}
			dst map[string]int
		)

		assert.NotPanics(t, func() {
			dst = Clear(src)
		})

		assert.NotNil(t, dst)
		assert.Empty(t, dst)
	})
}
