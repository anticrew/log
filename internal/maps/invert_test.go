package maps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvert(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		t.Parallel()

		var (
			src map[string]int
			dst map[int]string
		)

		assert.NotPanics(t, func() {
			dst = Invert(src)
		})

		assert.Nil(t, dst)
	})

	t.Run("empty", func(t *testing.T) {
		t.Parallel()

		var (
			src = make(map[string]int)
			dst map[int]string
		)

		assert.NotPanics(t, func() {
			dst = Invert(src)
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
			dst map[int]string
		)

		assert.NotPanics(t, func() {
			dst = Invert(src)
		})

		assert.NotNil(t, dst)
		assert.Equal(t, map[int]string{
			1: "a",
			2: "b",
			3: "c",
		}, dst)
	})
}
