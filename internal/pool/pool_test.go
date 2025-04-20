package pool

import (
	"runtime/debug"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type value[T any] struct {
	value T
}

func TestNew(t *testing.T) {
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	p := NewPool(func() *value[string] {
		return &value[string]{
			value: "new",
		}
	})

	// Probabilistically, 75% of sync.Pool.Put calls will succeed when -race
	// is enabled (see ref below); attempt to make this quasi-deterministic by
	// brute force (i.e., put significantly more objects in the pool than we
	// will need for the test) in order to avoid testing without race enabled.
	//
	// ref: https://cs.opensource.google/go/go/+/refs/tags/go1.20.2:src/sync/pool.go;l=100-103
	for range 1_000 {
		p.Put(&value[string]{
			value: t.Name(),
		})
	}

	// Ensure that we always get the expected value. Note that this must only
	// run a fraction of the number of times that Put is called above.
	for range 10 {
		func() {
			x := p.Get()
			defer p.Put(x)
			require.Equal(t, t.Name(), x.value)
		}()
	}

	for range 1_000 {
		p.Get()
	}

	require.Equal(t, "new", p.Get().value)
}

func TestNew_Race(t *testing.T) {
	p := NewPool(func() *value[int] {
		return &value[int]{
			value: -1,
		}
	})

	var wg sync.WaitGroup
	defer wg.Wait()

	// Run a number of goroutines that read and write pool object fields to
	// tease out races.
	for i := range 1_000 {
		wg.Add(1)
		go func() {
			defer wg.Done()

			x := p.Get()
			defer p.Put(x)

			// Must both read and write the field.
			if n := x.value; n >= -1 {
				x.value = i
			}
		}()
	}
}

type resetValue[T any] struct {
	called *atomic.Bool
}

func (r resetValue[T]) Reset() {
	r.called.Store(true)
}

func TestNewResettablePool(t *testing.T) {
	p := NewResettablePool(func() *resetValue[int] {
		return &resetValue[int]{}
	})

	for range 1_000 {
		b := &atomic.Bool{}

		v := p.Get()
		v.called = b

		p.Put(v)

		assert.True(t, b.Load())
		v.called = nil
	}
}
